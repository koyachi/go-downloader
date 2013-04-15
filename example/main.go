package main

import (
	"../"
	"fmt"
	//"github.com/koyachi/go-downloader"
	"github.com/gregjones/gogress"
)

type Any interface{}
type ChannelReceiver func(Any)

type MailBox struct {
	channel chan *downloader.Downloader
	handler ChannelReceiver
}

func NewMailBox() *MailBox {
	mb := &MailBox{
		channel: make(chan *downloader.Downloader),
		handler: updateDownloadProgress,
	}
	return mb
}

var mailbox = NewMailBox()

func eventloop() {
	select {
	case d := <-mailbox.channel:
		mailbox.handler(d)
	}
}

func download(url string) {
	d := downloader.New(url)
	mailbox.channel <- d

	err := d.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("\nEND\n")
}

func updateDownloadProgress(data Any) {
	d := data.(*downloader.Downloader)
	fmt.Printf("download[%s] start\n", d.Url)
	p := gogress.NewProgressBar(100)
	p.Start()
	total := 0
	for progress := range d.ProgressCh {
		total += progress
		//fmt.Printf("received:%d, total = (%d / %d)\n", progress, total, d.Size)
		p.Update(int64(100.0 * total / d.Size))
	}
	fmt.Printf("download(%s) end.\n", d.Url)
}

func main() {
	//
	// mock user input
	//

	//url := "http://golang.org/pkg/net/http/"
	//url := "http://unkar.org/r/mnewsplus/1284379385"
	//url := "http://www.yahoo.co.jp"
	// http://ubu.com/sound/snow.html
	//url := "http://ubumexico.centro.org.mx/sound/snow_michael/music-for-piano-whistling/Snow-Michael_Music-For-Piano-Whistling_1-01-Falling-Starts_Begining.mp3"
	url := "http://ubumexico.centro.org.mx/sound/snow_michael/last-lp/Snow-Michael_Last-LP_01-Wu.mp3"
	command := "download"

	go eventloop()

	switch command {
	case "download":
		download(url)
	default:
	}
}
