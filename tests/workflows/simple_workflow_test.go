package workflow_tests

import (
	"github.com/microsoft/BladeMonRT/nodes"
	"github.com/microsoft/BladeMonRT/workflows"
	"testing"
	"gotest.tools/assert"
)

func TestWorkflow(t *testing.T) {
	var dummyNodeA nodes.DummyNode = nodes.DummyNode{Node: nodes.Node{}}
	var dummyNodeB nodes.DummyNode = nodes.DummyNode{Node: nodes.Node{}}
	var dummyNodeC nodes.DummyNode = nodes.DummyNode{Node: nodes.Node{}}

	var workflow workflows.SimpleWorkflow = workflows.SimpleWorkflow{Workflow: workflows.Workflow{}}
	workflow.AddNode(&dummyNodeA)
	workflow.AddNode(&dummyNodeB)
	workflow.AddNode(&dummyNodeC)

	workflow.RunVirt()

	// Check that the result at each node includes the predecessor results and the expected hard-coded value.
	resultA := dummyNodeA.GetResult()
	assert.Equal(t, resultA, "dummy-node-result");
	resultB := dummyNodeB.GetResult()
	assert.Equal(t, resultB, "dummy-node-result|dummy-node-result");
	resultC := dummyNodeC.GetResult()
	assert.Equal(t, resultC, "dummy-node-result|dummy-node-result|dummy-node-result|dummy-node-result");
}
