package main

import (
	"os"
)

func writeMp4(name string, data []byte) error {
	err := os.WriteFile(name, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
