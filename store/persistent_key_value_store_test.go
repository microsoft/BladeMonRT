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

func TestGetValue_KeyExists(t *testing.T) {
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

func TestGetValue_KeyDoesNotExist(t *testing.T) {
	// Assume
	PersistentKeyValueStore, _ := NewPersistentKeyValueStore("./BookmarkStore.sqlite", "KeyValueTable")
	PersistentKeyValueStore.InitTable()
	PersistentKeyValueStore.Clear()
	PersistentKeyValueStore.SetValue("key1", "value1")

	// Action
	val, err := PersistentKeyValueStore.GetValue("key2")

	// Test
	assert.Equal(t, val, "")
	assert.DeepEqual(t, err.Error(), "Key=key2 not found in the store.")
}

func TestClear(t *testing.T) {
	// Assume
	PersistentKeyValueStore, _ := NewPersistentKeyValueStore("./BookmarkStore.sqlite", "KeyValueTable")
	PersistentKeyValueStore.InitTable()
	PersistentKeyValueStore.SetValue("key1", "value1")

	// Action
	PersistentKeyValueStore.Clear()

	// Action
	valKey1, errKey1 := PersistentKeyValueStore.GetValue("key1")

	assert.Equal(t, valKey1, "")
	assert.Assert(t, errKey1 != nil)
}
