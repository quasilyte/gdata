//go:build darwin || linux || windows

package gdata

import (
	"io"
	"os"
	"path/filepath"
)

// filesystemDataManager implements FS-based storage.
// It works on Windows, MacOS, Linux, and Android.
//
// Objects are folders, files are their properties.
type filesystemDataManager struct {
	root     *os.Root
	dataPath string
}

func (m *filesystemDataManager) objectPath(objectKey string) string {
	return filepath.Join(m.dataPath, objectKey)
}

func (m *filesystemDataManager) ObjectPropPath(objectKey, propKey string) string {
	return filepath.Join(m.dataPath, objectKey, propKey)
}

func (m *filesystemDataManager) ListObjectProps(objectKey string) ([]string, error) {
	p := m.objectPath(objectKey)
	if !fileExists(m.root, p) {
		return nil, nil
	}
	dir, err := m.root.Open(p)
	if err != nil {
		return nil, err
	}
	files, err := dir.ReadDir(-1)
	if err != nil {
		return nil, err
	}
	result := make([]string, len(files))
	for i, f := range files {
		result[i] = f.Name()
	}
	return result, nil
}

func (m *filesystemDataManager) SaveObjectProp(objectKey, propKey string, data []byte) error {
	p := m.objectPath(objectKey)
	if !fileExists(m.root, p) {
		if err := m.root.Mkdir(p, os.ModePerm); err != nil {
			return err
		}
	}

	file, err := m.root.Create(p)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

func (m *filesystemDataManager) LoadObjectProp(objectKey, propKey string) ([]byte, error) {
	p := m.ObjectPropPath(objectKey, propKey)
	if !fileExists(m.root, p) {
		return nil, nil
	}

	file, err := m.root.Open(p)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

func (m *filesystemDataManager) ObjectPropExists(objectKey, propKey string) bool {
	return fileExists(m.root, m.ObjectPropPath(objectKey, propKey))
}

func (m *filesystemDataManager) ObjectExists(objectKey string) bool {
	return fileExists(m.root, m.objectPath(objectKey))
}

func (m *filesystemDataManager) DeleteObjectProp(objectKey, propKey string) error {
	p := m.ObjectPropPath(objectKey, propKey)
	if !fileExists(m.root, p) {
		return nil
	}
	return m.root.Remove(p)
}

func (m *filesystemDataManager) DeleteObject(objectKey string) error {
	p := m.objectPath(objectKey)
	// Since RemoveAll returns a nil error for a non-existing
	// path, we can avoid doing an explicit fileExists call.
	return removeAllInRoot(m.root, p)
}
