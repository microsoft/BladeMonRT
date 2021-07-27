package nodes

import (
	"fmt"
)

type WorkflowContext struct {
	PredecessorNodesResults []interface{}
	NodeToResult map[InterfaceNode]interface{}
}

// Interface for defining unit of work to be processed by event loop. Classes that implement InterfaceNode can provide their own constructor.
type InterfaceNode interface {
 	Process(interfaceNode InterfaceNode, workflowContext *WorkflowContext)
	processVirt(workflowContext *WorkflowContext)
	SaveResult(interfaceNode InterfaceNode, workflowContext *WorkflowContext, result interface{})
	GetResult(interfaceNode InterfaceNode, workflowContext *WorkflowContext) interface{}
}
 
// Concrete type for defining unit of work to be processed by event loop
type Node struct {
	dummyVariable interface{} // This variable makes the Node struct non-empty, to prevent GO's behavior in the allocation of zero-sized objects. 
}

func (node *Node) Process(interfaceNode InterfaceNode, workflowContext *WorkflowContext) {
	// TODO: Add logging.
	interfaceNode.processVirt(workflowContext)
}

func (node *Node) SaveResult(interfaceNode InterfaceNode, workflowContext *WorkflowContext, result interface{}) {
	fmt.Println("Running SaveResult method.")
	workflowContext.NodeToResult[interfaceNode] = result
	workflowContext.PredecessorNodesResults = append(workflowContext.PredecessorNodesResults, result)
	fmt.Println(workflowContext.NodeToResult)
}

func (node *Node) GetResult(interfaceNode InterfaceNode, workflowContext *WorkflowContext) interface{} {
	return workflowContext.NodeToResult[interfaceNode]
}