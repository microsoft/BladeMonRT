package main

import (
	"gotest.tools/assert"
	"testing"
	"strings"
	"github.com/microsoft/BladeMonRT/AzPubSub"
)

func TestNewAzPubSubSimpleClient(t *testing.T) {
	// Action
	client := azpubsub.NewAzPubSubSimpleClient(true, "127.0.0.1")

	// Assert
	assert.Assert(t, client.Hclient != azpubsub.HCLIENT(0))
	assert.Assert(t, client.Hconfig != azpubsub.HCONFIG(0))
	assert.Assert(t, client.Hproducer != azpubsub.HPRODUCER(0)) 

}

func TestSendMessage(t *testing.T) {
	// Assume
	client := azpubsub.NewAzPubSubSimpleClient(true, "127.0.0.1")

	// Action
	response, err := client.SendMessage("AzureCompute.Anvil.Request","test_message_1")

	// Assert
	assert.Equal(t, err, nil)
	assert.Equal(t, strings.Contains(response.Message, "The operation timed out") || strings.Contains(response.Message, "The server name or address could not be resolved"), true)
}