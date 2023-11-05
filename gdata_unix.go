//go:build (darwin || linux) && !android

package gdata

import (
	"os"
	"path/filepath"
)

type dataManager struct {
	dataPath string
}

func newDataManager(config Config) (dataManagerImpl, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dataPath := filepath.Join(home, ".local", "share", config.AppName)
	if err := mkdirAll(dataPath); err != nil {
		return nil, err
	}
	m := &dataManager{
		dataPath: dataPath,
	}
	return m, nil
}

func (m *dataManager) DataPath(itemKey string) string {
	return filepath.Join(m.dataPath, itemKey)
}

func (m *dataManager) DataExists(itemKey string) bool {
	return fileExists(m.DataPath(itemKey))
}

func (m *dataManager) SaveData(itemKey string, data []byte) error {
	return os.WriteFile(m.DataPath(itemKey), data, 0o666)
}

func (m *dataManager) LoadData(itemKey string) ([]byte, error) {
	itemPath := m.DataPath(itemKey)
	if !fileExists(itemPath) {
		return nil, nil
	}
	return os.ReadFile(itemPath)
}
