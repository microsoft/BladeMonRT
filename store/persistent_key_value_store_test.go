package store

import (
	"gotest.tools/assert"
	"testing"
)

func TestSetValue(t *testing.T) {
	// Assume
	PersistentKeyValueStore, _ := NewPersistentKeyValueStore("./BookmarkStore.sqlite", "KeyValueTable")
	PersistentKeyValueStore.Clear()

	// Action
	errKey1 := PersistentKeyValueStore.SetValue("key1", "value1")
	errKey2 := PersistentKeyValueStore.SetValue("key2", "value2")

	// Assert
	assert.Equal(t, errKey1, nil)
	assert.Equal(t, errKey2, nil)
}

func TestGetValue_KeyExists(t *testing.T) {
	// Assume
	PersistentKeyValueStore, _ := NewPersistentKeyValueStore("./BookmarkStore.sqlite", "KeyValueTable")
	PersistentKeyValueStore.Clear()
	PersistentKeyValueStore.SetValue("key1", "value1")
	PersistentKeyValueStore.SetValue("key2", "value2")

	// Action
	valKey1, errKey1 := PersistentKeyValueStore.GetValue("key1")
	valKey2, errKey2 := PersistentKeyValueStore.GetValue("key2")

	// Assert
	assert.Equal(t, valKey1, "value1")
	assert.Equal(t, errKey1, nil)
	assert.Equal(t, valKey2, "value2")
	assert.Equal(t, errKey2, nil)
}

func TestGetValue_KeyDoesNotExist(t *testing.T) {
	// Assume
	PersistentKeyValueStore, _ := NewPersistentKeyValueStore("./BookmarkStore.sqlite", "KeyValueTable")
	PersistentKeyValueStore.Clear()
	PersistentKeyValueStore.SetValue("key1", "value1")

	// Action
	val, err := PersistentKeyValueStore.GetValue("key2")

	// Assert
	assert.Equal(t, val, "")
	assert.Equal(t, err, nil)
}

func TestClear(t *testing.T) {
	// Assume
	PersistentKeyValueStore, _ := NewPersistentKeyValueStore("./BookmarkStore.sqlite", "KeyValueTable")
	PersistentKeyValueStore.SetValue("key1", "value1")

	// Action
	PersistentKeyValueStore.Clear()

	// Assert
	valKey1, errKey1 := PersistentKeyValueStore.GetValue("key1")
	assert.Equal(t, valKey1, "")
	assert.Assert(t, errKey1, nil)
}
