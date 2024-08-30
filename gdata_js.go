package gdata

import (
	"errors"
	"syscall/js"

	_ "embed"
)

// Since we can't store objects inside the localStorage,
// we store a metadata string per every object.
// This metadata is a single string consisting of all
// "filenames" (related property keys).

//go:embed "gdata.js"
var jsCode string

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
	lib := js.Global().Get("window").Call("eval", jsCode)
	js.Global().Set("___gdata", lib)
	return m, nil
}

func (m *dataManager) getLib() js.Value {
	return js.Global().Get("___gdata")
}

func (m *dataManager) ObjectPropPath(objectKey, propKey string) string {
	return m.getLib().Call("objectPropPath", m.appName, objectKey, propKey).String()
}

func (m *dataManager) ListObjectProps(objectKey string) ([]string, error) {
	v := m.getLib().Call("listObjectProps", m.appName, objectKey)
	if v.IsNull() {
		return nil, nil
	}
	result := make([]string, v.Length())
	for i := 0; i < v.Length(); i++ {
		result[i] = v.Index(i).String()
	}
	return result, nil
}

func (m *dataManager) SaveObjectProp(objectKey, propKey string, data []byte) error {
	m.getLib().Call("saveObjectProp", m.appName, objectKey, propKey, string(data))
	return nil
}

func (m *dataManager) LoadObjectProp(objectKey, propKey string) ([]byte, error) {
	result := m.getLib().Call("loadObjectProp", m.appName, objectKey, propKey)
	if result.IsNull() {
		return nil, nil
	}
	return []byte(result.String()), nil
}

func (m *dataManager) ObjectPropExists(objectKey, propKey string) bool {
	return m.getLib().Call("objectPropExists", m.appName, objectKey, propKey).Bool()
}

func (m *dataManager) ObjectExists(objectKey string) bool {
	return m.getLib().Call("objectExists", m.appName, objectKey).Bool()
}

func (m *dataManager) DeleteObjectProp(objectKey, propKey string) error {
	m.getLib().Call("deleteObjectProp", m.appName, objectKey, propKey)
	return nil
}

func (m *dataManager) DeleteObject(objectKey string) error {
	m.getLib().Call("deleteObject", m.appName, objectKey)
	return nil
}
