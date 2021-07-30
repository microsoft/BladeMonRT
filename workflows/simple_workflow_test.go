package workflows

import (
	"github.com/microsoft/BladeMonRT/nodes"
	"github.com/microsoft/BladeMonRT/nodes/dummy_node_a"
	"testing"
	"gotest.tools/assert"
)

func TestWorkflow(t *testing.T) {
	var dummyNodeA nodes.InterfaceNode = dummy_node_a.NewDummyNodeA()
	var dummyNodeB nodes.InterfaceNode = dummy_node_a.NewDummyNodeA()
	var dummyNodeC nodes.InterfaceNode = dummy_node_a.NewDummyNodeA()

	var workflow *SimpleWorkflow = NewSimpleWorkflow()
	workflow.AddNode(dummyNodeA)
	workflow.AddNode(dummyNodeB)
	workflow.AddNode(dummyNodeC)

	var workflowContext *nodes.WorkflowContext = nodes.NewWorkflowContext()
	workflow.Run(workflow, workflowContext)

	// Check that the result at each node includes the predecessor results and the expected hard-coded value.
	resultA := dummyNodeA.GetResult(dummyNodeA, workflowContext)
	assert.Equal(t, resultA, "dummy-node-result");
	resultB := dummyNodeB.GetResult(dummyNodeB, workflowContext)
	assert.Equal(t, resultB, "dummy-node-result|dummy-node-result");
	resultC := dummyNodeC.GetResult(dummyNodeC, workflowContext)
	assert.Equal(t, resultC, "dummy-node-result|dummy-node-result|dummy-node-result|dummy-node-result");
}
