package azpubsub

import (
	"gotest.tools/assert"
	"testing"
	"github.com/microsoft/BladeMonRT/test_configs"
)

func TestFetchAzPubSubPfVIP(t *testing.T) {
	// Assume
	utils := NewUtils()

	// Action
	vip, err := utils.FetchAzPubSubPfVIP(test_configs.TEST_PF_AZ_PUB_SUB_VIP_FILE)

	// Assert
	assert.Equal(t, err, nil)
	assert.Equal(t, vip, "127.0.0.1")
}