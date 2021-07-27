package nodes

/** Node that has a hard-coded value for its result. */
type DummyNodeB struct {
	Node
}

func NewDummyNodeB() *DummyNodeB {
	// No fields to initialize.
	return &DummyNodeB{}
}

func (dummyNode *DummyNodeB) processVirt(workflowContext *WorkflowContext) {
	dummyNode.SaveResult(dummyNode, workflowContext, "dummy-node-b-result")
  }