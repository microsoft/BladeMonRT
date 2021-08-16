package dummy_node_b

import (
	"github.com/microsoft/BladeMonRT/nodes"
	"github.com/microsoft/BladeMonRT/logging"
	"log"
)

/** Interface for a node that has a hard-coded value for its result. */
type InterfaceDummyNodeB interface {
	Process(interfaceNode nodes.InterfaceNode, workflowContext *nodes.WorkflowContext) error
	ProcessVirt(workflowContext *nodes.WorkflowContext) error
	SaveResult(interfaceNode nodes.InterfaceNode, workflowContext *nodes.WorkflowContext, result interface{})
	GetResult(interfaceNode nodes.InterfaceNode, workflowContext *nodes.WorkflowContext) interface{}
	GetPredecessorNodes(interfaceNode nodes.InterfaceNode, workflowContext *nodes.WorkflowContext) []nodes.InterfaceNode
	GetPredecessorNodesResults(interfaceNode nodes.InterfaceNode, workflowContext *nodes.WorkflowContext) []interface{}
	result() string
}

/** Node that has a hard-coded value for its result. */
type DummyNodeB struct {
	nodes.Node
}

func NewDummyNodeB() *DummyNodeB {
	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("DummyNodeB")
	var dummyNode DummyNodeB = DummyNodeB{Node : nodes.Node{Logger : logger}}
	return &dummyNode
}

func (dummyNode *DummyNodeB) result() string {
	return "dummy-node-b-result"
}

func (dummyNode *DummyNodeB) ProcessVirt(workflowContext *nodes.WorkflowContext) error {
	dummyNode.Logger.Println("Running ProcessVirt method.")
	dummyNode.SaveResult(dummyNode, workflowContext, dummyNode.result())
	return nil
}