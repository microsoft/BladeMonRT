package dummy_node_with_panic

import (
	"log"

	"github.com/microsoft/BladeMonRT/logging"
	"github.com/microsoft/BladeMonRT/nodes"
)

/** Node that throws two panics in its ProcessVirt method. */
type DummyNodeWithPanic struct {
	nodes.Node
}

func NewDummyNodeWithPanic() *DummyNodeWithPanic {
	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("DummyNodeWithPanic")
	var dummyNode DummyNodeWithPanic = DummyNodeWithPanic{Node : nodes.Node{Logger : logger}}
	return &dummyNode
}

func (dummyNode *DummyNodeWithPanic) ProcessVirt(workflowContext *nodes.WorkflowContext) error {
	panic("Throwing a first panic for testing.")
	panic("Throwing a second panic for testing.")
	return nil
}