package store

import (
	"gotest.tools/assert"
	"testing"
)

func TestInitTable(t *testing.T) {
	// Action
	PersistentKeyValueStore, err := NewPersistentKeyValueStore("./BookmarkStore.sqlite", "KeyValueTable")
	PersistentKeyValueStore.InitTable()
	PersistentKeyValueStore.Clear()

	// Assert
	assert.Equal(t, err, nil)
}

func TestSetValue(t *testing.T) {
	// Assume
	PersistentKeyValueStore, _ := NewPersistentKeyValueStore("./BookmarkStore.sqlite", "KeyValueTable")
	PersistentKeyValueStore.InitTable()
	PersistentKeyValueStore.Clear()

	// Action
	errKey1 := PersistentKeyValueStore.SetValue("key1", "value1")
	errKey2 := PersistentKeyValueStore.SetValue("key2", "value2")

	assert.Equal(t, errKey1, nil)
	assert.Equal(t, errKey2, nil)
}

func TestGetValue(t *testing.T) {
	// Assume
	PersistentKeyValueStore, _ := NewPersistentKeyValueStore("./BookmarkStore.sqlite", "KeyValueTable")
	PersistentKeyValueStore.InitTable()
	PersistentKeyValueStore.Clear()
	PersistentKeyValueStore.SetValue("key1", "value1")
	PersistentKeyValueStore.SetValue("key2", "value2")

	// Action
	valKey1, errKey1 := PersistentKeyValueStore.GetValue("key1")
	valKey2, errKey2 := PersistentKeyValueStore.GetValue("key2")

	// Test
	assert.Equal(t, valKey1, "value1")
	assert.Equal(t, errKey1, nil)
	assert.Equal(t, valKey2, "value2")
	assert.Equal(t, errKey2, nil)
}

func TestClear(t *testing.T) {
	// Assume
	PersistentKeyValueStore, _ := NewPersistentKeyValueStore("./BookmarkStore.sqlite", "KeyValueTable")
	PersistentKeyValueStore.InitTable()
	PersistentKeyValueStore.SetValue("key1", "value1")
	PersistentKeyValueStore.SetValue("key2", "value2")

	// Action
	PersistentKeyValueStore.Clear()

	// Action
	valKey1, errKey1 := PersistentKeyValueStore.GetValue("key1")
	valKey2, errKey2 := PersistentKeyValueStore.GetValue("key2")

	assert.Equal(t, valKey1, "")
	assert.Assert(t, errKey1 != nil)
	assert.Equal(t, valKey2, "")
	assert.Assert(t, errKey2 != nil)

}


