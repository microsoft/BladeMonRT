package workflows

import (
	"example.com/nodes"
)

// Interface for defining execution sequence of nodes
type InterfaceWorkflow interface {
 	RunVirt(workflowContextResult map[string]string)
	AddNode(node nodes.InterfaceNode)
}

// Concrete type for defining execution sequence of nodes
type Workflow struct {
	Nodes []nodes.InterfaceNode
}
