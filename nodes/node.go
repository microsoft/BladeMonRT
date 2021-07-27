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
	SaveResult(interfaceNode InterfaceNode, workflowContext *WorkflowContext, result interface{})
	GetResult(interfaceNode InterfaceNode, workflowContext *WorkflowContext) interface{}
	GetPredecessorResults(interfaceNode InterfaceNode, workflowContext *WorkflowContext) []interface{}
}
 
// Concrete type for defining unit of work to be processed by event loop.
type Node struct {
	dummyVariable interface{} // This variable makes the Node struct non-empty, to prevent GO's behavior in the allocation of zero-sized objects. 
}

func (node *Node) Process(interfaceNode InterfaceNode, workflowContext *WorkflowContext) {
	// TODO: Add logging.
	interfaceNode.processVirt(workflowContext)
}

func (node *Node) SaveResult(interfaceNode InterfaceNode, workflowContext *WorkflowContext, result interface{}) {
	fmt.Println("Running SaveResult method.")
	workflowContext.nodeToResult[interfaceNode] = result
}

func (node *Node) GetResult(interfaceNode InterfaceNode, workflowContext *WorkflowContext) interface{} {
	fmt.Println("Running GetResult method.")
	return workflowContext.nodeToResult[interfaceNode]
}

func (node *Node) GetPredecessorResults(interfaceNode InterfaceNode, workflowContext *WorkflowContext) []interface{} {
	var predecessorNodeResults []interface{}
	for _, predecessorNode := range workflowContext.GetNodes() {
		if (interfaceNode == predecessorNode) {
			break
		}
		predecessorNodeResults = append(predecessorNodeResults,  workflowContext.nodeToResult[predecessorNode])
	}
	return predecessorNodeResults
}