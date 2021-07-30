package dummy_node_b

import (
	"github.com/microsoft/BladeMonRT/nodes"
)

/** Node that has a hard-coded value for its result. */
type DummyNodeB struct {
	nodes.Node
}

func NewDummyNodeB() *DummyNodeB {
	// No fields to initialize.
	return &DummyNodeB{}
}

func (dummyNode *DummyNodeB) ProcessVirt(workflowContext *nodes.WorkflowContext) {
	dummyNode.SaveResult(dummyNode, workflowContext, "dummy-node-b-result")
  }