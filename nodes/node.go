package nodes

import (
	"fmt"
)

// Interface for defining unit of work to be processed by event loop
type InterfaceNode interface {
	InitializeFields()
 	Process(interfaceNode InterfaceNode, predecessorNodesResults []interface{})
	processVirt(predecessorNodesResults []interface{})
	GetResult() interface{}
	SaveResult(result interface{}) 
}
 
// Concrete type for defining unit of work to be processed by event loop
type Node struct {
	nodeResult interface{}
}

func (node *Node) Process(interfaceNode InterfaceNode, predecessorNodesResults []interface{}) {
	// TODO: Add logging.
	interfaceNode.processVirt(predecessorNodesResults)
}

func (node *Node) SaveResult(result interface{}) {
	fmt.Println("Running SaveResult method.")
	node.nodeResult = result
}

func (node *Node) GetResult() interface{} {
	return node.nodeResult
}