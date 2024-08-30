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
	if err := mkdirAll(dataPath); err != nil {
		return nil, err
	}
	m := &filesystemDataManager{
		dataPath: dataPath,
	}
	return m, nil
}
