package azpubsub

import (
	"gotest.tools/assert"
	"testing"
	"strings"
)

func TestNewAzPubSubSimpleClient(t *testing.T) {
	// Action
	client := NewAzPubSubSimpleClient(true, "127.0.0.1")

	// Assert
	assert.Assert(t, client.hclient != HCLIENT(0))
	assert.Assert(t, client.hconfig != HCONFIG(0))
	assert.Assert(t, client.hproducer != HPRODUCER(0)) 

}

func TestSendMessage(t *testing.T) {
	// Assume
	client := NewAzPubSubSimpleClient(true, "127.0.0.1")

	// Action
	response, err := client.SendMessage("AzureCompute.Anvil.Request","test_message_1")

	// Assert
	assert.Equal(t, err, nil)
	assert.Equal(t, strings.Contains(response.message, "The operation timed out"), true)
}