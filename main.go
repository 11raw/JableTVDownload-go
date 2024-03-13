package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"log"
	"path"
	"strings"
	"sync"
)

var videoUrl = flag.String("u", "", "input your url")
var concurrentDownloadNum = flag.Int64("c", 5, "concurrent download num range 0-30")

func main() {
	flag.Parse()
	if *videoUrl == "" {
		fmt.Println("empty url")
		return
	}
	if *concurrentDownloadNum <= 0 || *concurrentDownloadNum > 30 {
		fmt.Println("error set")
		return
	}
	fmt.Println("video_url ", *videoUrl)
	fmt.Println("concurrent download num ", *concurrentDownloadNum)

	fmt.Println("getting playlist...")
	m3u8Url := mustGetM3u8Url(*videoUrl)
	m3u8UrlDir := strings.Replace(m3u8Url, path.Base(m3u8Url), "", 1)

	m3u8Media := mustDownloadM3u8(m3u8Url)
	log.Println("got playlist from ", m3u8Url)
	makeTsWg := sync.WaitGroup{}
	// 获取 ts 文件列表
	var segUriList []string
	total := 0
	for _, segment := range m3u8Media.Segments {
		if segment != nil {
			segUriList = append(segUriList, m3u8UrlDir+segment.URI)
			makeTsWg.Add(1)
			total++
		}
	}

	if len(segUriList) == 0 {
		log.Fatal("segUriList empty")
	}

	bar := progressbar.Default(int64(len(segUriList)))
	// 获取解密文件
	m3u8Uri := m3u8Media.Key.URI
	m3uiIv := m3u8Media.Key.IV
	vt := getVt(m3uiIv)
	decryptKey := mustGetDecryptKey(m3u8UrlDir + "/" + m3u8Uri)

	waiter := sync.WaitGroup{}
	limiter := make(chan struct{}, *concurrentDownloadNum)
	// 下载并解密
	id2Ts := sync.Map{}
	for i, TsUrl := range segUriList {
		limiter <- struct{}{}
		waiter.Add(1)

		go func(i int, TsUrl string) {
			defer func() { <-limiter }()
			defer waiter.Done()

			segTsResponse := download(TsUrl)

			decrypted, err := decrypt(decryptKey, vt, segTsResponse)
			if err != nil {
				log.Fatal(err)
			}
			buf := bytes.Buffer{}
			buf.Write(decrypted)
			id2Ts.Store(i, buf)

			_ = bar.Add(1)

		}(i, TsUrl)
	}

	// 拼接视频
	waiter.Wait()
	fmt.Println("merging videos...")

	var TsFileList []bytes.Buffer
	for i := 0; i < total; i++ {
		tsBuf, exist := id2Ts.Load(i)
		if exist {
			t := tsBuf.(bytes.Buffer)
			TsFileList = append(TsFileList, t)
		}
	}

	path.Base(*videoUrl)
	outMp4 := path.Base(*videoUrl) + ".mp4"
	m := MergeTsFileListToSingleMp4_Req{
		TsFileList: TsFileList,
		OutputMp4:  outMp4,
		Ctx:        context.Background(),
	}

	err := MergeTsFileListToSingleMp4(m)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("downloaded ", outMp4, " !")
}
