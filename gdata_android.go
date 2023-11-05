//go:build android

package gdata

import (
	"errors"
	"os"
	"path/filepath"
)

type dataManager struct {
	dataPath string
}

func newDataManager(config Config) (dataManagerImpl, error) {
	app, err := detectAndroidApp()
	if err != nil {
		return nil, err
	}
	dataPath := filepath.Join("/data/data/", app)
	if !fileExists(dataPath) {
		return nil, errors.New("can't find the app data directory")
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

func (m *dataManager) DeleteData(itemKey string) error {
	itemPath := m.DataPath(itemKey)
	if !fileExists(itemPath) {
		return nil
	}
	return os.Remove(itemPath)
}

func (m *dataManager) LoadData(itemKey string) ([]byte, error) {
	itemPath := m.DataPath(itemKey)
	if !fileExists(itemPath) {
		return nil, nil
	}
	return os.ReadFile(itemPath)
}

func detectAndroidApp() (string, error) {
	data, err := os.ReadFile("/proc/self/cmdline")
	if err != nil {
		return "", err
	}
	// Trim any potential "\n" and remove the null bytes.
	copied := make([]byte, 0, len(data))
	for _, ch := range data {
		switch ch {
		case 0, '\n':
			continue
		}
		copied = append(copied, ch)
	}
	result := string(copied)
	if result == "" {
		return "", errors.New("got empty output from /proc/self/cmdline")
	}
	return result, nil
}
