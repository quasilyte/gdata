package gdata

import (
	"os"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func mkdirAll(path string) error {
	if fileExists(path) {
		return nil
	}
	return os.MkdirAll(path, 0755)
}
