
package nodes

/** Node that has a hard-coded value for its result. */
type DummyNodeC struct {
	Node
	result string
}

func NewDummyNodeC() *DummyNodeC {
	var dummyNode DummyNodeC = DummyNodeC{}
	dummyNode.result = "dummy-result-c"
	return &dummyNode
}

func (dummyNode *DummyNodeC) processVirt(workflowContext *WorkflowContext) {
  dummyNode.saveResult(dummyNode, workflowContext, dummyNode.result)
}