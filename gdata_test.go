package gdata_test

import (
	"bytes"
	"testing"

	"github.com/quasilyte/gdata"
)

func TestSaveLoad(t *testing.T) {
	m, err := gdata.Open(gdata.Config{
		AppName: "gdata_test",
	})
	if err != nil {
		t.Fatal(err)
	}

	const testItemKey = "testitem.txt"

	// Deleting a potentially non-existing item.
	// For a second run it acts like a state reset function.
	if err := m.DeleteItem(testItemKey); err != nil {
		t.Fatalf("delete should never result in an error here (got %v)", err)
	}

	if m.ItemExists(testItemKey) {
		t.Fatalf("%s item should not exist yet", testItemKey)
	}

	data := []byte("example data")
	if err := m.SaveItem(testItemKey, data); err != nil {
		t.Fatalf("error saving %s data", testItemKey)
	}

	if !m.ItemExists(testItemKey) {
		t.Fatalf("%s item should exist after a successful save operation", testItemKey)
	}

	loadedData, err := m.LoadItem(testItemKey)
	if err != nil {
		t.Fatalf("loading %s error: %v", testItemKey, err)
	}

	if !bytes.Equal(data, loadedData) {
		t.Fatalf("saved and loaded data mismatch:\nwant: %q\n have: %q", string(data), string(loadedData))
	}

	// Now we're deleting an existing item.
	if err := m.DeleteItem(testItemKey); err != nil {
		t.Fatalf("delete should never result in an error here (got %v)", err)
	}

	if m.ItemExists(testItemKey) {
		t.Fatalf("%s item should not exist after a successful delete operation", testItemKey)
	}
}
