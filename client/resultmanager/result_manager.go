package resultmanager

import (
	"sync"
)

type ResultManager struct {
	m            *sync.Mutex
	bestPosition int
	toDownload   []string
}

func New() *ResultManager {
	return &ResultManager{m: &sync.Mutex{}}
}

func (r *ResultManager) Add(postion int, filename string) bool {
	r.m.Lock()
	defer r.m.Unlock()

	if postion == r.bestPosition {
		r.toDownload = append(r.toDownload, filename)
		return true
	}

	if r.bestPosition == 0 || postion < r.bestPosition {
		r.bestPosition = postion
		r.toDownload = []string{filename}
		return true
	}

	return false
}

func (r *ResultManager) GetFilesToDownload() []string {
	return r.toDownload
}

func (r *ResultManager) GetPosition() int {
	return r.bestPosition
}
