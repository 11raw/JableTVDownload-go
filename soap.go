package main

import (
	"bufio"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/grafov/m3u8"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

var client = &http.Client{}

var reg = regexp.MustCompile("https:\\/\\/[^\\\"\\s]+\\.m3u8")

// mustGetM3u8Url 获取 m3u8 链接
func mustGetM3u8Url(link, folderPath string) string {
	u := launcher.New().Bin("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome").MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()

	pageResponse := browser.MustPage(link).MustWaitStable()
	h5 := pageResponse.MustHTML()
	rawSearch := reg.Find([]byte(h5))
	if len(rawSearch) == 0 {
		log.Fatal("no m3u8 url")
	}

	browser.MustClose()

	getCover(strings.NewReader(h5), path.Join(folderPath, folderPath+".jpg"))

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

func getCover(htmlFile io.Reader, coverPath string) {
	doc, err := goquery.NewDocumentFromReader(htmlFile)
	if err != nil {
		log.Fatal("Error parsing HTML:", err)
	}

	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		metaContent, _ := s.Attr("content")
		if metaContent != "" && strings.Contains(metaContent, "preview.jpg") {
			resp, err := http.Get(metaContent)
			if err != nil {
				log.Fatal("Error downloading cover:", err)
			}
			defer func() { _ = resp.Body.Close() }()
			readAll, _ := io.ReadAll(resp.Body)
			if err = os.WriteFile(coverPath, readAll, 0755); err != nil {
				log.Fatal("Error writing cover:", err)
			}

			log.Printf("Cover downloaded as %s\n", coverPath)
		}
	})
}
