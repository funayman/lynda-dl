package downloader

import (
	"net/http"

	"github.com/funayman/lynda-dl/course"
)

type Downloader struct {
	client *http.Client
}

func New() *Downloader {
	return &Downloader{client: &http.Client{}}
}

func (d *Downloader) Get(c course.LyndaCourse) error {
	return nil
}
