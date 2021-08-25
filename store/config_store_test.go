package store

import (
	"gotest.tools/assert"
	"testing"
)

func TestInitTable(t *testing.T) {
	// Action
	PersistentKeyValueStore, err := NewPersistentKeyValueStore("./BookmarkStore.sqlite", "ConfigTable")
	PersistentKeyValueStore.InitTable()

	// Assert
	assert.Equal(t, err, nil)
}

func TestSetConfigValue(t *testing.T) {
	// Assume
	PersistentKeyValueStore, _ := NewPersistentKeyValueStore("./BookmarkStore.sqlite", "ConfigTable")
	PersistentKeyValueStore.InitTable()

	// Action
	errKey1 := PersistentKeyValueStore.SetConfigValue("key1", "value1")
	errKey2 := PersistentKeyValueStore.SetConfigValue("key2", "value2")

	assert.Equal(t, errKey1, nil)
	assert.Equal(t, errKey2, nil)
}

func TestGetConfigValue(t *testing.T) {
	// Assume
	PersistentKeyValueStore, _ := NewPersistentKeyValueStore("./BookmarkStore.sqlite", "ConfigTable")
	PersistentKeyValueStore.InitTable()
	PersistentKeyValueStore.SetConfigValue("key1", "value1")
	PersistentKeyValueStore.SetConfigValue("key2", "value2")

	// Action
	valKey1, errKey1 := PersistentKeyValueStore.GetConfigValue("key1")
	valKey2, errKey2 := PersistentKeyValueStore.GetConfigValue("key2")

	// Test
	assert.Equal(t, valKey1, "value1")
	assert.Equal(t, errKey1, nil)
	assert.Equal(t, valKey2, "value2")
	assert.Equal(t, errKey2, nil)
}
