package main

import (
	"bytes"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMerge(t *testing.T) {
	buf := bytes.Buffer{}
	filepath.WalkDir("./output", func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(d.Name(), ".mp4") {
			readFile, _ := os.ReadFile("output/" + d.Name())
			buf.Write(readFile)
		}

		return nil
	})

	log.Println("开始转档")
	fileName := "output/output.mp4"
	err := os.WriteFile(fileName, buf.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}

	mergeVideo(fileName)
}
