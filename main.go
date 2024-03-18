package main

import (
	"bytes"
	"context"
	"flag"
	"github.com/schollz/progressbar/v3"
	"log"
	"os"
	"path"
	"strings"
	"sync"
)

var videoUrl = flag.String("u", "", "input your url")
var concurrentDownloadNum = flag.Int64("c", 5, "concurrent download num range 0-30")

func main() {
	defer func() {
		r := recover()
		if r != nil {
			log.Printf("Fatal panic: %v", r)
			log.Println("cleaning file...")
			_ = os.RemoveAll(path.Base(*videoUrl))
			os.Exit(1)
		}
	}()

	flag.Parse()
	if *videoUrl == "" {
		log.Panicln("empty url")
	}
	if *concurrentDownloadNum <= 0 || *concurrentDownloadNum > 30 {
		log.Panicln("error set")
	}
	log.Println("video_url ", *videoUrl)
	log.Println("concurrent download num ", *concurrentDownloadNum)

	log.Println("getting playlist and video cover...")
	// make dir
	folderPath := path.Base(*videoUrl)
	videoPath := path.Join(folderPath, folderPath+".mp4")
	err := os.Mkdir(folderPath, 0755)
	if err != nil {
		log.Panicln("Error creating folder:", err)
	}

	m3u8Url := mustGetM3u8Url(*videoUrl, folderPath)
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
		log.Panicln("segUriList empty")
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
				log.Panicln(err)
			}
			buf := bytes.Buffer{}
			buf.Write(decrypted)
			id2Ts.Store(i, buf)

			_ = bar.Add(1)

		}(i, TsUrl)
	}

	// 拼接视频
	waiter.Wait()
	log.Println("merging videos...")

	var TsFileList []bytes.Buffer
	for i := 0; i < total; i++ {
		tsBuf, exist := id2Ts.Load(i)
		if exist {
			t := tsBuf.(bytes.Buffer)
			TsFileList = append(TsFileList, t)
		}
	}

	m := MergeTsFileListToSingleMp4Req{
		TsFileList: TsFileList,
		OutputMp4:  videoPath,
		Ctx:        context.Background(),
	}

	err = MergeTsFileListToSingleMp4(m)
	if err != nil {
		log.Panicln("err merging:", err)
	}

	log.Println("downloaded ", videoPath, " !")
}
