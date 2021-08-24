package dummy_node_a

import (
	"github.com/microsoft/BladeMonRT/logging"
	"github.com/microsoft/BladeMonRT/nodes"
	"log"
)

/** Node that has the concatenation of its predecessors' results and a hard-coded value for its result. */
type DummyNodeA struct {
	nodes.Node
	Result string
}

func NewDummyNodeA() *DummyNodeA {
	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("DummyNodeA")
	var dummyNode DummyNodeA = DummyNodeA{Node: nodes.Node{Logger: logger}}
	dummyNode.Result = "dummy-node-result"
	return &dummyNode
}

func (dummyNode *DummyNodeA) ProcessVirt(workflowContext *nodes.WorkflowContext) error {
	dummyNode.Logger.Println("Running ProcessVirt method.")
	var result string

	// Add the predecessor results.
	for _, predecessorNodeResult := range dummyNode.GetPredecessorNodesResults(dummyNode, workflowContext) {
		result += predecessorNodeResult.(string) + "|"
	}

	// Add the result at the current nodes.
	result += dummyNode.Result

	dummyNode.SaveResult(dummyNode, workflowContext, result)
	return nil
}
