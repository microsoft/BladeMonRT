package node_tests

import (
	"example.com/nodes"
	"testing"
	"gotest.tools/assert"
)

func TestDummyNode(t *testing.T) {
	dummyNode := nodes.DummyNode{Node: nodes.Node{Name : "dummyNode"}}
	nodeToResult := make(map[string]string)
	dummyNode.ProcessVirt(nodeToResult)
	result, ok := nodeToResult[dummyNode.Name]

	assert.Equal(t, ok, true)
	assert.Equal(t, result, "dummy-node-result");
}
