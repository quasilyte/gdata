//go:build android

package gdata

import (
	"errors"
	"os"
	"path/filepath"
)

func newDataManager(config Config) (dataManagerImpl, error) {
	app, err := detectAndroidApp()
	if err != nil {
		return nil, err
	}
	dataPath := filepath.Join("/data/data/", app)
	if !fileExists(dataPath) {
		return nil, errors.New("can't find the app data directory")
	}
	m := &filesystemDataManager{
		dataPath: dataPath,
	}
	return m, nil
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
