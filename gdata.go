package gdata

import (
	"errors"
)

// Manager implements the main gamedata operations.
// You can create a new manager by using Open function.
type Manager struct {
	impl dataManagerImpl
}

// dataManagerImpl is an interface every platform-specific storage provider implements.
type dataManagerImpl interface {
	SaveData(itemKey string, data []byte) error
	LoadData(itemKey string) ([]byte, error)
	DataExists(itemKey string) bool
	DataPath(itemKey string) string
}

// Config affects the created gamedata manager behavior.
type Config struct {
	// AppName is used as a part of the key used to store the game data.
	//
	// The exact effect depends on the platform, but generally it doesn't have
	// to reflect the application name perfectly.
	//
	// You need to use the same AppName to make sure that the game can
	// then load the previously saved data.
	// If you want to separate the data, use suffixes: "app" and "app2" data
	// will be stored completely independently.
	//
	// An empty app name is not allowed.
	AppName string
}

// Open attempts to create a gamedata manager.
//
// There are various cases when it can fail and you need to be
// ready to handle that situation and run the game without the save states.
// For instance, on wasm platforms it's using a localStorage which can be disabled.
// In this case, a non-nil error will be returned and the game should continue
// without any attempts to load or save data.
//
// One gamedata manager per game is enough.
// You need to pass it explicitely as a part of your game's context.
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

// SaveItem writes the data with the associated itemKey.
//
// Saving an item to the existing key overwrites it.
//
// The returned error is usually a file write error.
func (m *Manager) SaveItem(itemKey string, data []byte) error {
	return m.impl.SaveData(itemKey, data)
}

// LoadItem reads the data saved with the itemKey.
//
// Loading a non-existing itemKey is not an error, a nil slice will be returned.
// If you want to know whether some itemKey exists or not, use ItemExists() method.
//
// The returned error is usually a file read error.
func (m *Manager) LoadItem(itemKey string) ([]byte, error) {
	return m.impl.LoadData(itemKey)
}

// ItemExists reports whether the itemKey was saved before.
// An existing key will result in a non-nil data being read with LoadItem().
func (m *Manager) ItemExists(itemKey string) bool {
	return m.impl.DataExists(itemKey)
}

// ItemPath returns a unique itemKey path.
//
// On platforms with filesystem storage, it's an absolute file path.
// On other platforms it just some unique resource identifier.
//
// You can't treat it as a filesystem path unless you're knowing what you're doing.
// It's safe to use it for debugging and for things like map keys.
//
// Note that ItemPath returns a potential item path, it doesn't care if
// the item actually exists or not.
// Use ItemExists() method first if you need to know whether this itemKey is used.
func (m *Manager) ItemPath(itemKey string) string {
	return m.impl.DataPath(itemKey)
}
