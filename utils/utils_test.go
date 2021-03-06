package utils

import (
	"github.com/microsoft/BladeMonRT/test_configs"
	"gotest.tools/assert"
	"testing"
	"time"
)

func TestParseEventXML(t *testing.T) {
	// Assume
	expectedProvider := "CpuSpeedMonitoring"
	expectedEventId := 999
	expectedTimeCreated := time.Date(2021, 8, 10, 19, 10, 29, 0, time.UTC)
	expectedEventRecordId := 19818

	// Action
	event := NewUtils().ParseEventXML(test_configs.ARBITRARY_EVT_XML)

	// Assert
	assert.Equal(t, event.Provider, expectedProvider)
	assert.Equal(t, event.EventID, expectedEventId)
	assert.Equal(t, event.TimeCreated, expectedTimeCreated)
	assert.Equal(t, event.EventRecordID, expectedEventRecordId)
}
