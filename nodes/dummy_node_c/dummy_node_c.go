package dummy_node_c

import (
	"log"

	"github.com/microsoft/BladeMonRT/logging"
	"github.com/microsoft/BladeMonRT/nodes"
)

/** Node that has a hard-coded value for its result. */
type DummyNodeC struct {
	nodes.Node
	result         string
	resultProvider InterfaceResultProvider
}

/** Interface that provides the result for the dummy node. */
type InterfaceResultProvider interface {
	result() string
}

/** Class that provides the result for the dummy node in production. */
type ResultProvider struct {
}

func NewDummyNodeC() *DummyNodeC {
	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("DummyNodeC")
	var dummyNode DummyNodeC = DummyNodeC{Node: nodes.Node{Logger: logger}, resultProvider: ResultProvider{}}
	return &dummyNode
}

func (dummyNode *DummyNodeC) ProcessVirt(workflowContext *nodes.WorkflowContext) {
	dummyNode.Logger.Println("Running ProcessVirt method.")
	dummyNode.SaveResult(dummyNode, workflowContext, dummyNode.resultProvider.result())
}

func (dummyNode ResultProvider) result() string {
	return "dummy-result-c"
}
