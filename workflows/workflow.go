package workflows

import (
	"github.com/microsoft/BladeMonRT/nodes"
)

/** Interface for defining execution sequence of nodes. */
type InterfaceWorkflow interface {
	AddNode(node nodes.InterfaceNode)
	RunVirt()
}

/** Concrete type for defining execution sequence of nodes. */
type Workflow struct {
}
