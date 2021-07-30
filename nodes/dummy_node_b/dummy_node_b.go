package dummy_node_b

import (
	"github.com/microsoft/BladeMonRT/nodes"
)

/** Interface for a node that has a hard-coded value for its result. */
type InterfaceDummyNodeB interface {
	Process(interfaceNode nodes.InterfaceNode, workflowContext *nodes.WorkflowContext)
	ProcessVirt(workflowContext *nodes.WorkflowContext)
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
	// No fields to initialize.
	return &DummyNodeB{}
}

func (dummyNode *DummyNodeB) result() string {
	return "dummy-node-b-result"
}

func (dummyNode *DummyNodeB) ProcessVirt(workflowContext *nodes.WorkflowContext) {
	dummyNode.SaveResult(dummyNode, workflowContext, dummyNode.result())
}