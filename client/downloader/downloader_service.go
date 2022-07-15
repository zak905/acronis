package downloader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/zak905/acronis-assignement/client/resultmanager"
	"golang.org/x/net/html"
)

type DownloaderService struct {
	resultManager       *resultmanager.ResultManager
	serverAddress       string
	wg                  sync.WaitGroup
	client              http.Client
	charachter          rune
	downloadDestination string
}

func New(serverAddress, downloadDestination string, charachter rune) *DownloaderService {
	return &DownloaderService{serverAddress: serverAddress, resultManager: resultmanager.New(), wg: sync.WaitGroup{}, client: *http.DefaultClient, charachter: charachter, downloadDestination: downloadDestination}
}

func (d *DownloaderService) DownloadFilesWithFirstCharOccurence() error {
	res, err := http.DefaultClient.Get(d.serverAddress)
	if err != nil {
		return fmt.Errorf("unable to GET from server %s: %w", d.serverAddress, err)
	}

	doc, err := html.Parse(res.Body)
	if err != nil {
		return fmt.Errorf("error while parsing html received from server: %w", err)
	}

	if err := d.parseAndProcess(doc); err != nil {
		return fmt.Errorf("error processing file list: %w", err)
	}

	log.Println("========================== Result:")
	log.Printf("Downloading the following files with charachter %c occuring at position %d: \n", d.charachter, d.resultManager.GetPosition())
	for _, fileURL := range d.resultManager.GetFilesToDownload() {
		log.Println(fileURL)
		if err := d.download(fileURL); err != nil {
			return fmt.Errorf("error downloading results: %w", err)
		}
	}
	log.Println("Download done")

	return nil
}

//parse file serve file list html response
//and process every link to a file
func (d *DownloaderService) parseAndProcess(root *html.Node) error {
	//this is because Parse() wrapps the markup into <html> to make a full html doc
	// html -> body -> pre -> a
	link := root.FirstChild.FirstChild.NextSibling.FirstChild.FirstChild

	multiError := newMultiError()

	//we go through all the files
	for link != nil {
		if link.Data == "a" {
			fileHref := getHrefFromAnchor(link)
			if fileHref == "" {
				return fmt.Errorf("did not find any href attribute for %s", link.DataAtom)
			}
			fileURL := d.serverAddress + "/" + getHrefFromAnchor(link)

			d.wg.Add(1)

			//download in parallel
			go func(url string) {
				if err := d.process(url); err != nil {
					multiError.add(url, err)
				}
				d.wg.Done()
			}(fileURL)
		}

		link = link.NextSibling
	}

	//wait for all the go routines to finish
	d.wg.Wait()

	if multiError.hasErrors() {
		return multiError
	}

	return nil
}

//the download logic
//it checks the content length first using head to avoid downloading the whole file
//then it downlaods byte by byte and checks whether it matches the charachter
func (d *DownloaderService) process(url string) error {
	//we use head to get the file size
	res, err := d.client.Head(url)
	if err != nil {
		return err
	}

	i := 0

	for i < int(res.ContentLength) {
		//if there is already a better result, we escape here
		if i > d.resultManager.GetPosition() && d.resultManager.GetPosition() != 0 {
			break
		}
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		//more details on https://developer.mozilla.org/en-US/docs/Web/HTTP/Range_requests
		req.Header.Set("Accept-Ranges", "bytes")
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", i, i))

		res2, err := d.client.Do(req)
		if err != nil {
			return err
		}

		b, err := io.ReadAll(res2.Body)
		if err != nil {
			return err
		}

		if string(b) == string(d.charachter) {
			d.resultManager.Add(i+1, url)
			break
		}

		i++
	}

	//for debug only
	log.Printf("downloaded %d bytes out of %d from file %s before making a decision\n", i, res.ContentLength, url)

	return nil
}

//download the file and write it to the destination
func (d *DownloaderService) download(fileURL string) error {
	pathTokens := strings.Split(fileURL, "/")

	fileName := pathTokens[len(pathTokens)-1]

	res, err := d.client.Get(fileURL)
	if err != nil {
		return err
	}

	f, err := os.Create(path.Clean(d.downloadDestination + "/" + fileName))
	if err != nil {
		return err
	}

	defer f.Close()

	t := io.TeeReader(res.Body, f)

	_, err = io.ReadAll(t)

	return err
}

//find the link pointed to by href attribute
func getHrefFromAnchor(link *html.Node) string {
	for _, a := range link.Attr {
		if a.Key == "href" {
			return a.Val
		}
	}

	return ""
}
