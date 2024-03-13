package main

import (
	"fmt"
	"os/exec"
)

func mergeVideo(fileName string) {
	cmdArgs := []string{
		"-i", fileName,
		"-c:v", "h264_videotoolbox",
		"-b:v", "3M",
		"-threads", "5",
		"-preset", "superfast",
		"result.mp4",
	}

	cmd := exec.Command("ffmpeg", cmdArgs...)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Failed to encode video:", err)
		return
	}
}
