package nodes

import (
	"fmt"
)

/** Class that stores information about the current state of a running workflow. */
type WorkflowContext struct {
	nodes []InterfaceNode
	nodeToResult map[InterfaceNode]interface{}
}

func NewWorkflowContext(nodes []InterfaceNode) *WorkflowContext {
	var nodeToResult map[InterfaceNode]interface{} = make(map[InterfaceNode]interface{})
	var workflowContext WorkflowContext = WorkflowContext{nodeToResult : nodeToResult, nodes : nodes}
	return &workflowContext
}

func (workflowContext *WorkflowContext) GetNodes() []InterfaceNode {
	return workflowContext.nodes
}

// Interface for defining unit of work to be processed by event loop. Classes that implement InterfaceNode can provide their own constructor.
type InterfaceNode interface {
 	Process(interfaceNode InterfaceNode, workflowContext *WorkflowContext)
	processVirt(workflowContext *WorkflowContext)
	saveResult(interfaceNode InterfaceNode, workflowContext *WorkflowContext, result interface{})
	GetResult(interfaceNode InterfaceNode, workflowContext *WorkflowContext) interface{}
	getPredecessorNodes(interfaceNode InterfaceNode, workflowContext *WorkflowContext) []InterfaceNode
	getPredecessorNodesResults(interfaceNode InterfaceNode, workflowContext *WorkflowContext) []interface{}
}
 
// Concrete type for defining unit of work to be processed by event loop.
type Node struct {
	dummyVariable interface{} // This variable makes the Node struct non-empty, to prevent GO's behavior in the allocation of zero-sized objects. 
}

func (node *Node) Process(interfaceNode InterfaceNode, workflowContext *WorkflowContext) {
	// TODO: Add logging.
	interfaceNode.processVirt(workflowContext)
}

func (node *Node) saveResult(interfaceNode InterfaceNode, workflowContext *WorkflowContext, result interface{}) {
	fmt.Println("Running SaveResult method.")
	workflowContext.nodeToResult[interfaceNode] = result
}

func (node *Node) GetResult(interfaceNode InterfaceNode, workflowContext *WorkflowContext) interface{} {
	fmt.Println("Running GetResult method.")
	return workflowContext.nodeToResult[interfaceNode]
}

func (node *Node) getPredecessorNodes(interfaceNode InterfaceNode, workflowContext *WorkflowContext) []InterfaceNode {
	var predecessorNodes []InterfaceNode
	var nodes []InterfaceNode = workflowContext.GetNodes()
	for _, currNode := range nodes {
		if (interfaceNode == currNode) {
			break
		}
		predecessorNodes = append(predecessorNodes, currNode)
	}
	return predecessorNodes
}

func (node *Node) getPredecessorNodesResults(interfaceNode InterfaceNode, workflowContext *WorkflowContext) []interface{} {
	var predecessorNodesResults []interface{}
	var predecessorNodes []InterfaceNode = interfaceNode.getPredecessorNodes(interfaceNode, workflowContext)
	for _, predecessorNode := range predecessorNodes {
		predecessorNodesResults = append(predecessorNodesResults, workflowContext.nodeToResult[predecessorNode])
	}
	return predecessorNodesResults
}