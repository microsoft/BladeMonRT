package main

import (
	"gotest.tools/assert"
	"testing"
	"github.com/microsoft/BladeMonRT/AzPubSub"
)

func TestNewAzPubSubGlobalClient(t *testing.T) {
	// Action
	client := azpubsub.NewAzPubSubGlobalClient([]string{"AnvilRepairRequest"}, true, "10.0.0.155:9092")

	// Assert
	assert.Assert(t, client.Hclient != azpubsub.HCLIENT(0))
	assert.Assert(t, client.Hconfig != azpubsub.HCONFIG(0))
	// assert.Assert(t, client.Hproducer != azpubsub.HPRODUCER(0)) 

}