package store

import (
	"gotest.tools/assert"
	"testing"
)

func TestInitTable(t *testing.T) {
	// Action
	configStore, err := NewConfigStore("./BookmarkStore.sqlite", "ConfigTable")
	configStore.InitTable()

	// Assert
	assert.Equal(t, err, nil)
}

func TestSetConfigValue(t *testing.T) {
	// Assume
	configStore, _ := NewConfigStore("./BookmarkStore.sqlite", "ConfigTable")
	configStore.InitTable()

	// Action
	errKey1 := configStore.SetConfigValue("key1", "value1")
	errKey2 := configStore.SetConfigValue("key2", "value2")

	assert.Equal(t, errKey1, nil)
	assert.Equal(t, errKey2, nil)
}

func TestGetConfigValue(t *testing.T) {
	// Assume
	configStore, _ := NewConfigStore("./BookmarkStore.sqlite", "ConfigTable")
	configStore.InitTable()
	configStore.SetConfigValue("key1", "value1")
	configStore.SetConfigValue("key2", "value2")

	// Action
	valKey1, errKey1 := configStore.GetConfigValue("key1")
	valKey2, errKey2 := configStore.GetConfigValue("key2")

	// Test
	assert.Equal(t, valKey1, "value1")
	assert.Equal(t, errKey1, nil)
	assert.Equal(t, valKey2, "value2")
	assert.Equal(t, errKey2, nil)
}
