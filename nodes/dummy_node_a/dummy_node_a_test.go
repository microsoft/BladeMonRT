package dummy_node_a

import (
	"github.com/microsoft/BladeMonRT/nodes"
	"gotest.tools/assert"
	"testing"
)

func TestDummyNodeA(t *testing.T) {
	var dummyNode *DummyNodeA = NewDummyNodeA()
	var workflowContext *nodes.WorkflowContext = nodes.NewWorkflowContext()
	dummyNode.Process(dummyNode, workflowContext)
	result := dummyNode.GetResult(dummyNode, workflowContext)

	assert.Equal(t, result, "dummy-node-result")
}
