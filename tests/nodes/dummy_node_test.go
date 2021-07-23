package node_tests

import (
	"github.com/microsoft/BladeMonRT/nodes"
	"testing"
	"gotest.tools/assert"
)

func TestDummyNode(t *testing.T) {
	var dummyNode nodes.DummyNode = nodes.DummyNode{Node: nodes.Node{}}

	var predecessorNodesResults []interface{}
	dummyNode.ProcessVirt(predecessorNodesResults)
	result := dummyNode.GetResult()

	assert.Equal(t, result, "dummy-node-result");
}
