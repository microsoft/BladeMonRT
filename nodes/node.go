package nodes

import (
	"fmt"
)

// Interface for defining unit of work to be processed by event loop
type InterfaceNode interface {
	GetName() string
	SetName(name string)
	InitializeFields()
 	ProcessVirt(predecessorNodeResults []interface{})
	GetResult() interface{}
	SaveResult(result interface{}) 
}
 
// Concrete type for defining unit of work to be processed by event loop
type Node struct {
	nodeResult interface{}
	name string
}

func (node *Node) SaveResult(result interface{}) {
	fmt.Println("Running SaveResult method.")
	node.nodeResult = result
}

func (node *Node) GetResult() interface{} {
	return node.nodeResult
}

func (node *Node) SetName(name string) {
	node.name = name
}
func (node *Node) GetName() string {
	return node.name
}
