
package dummy_node_c

import (
	"github.com/microsoft/BladeMonRT/nodes"
)

/** Node that has a hard-coded value for its result. */
type DummyNodeC struct {
	nodes.Node
	result string
}

func NewDummyNodeC() *DummyNodeC {
	var dummyNode DummyNodeC = DummyNodeC{}
	dummyNode.result = "dummy-result-c"
	return &dummyNode
}

func (dummyNode *DummyNodeC) ProcessVirt(workflowContext *nodes.WorkflowContext) {
  dummyNode.SaveResult(dummyNode, workflowContext, dummyNode.result)
}