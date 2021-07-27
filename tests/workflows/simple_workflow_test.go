package workflow_tests

import (
	"github.com/microsoft/BladeMonRT/nodes"
	"github.com/microsoft/BladeMonRT/workflows"
	"testing"
	"gotest.tools/assert"
)

func TestWorkflow(t *testing.T) {
	var dummyNodeA nodes.InterfaceNode = nodes.NewDummyNode()
	var dummyNodeB nodes.InterfaceNode = nodes.NewDummyNode()
	var dummyNodeC nodes.InterfaceNode = nodes.NewDummyNode()

	var workflow *workflows.SimpleWorkflow = workflows.NewSimpleWorkflow()
	workflow.AddNode(dummyNodeA)
	workflow.AddNode(dummyNodeB)
	workflow.AddNode(dummyNodeC)

	workflow.Run(workflow)

	// Check that the result at each node includes the predecessor results and the expected hard-coded value.
	var workflowContext *nodes.WorkflowContext = workflow.GetWorkflowContext()
	resultA := dummyNodeA.GetResult(dummyNodeA, workflowContext)
	assert.Equal(t, resultA, "dummy-node-result");
	resultB := dummyNodeB.GetResult(dummyNodeB, workflowContext)
	assert.Equal(t, resultB, "dummy-node-result|dummy-node-result");
	resultC := dummyNodeC.GetResult(dummyNodeC, workflowContext)
	assert.Equal(t, resultC, "dummy-node-result|dummy-node-result|dummy-node-result|dummy-node-result");
}
