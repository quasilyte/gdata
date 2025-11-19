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
	if !fileExists(p) {
		return nil, nil
	}
	files, err := os.ReadDir(p)
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
	if !fileExists(p) {
		if err := os.MkdirAll(p, os.ModePerm); err != nil {
			return err
		}
	}
	return os.WriteFile(filepath.Join(p, propKey), data, 0o666)
}

func (m *filesystemDataManager) LoadObjectProp(objectKey, propKey string) ([]byte, error) {
	p := m.ObjectPropPath(objectKey, propKey)
	if !fileExists(p) {
		return nil, nil
	}
	return os.ReadFile(p)
}

func (m *filesystemDataManager) ReadObjectProp(objectKey, propKey string, buf []byte) (int, error) {
	p := m.ObjectPropPath(objectKey, propKey)

	f, err := os.Open(p)
	if err != nil {
		if !os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}
	defer f.Close()

	n, err := io.ReadFull(f, buf)
	switch {
	case err == io.ErrUnexpectedEOF:
		// Not an error.
	case err != nil:
		return n, err
	}

	return n, nil
}

func (m *filesystemDataManager) ObjectPropExists(objectKey, propKey string) bool {
	return fileExists(m.ObjectPropPath(objectKey, propKey))
}

func (m *filesystemDataManager) ObjectExists(objectKey string) bool {
	return fileExists(m.objectPath(objectKey))
}

func (m *filesystemDataManager) DeleteObjectProp(objectKey, propKey string) error {
	p := m.ObjectPropPath(objectKey, propKey)
	if !fileExists(p) {
		return nil
	}
	return os.Remove(p)
}

func (m *filesystemDataManager) DeleteObject(objectKey string) error {
	p := m.objectPath(objectKey)
	// Since RemoveAll returns a nil error for a non-existing
	// path, we can avoid doing an explicit fileExists call.
	return os.RemoveAll(p)
}
