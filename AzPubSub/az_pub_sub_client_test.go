package main

import (
	"gotest.tools/assert"
	"testing"
	"fmt"
)

func TestNewAzPubSubSimpleClient(t *testing.T) {
	// Action
	client := NewAzPubSubSimpleClient(true, "127.0.0.7")

	// Assert
	assert.Assert(t, client.hclient != HCLIENT(0))
	assert.Assert(t, client.hconfig != HCONFIG(0))
	assert.Assert(t, client.hproducer != HPRODUCER(0)) 

}

func TestSendMessage(t *testing.T) {
	// Assume
	client := NewAzPubSubSimpleClient(true, "127.0.0.7")

	// Action
	response := client.sendMessage("AzureCompute.Anvil.Request","test_message_1")

	// Assert
	assert.Equal(t, response.message, "The operation timed out")
	fmt.Println(response)
}
  
/*

    self.assertTrue(self.simple_client.send_message("AzureCompute.Anvil.Request","test_message_1"),
                     "Failed to send message to AzPubSub EventServer")

    self.assertTrue(self.simple_client.send_message("AzureCompute.Anvil.Request","test_message_2", "2375f39f-2147-4245-b4b7-71db290bc194"),
                     "Failed to send message with key to AzPubSub EventServer")
					 */