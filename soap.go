package main

import (
	"bufio"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/grafov/m3u8"
	"io"
	"log"
	"net/http"
	"regexp"
)

var client = &http.Client{}

var reg = regexp.MustCompile("https:\\/\\/[^\\\"\\s]+\\.m3u8")

// mustGetM3u8Url 获取 m3u8 链接
func mustGetM3u8Url(link string) string {
	u := launcher.New().Bin("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome").MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()

	pageResponse := browser.MustPage(link).MustWaitStable()
	h5 := pageResponse.MustHTML()
	rawSearch := reg.Find([]byte(h5))
	if len(rawSearch) == 0 {
		log.Fatal("no m3u8 url")
	}

	browser.MustClose()

	return string(rawSearch)
}

func mustDownloadM3u8(link string) *m3u8.MediaPlaylist {
	response, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = response.Body.Close() }()

	playlist, _, err := m3u8.DecodeFrom(bufio.NewReader(response.Body), true)
	if err != nil {
		log.Fatal(err)
	}

	media, exist := playlist.(*m3u8.MediaPlaylist)
	if !exist {
		log.Fatal("error m3u8")
	}

	return media
}

func mustGetDecryptKey(tsUrl string) []byte {
	response, err := http.Get(tsUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = response.Body.Close() }()
	readAll, _ := io.ReadAll(response.Body)

	return readAll
}

func download(link string) []byte {
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Close = true

	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = response.Body.Close() }()

	readAll, _ := io.ReadAll(response.Body)

	return readAll
}
