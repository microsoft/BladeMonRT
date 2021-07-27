package node_tests

import (
	"github.com/microsoft/BladeMonRT/nodes"
	"testing"
	"gotest.tools/assert"
)

func TestDummyNode(t *testing.T) {
	var dummyNode nodes.DummyNode = nodes.DummyNode{Node: nodes.Node{}}
	var workflowNodes []nodes.InterfaceNode
	workflowNodes = append(workflowNodes, &dummyNode)

	var workflowContext *nodes.WorkflowContext = nodes.NewWorkflowContext(workflowNodes)
	dummyNode.Process(&dummyNode, workflowContext)
	result := dummyNode.GetResult(&dummyNode, workflowContext)

	assert.Equal(t, result, "dummy-node-result");
}
