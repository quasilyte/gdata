package gdata

import (
	"errors"
	"os"
	"path/filepath"
)

func newDataManager(config Config) (dataManagerImpl, error) {
	appData := os.Getenv("AppData")
	if appData == "" {
		return nil, errors.New("AppData env var is undefined")
	}
	dataPath := filepath.Join(appData, config.AppName)
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
