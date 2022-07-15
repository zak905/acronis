package main

import (
	"flag"
	"log"

	"github.com/zak905/acronis-assignement/client/downloader"
)

func main() {

	var serverAddress string
	var downloadDestination string

	flag.StringVar(&serverAddress, "server", "http://localhost:8080", "the address of the file server")
	flag.StringVar(&downloadDestination, "dest", "./", "where to download the files")

	flag.Parse()

	if err := downloader.New(serverAddress, downloadDestination, 'A').DownloadFilesWithFirstCharOccurence(); err != nil {
		log.Fatal(err.Error())
	}
}
