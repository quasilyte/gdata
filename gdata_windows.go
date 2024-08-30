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
	if err := mkdirAll(dataPath); err != nil {
		return nil, err
	}
	m := &filesystemDataManager{
		dataPath: dataPath,
	}
	return m, nil
}
