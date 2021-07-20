package nodes

import (
	"fmt"
)

// Interface for defining unit of work to be processed by event loop
type InterfaceNode interface {
 	ProcessVirt(workflowContextResult map[string]string)
}
 
// Concrete type for defining unit of work to be processed by event loop
type Node struct {
	NodeResult string
	Name string
}

func (node *Node) ProcessVirt(workflowContextResult map[string]string) {
	fmt.Println("Running processVirt method.")
}

func (node *Node) SaveResult(workflowContextResult map[string]string) {
	fmt.Println("Running SaveResult method.")
	fmt.Println(node.NodeResult)
	workflowContextResult[node.Name] = node.NodeResult
}