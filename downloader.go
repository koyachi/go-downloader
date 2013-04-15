package downloader

import (
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Downloader struct {
	Url        string
	ProgressCh chan int
	Size       int
}

func New(url string) *Downloader {
	return &Downloader{
		Url:        url,
		ProgressCh: make(chan int),
		Size:       0,
	}
}

func (d *Downloader) updateSize(header http.Header) error {
	strVal := header.Get("Content-Length")
	if len(strVal) > 0 {
		val, err := strconv.Atoi(strVal)
		if err != nil {
			return err
		}
		d.Size = val
		return nil
	}
	return nil
}

func (d *Downloader) updateSizeByHead() error {
	res, err := http.Head(d.Url)
	if err != nil {
		return err
	}
	d.updateSize(res.Header)
	return nil
}

// start download and save as file.
func (d *Downloader) Start() error {
	d.updateSizeByHead()

	dialFunc := func(network, addr string) (net.Conn, error) {
		conn, err := net.Dial(network, addr)
		if err != nil {
			return nil, err
		}
		return NewCountableConnection(conn, d.ProgressCh), nil
	}
	transport := &http.Transport{
		Dial: dialFunc,
	}
	client := &http.Client{Transport: transport}

	res, err := client.Get(d.Url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if d.Size == 0 {
		d.updateSize(res.Header)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fileNameParts := strings.Split(d.Url, "/")
	fileName := fileNameParts[len(fileNameParts)-1]
	if fileName == "" {
		// TODO: tailling slash
		fileName = "a.downloaded.file"
	}
	if err := ioutil.WriteFile(fileName, data, os.FileMode(0666)); err != nil {
		return err
	}

	return nil
}
