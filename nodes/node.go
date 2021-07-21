package nodes

import (
	"fmt"
)

// Interface for defining unit of work to be processed by event loop
type InterfaceNode interface {
 	ProcessVirt(predecessorNodeResults []interface{})
	GetResult() interface{}
	SaveResult(result interface{}) 
}
 
// Concrete type for defining unit of work to be processed by event loop
type Node struct {
	nodeResult interface{}
	Name string
}

func (node *Node) SaveResult(result interface{}) {
	fmt.Println("Running SaveResult method.")
	node.nodeResult = result
}

func (node *Node) GetResult() interface{} {
	return node.nodeResult
}