
package dummy_node_c

import (
	"github.com/microsoft/BladeMonRT/nodes"
	"github.com/microsoft/BladeMonRT/logging"
	"log"
)

/** Node that has a hard-coded value for its result. */
type DummyNodeC struct {
	nodes.Node
	result string
	resultProvider InterfaceResultProvider
}

func NewDummyNodeC(resultProvider InterfaceResultProvider) *DummyNodeC {
	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("DummyNodeB")
	var dummyNode DummyNodeC = DummyNodeC{Node : nodes.Node{Logger : logger}, resultProvider : resultProvider}
	return &dummyNode
}

func (dummyNode *DummyNodeC) ProcessVirt(workflowContext *nodes.WorkflowContext) {
  dummyNode.SaveResult(dummyNode, workflowContext, dummyNode.resultProvider.result())
}

/** Interface that provides the result for the dummy node. */
type InterfaceResultProvider interface {
	result() string
}

/** Class that provides the result for the dummy node in production. */
type ResultProvider struct {
}

func (dummyNode ResultProvider) result() string {
	return "dummy-result-c"
}
