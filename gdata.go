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
	ObjectPropPath(objectKey, propKey string) string

	ListObjectProps(objectKey string) ([]string, error)

	SaveObjectProp(objectKey, propKey string, data []byte) error

	LoadObjectProp(objectKey, propKey string) ([]byte, error)
	ReadObjectProp(objectKey, propKey string, buf []byte) (int, error)

	ObjectPropExists(objectKey, propKey string) bool
	ObjectExists(objectKey string) bool

	DeleteObjectProp(objectKey, propKey string) error
	DeleteObject(objectKey string) error
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

// ListObjectProps returns all object property keys.
//
// If there is no such object, a nil slice is returned.
//
// The returned error is usually a file operation error.
func (m *Manager) ListObjectProps(objectKey string) ([]string, error) {
	return m.impl.ListObjectProps(objectKey)
}

// SaveItem writes to the object's property with the associated key.
//
// Using an empty propKey is allowed.

// Saving to an existing key overwrites it.
// If object did not exist, it will be created automatically.
//
// The returned error is usually a file write error.
func (m *Manager) SaveObjectProp(objectKey, propKey string, data []byte) error {
	return m.impl.SaveObjectProp(objectKey, fixPropKey(propKey), data)
}

// LoadObjectProp reads from the object's property using the provided key.
//
// Using an empty propKey is allowed.
//
// Loading a non-existing propKey is not an error, a nil slice will be returned.
// The same rule applies to a non-existing object.
//
// If you want to know whether some object or a property exists,
// use ObjectExists or ObjectPropExists method.
//
// The returned error is usually a file read error.
func (m *Manager) LoadObjectProp(objectKey, propKey string) ([]byte, error) {
	return m.impl.LoadObjectProp(objectKey, fixPropKey(propKey))
}

// ReadObjectProp is like LoadObjectProp, but reads only up
// to the specified size (in bytes). This size is specified by the buf's len.
//
// This is usually useful on FS-like environments where save files
// can be large, but the header reading operation is needed to perform
// an efficient list-like option.
//
// It can also be useful if you want to reuse the destination buffer for
// reading smaller data into a single reusable buffer.
//
// The size is usually a fixed-size header block size.
// If you have a dynamic header size, there should be a way for
// you to know the size before calling this function (e.g. put it in the propKey).
//
// If data is smaller than size, the returned byte slice will have a length
// equal to the actual data size read.
func (m *Manager) ReadObjectProp(objectKey, propKey string, buf []byte) (int, error) {
	return m.impl.ReadObjectProp(objectKey, fixPropKey(propKey), buf)
}

// ObjectExists reports whether the object was saved before.
func (m *Manager) ObjectExists(objectKey string) bool {
	return m.impl.ObjectExists(objectKey)
}

// ObjectPropExists reports whether an object has the specified property.
//
// Using an empty propKey is allowed.
//
// If object itself doesn't exists, the function will report false.
func (m *Manager) ObjectPropExists(objectKey, propKey string) bool {
	return m.impl.ObjectPropExists(objectKey, fixPropKey(propKey))
}

// DeleteObject removes the object along with all of its properties (if any).
//
// Be careful with this function: it removes the data permanently.
// There is no way to undo it.
//
// Trying to delete a non-existing object is not an error.
//
// The returned error is usually a file operation error.
func (m *Manager) DeleteObject(objectKey string) error {
	return m.impl.DeleteObject(objectKey)
}

// DeleteObjectProp removes the object's property data.
//
// Using an empty propKey is allowed.
//
// Be careful with this function: it removes the data permanently.
// There is no way to undo it.
//
// Trying to delete a non-existing propKey is not an error.
// If object doesn't exist, it's not an error either.
//
// Deleting the last object's property does not delete
// the object itself. Use DeleteObject if you want to
// delete the object completely.
//
// The returned error is usually a file operation error.
func (m *Manager) DeleteObjectProp(objectKey, propKey string) error {
	return m.impl.DeleteObjectProp(objectKey, fixPropKey(propKey))
}

// ObjectPropPath returns a unique object property path.
//
// Using an empty propKey is allowed.
//
// On platforms with filesystem storage, it's an absolute file path.
// On other platforms it's just some unique resource identifier.
//
// You can't treat it as a filesystem path unless you know what you're doing.
// It's safe to use it for debugging and for things like map keys.
//
// Note that ObjectPropPath returns a potential data path, it doesn't care if
// the item actually exists or not.
// Use ObjectPropExists() method first if you need to know whether this
// objectKey is used.
func (m *Manager) ObjectPropPath(objectKey, propKey string) string {
	return m.impl.ObjectPropPath(objectKey, fixPropKey(propKey))
}

func fixPropKey(propKey string) string {
	if propKey == "" {
		return "_objdat"
	}
	return propKey
}
