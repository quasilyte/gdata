package gdata_test

import (
	"bytes"
	"testing"

	"github.com/quasilyte/gdata/v2"
)

func TestSaveLoad(t *testing.T) {
	m, err := gdata.Open(gdata.Config{
		AppName: "gdata_test",
	})
	if err != nil {
		t.Fatal(err)
	}

	const (
		testObjectKey = "obj"
		testItemKey   = "testitem.txt"
	)

	for round := 0; round < 3; round++ {
		for i := 0; i < 2; i++ {
			// Deleting a potentially non-existing item.
			// For a second run it acts like a state reset function.
			if err := m.DeleteObjectProp(testObjectKey, testItemKey); err != nil {
				t.Fatalf("delete should never result in an error here (got %v)", err)
			}
		}
		for i := 0; i < 2; i++ {
			if err := m.DeleteObject(testObjectKey); err != nil {
				t.Fatalf("delete should never result in an error here (got %v)", err)
			}
		}

		if m.ObjectPropExists(testObjectKey, testItemKey) {
			t.Fatalf("%s.%s item should not exist yet", testObjectKey, testItemKey)
		}
		if m.ObjectExists(testObjectKey) {
			t.Fatalf("%s item should not exist yet", testObjectKey)
		}

		data := []byte("example data")
		if err := m.SaveObjectProp(testObjectKey, testItemKey, data); err != nil {
			t.Fatalf("error saving %s data", testItemKey)
		}

		if !m.ObjectPropExists(testObjectKey, testItemKey) {
			t.Fatalf("%s.%s item should exist after a successful save operation", testObjectKey, testItemKey)
		}
		if !m.ObjectExists(testObjectKey) {
			t.Fatalf("%s item should exist after a successful save operation", testObjectKey)
		}

		loadedData, err := m.LoadObjectProp(testObjectKey, testItemKey)
		if err != nil {
			t.Fatalf("loading %s error: %v", testItemKey, err)
		}

		if !bytes.Equal(data, loadedData) {
			t.Fatalf("saved and loaded data mismatch:\nwant: %q\n have: %q", string(data), string(loadedData))
		}

		// Now we're deleting an existing item.
		if err := m.DeleteObjectProp(testObjectKey, testItemKey); err != nil {
			t.Fatalf("delete should never result in an error here (got %v)", err)
		}

		if m.ObjectPropExists(testObjectKey, testItemKey) {
			t.Fatalf("%s item should not exist after a successful delete operation", testItemKey)
		}
		if !m.ObjectExists(testObjectKey) {
			t.Fatalf("%s object should remain existing after removing its property", testObjectKey)
		}

		if err := m.DeleteObject(testObjectKey); err != nil {
			t.Fatalf("error deleting %s object", testObjectKey)
		}
		if m.ObjectExists(testObjectKey) {
			t.Fatalf("%s object should not exist after being deleted", testObjectKey)
		}
	}
}
