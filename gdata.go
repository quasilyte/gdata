package gdata

import (
	"errors"
)

type Manager struct {
	impl dataManagerImpl
}

type dataManagerImpl interface {
	SaveData(itemKey string, data []byte) error
	LoadData(itemKey string) ([]byte, error)
	DataExists(itemKey string) bool
	DataPath(itemKey string) string
}

type Config struct {
	AppName string
}

func Open(config Config) (*Manager, error) {
	if config.AppName == "" {
		return nil, errors.New("config.AppName can't be empty")
	}
	m, err := newDataManager(config)
	if err != nil {
		return nil, err
	}
	return &Manager{impl: m}, nil
}

func (m *Manager) SaveItem(itemKey string, data []byte) error {
	return m.impl.SaveData(itemKey, data)
}

func (m *Manager) LoadItem(itemKey string) ([]byte, error) {
	return m.impl.LoadData(itemKey)
}

func (m *Manager) ItemExists(itemKey string) bool {
	return m.impl.DataExists(itemKey)
}

func (m *Manager) ItemPath(itemKey string) string {
	return m.impl.DataPath(itemKey)
}
