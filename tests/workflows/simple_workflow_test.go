package workflow_tests

import (
	"example.com/nodes"
	"example.com/workflows"
	"testing"
	"gotest.tools/assert"
)

func TestWorkflow(t *testing.T) {
	dummyNodeA := nodes.DummyNode{Node: nodes.Node{Name : "dummyNodeA"}}
	dummyNodeB := nodes.DummyNode{Node: nodes.Node{Name : "dummyNodeB"}}
	dummyNodeC := nodes.DummyNode{Node: nodes.Node{Name : "dummyNodeC"}}

	nodeToResult := make(map[string]string)

	workflow := workflows.SimpleWorkflow{Workflow: workflows.Workflow{}}
	workflow.AddNode(&dummyNodeA)
	workflow.AddNode(&dummyNodeB)
	workflow.AddNode(&dummyNodeC)

	workflow.RunVirt(nodeToResult)

	resultA, okA := nodeToResult[dummyNodeA.Name]
	assert.Equal(t, okA, true)
	assert.Equal(t, resultA, "dummy-node-result");

	resultB, okB := nodeToResult[dummyNodeB.Name]
	assert.Equal(t, okB, true)
	assert.Equal(t, resultB, "dummy-node-result");

	resultC, okC := nodeToResult[dummyNodeC.Name]
	assert.Equal(t, okC, true)
	assert.Equal(t, resultC, "dummy-node-result");
}
