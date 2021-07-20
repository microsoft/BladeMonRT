package workflows

import (
	"example.com/nodes"
)

type InterfaceWorkflow interface {
 	RunVirt(workflowContextResult map[string]string)
	AddNode(node nodes.InterfaceNode)
}

type Workflow struct {
	Nodes []nodes.InterfaceNode
}
