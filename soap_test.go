package main

import (
	"log"
	"testing"
)

func TestMustDownloadM3u8(t *testing.T) {
	m3u8Url := ""
	mustDownloadM3u8(m3u8Url)
}

func TestMustGetM3u8Url(t *testing.T) {
	videoUrl := ""
	m3u8Url := mustGetM3u8Url(videoUrl)
	log.Println("获取 m3u8 链接成功: ", m3u8Url)
}
