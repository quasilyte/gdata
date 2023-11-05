package gdata

import (
	"errors"
	"syscall/js"
)

type dataManager struct {
	appName string
}

func newDataManager(config Config) (dataManagerImpl, error) {
	const evalArgument = `
        let ___result = true;
        try {
            const ___key = "storage__test__";
            window.localStorage.setItem(___key, null);
            window.localStorage.removeItem(___key);
        } catch (e) {
            ___result = false;
        }
        ___result`
	hasLocalStorage := js.Global().Get("window").Call("eval", evalArgument).Bool()
	if !hasLocalStorage {
		return nil, errors.New("localStorage is not available")
	}
	m := &dataManager{
		appName: config.AppName,
	}
	return m, nil
}

func (m *dataManager) DataPath(itemKey string) string {
	return "_gdata_" + m.appName + "_" + itemKey
}

func (m *dataManager) DataExists(itemKey string) bool {
	result := js.Global().Get("localStorage").Call("getItem", m.DataPath(itemKey))
	return !result.IsNull()
}

func (m *dataManager) SaveData(itemKey string, data []byte) error {
	js.Global().Get("localStorage").Call("setItem", m.DataPath(itemKey), string(data))
	return nil
}

func (m *dataManager) DeleteData(itemKey string) error {
	js.Global().Get("localStorage").Call("removeItem", m.DataPath(itemKey))
	return nil
}

func (m *dataManager) LoadData(itemKey string) ([]byte, error) {
	result := js.Global().Get("localStorage").Call("getItem", m.DataPath(itemKey))
	if result.IsNull() {
		return nil, nil
	}
	return []byte(result.String()), nil
}
