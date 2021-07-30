
package dummy_node_c

import (
	"github.com/microsoft/BladeMonRT/nodes"
)

/** Node that has a hard-coded value for its result. */
type DummyNodeC struct {
	nodes.Node
	result string
}

func NewDummyNodeC(resultProvider InterfaceResultProvider) *DummyNodeC {
	var dummyNode DummyNodeC = DummyNodeC{}
	dummyNode.result = resultProvider.result()
	return &dummyNode
}

func (dummyNode *DummyNodeC) ProcessVirt(workflowContext *nodes.WorkflowContext) {
  dummyNode.SaveResult(dummyNode, workflowContext, dummyNode.result)
}

/** Interface that provides the result for the dummy node. */
type InterfaceResultProvider interface {
	result() string
}

/** Class that provides the result for the dummy node in production. */
type ResultProvider struct {
}

func (dummyNode *ResultProvider) result() string {
	return "dummy-result-c"
}
