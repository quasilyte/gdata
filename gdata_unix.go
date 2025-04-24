//go:build (darwin || linux) && !android

package gdata

import (
	"os"
	"path/filepath"
)

func newDataManager(config Config) (dataManagerImpl, error) {

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dataPath := filepath.Join(home, ".local", "share", config.AppName)
	root, err := os.OpenRoot(dataPath)
	if err != nil {
		return nil, err
	}
	if err := mkdirAll(root, dataPath); err != nil {
		return nil, err
	}
	m := &filesystemDataManager{
		root:     root,
		dataPath: dataPath,
	}
	return m, nil
}
