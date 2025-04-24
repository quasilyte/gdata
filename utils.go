package gdata

import (
	"os"
	"path/filepath"
)

func fileExists(root *os.Root, path string) bool {
	_, err := root.Stat(path)
	return err == nil
}

func mkdirAll(root *os.Root, path string) error {
	if fileExists(root, path) {
		return nil
	}
	return root.Mkdir(path, 0755)
}

func removeAllInRoot(root *os.Root, path string) error {
	dir, err := root.Open(path)
	if err != nil {
		return err
	}
	defer dir.Close()

	entries, err := dir.ReadDir(-1)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			if err := removeAllInRoot(root, entryPath); err != nil {
				return err
			}
		} else {
			if err := root.Remove(entryPath); err != nil {
				return err
			}
		}
	}

	return root.Remove(path)
}
