package workflows

import (
	"github.com/microsoft/BladeMonRT/nodes"
)

/** Interface for defining execution sequence of nodes. */
type InterfaceWorkflow interface {
	AddNode(node nodes.InterfaceNode)
	Run(interfaceWorkflow InterfaceWorkflow)
	runVirt()
	GetNodes() []nodes.InterfaceNode
}

/** Concrete type for defining execution sequence of nodes. */
type Workflow struct {
}

func (workflow *Workflow) Run(interfaceWorkflow InterfaceWorkflow) {
	// TODO: add logging
	interfaceWorkflow.runVirt()
}