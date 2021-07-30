package workflows

import (
	"github.com/microsoft/BladeMonRT/nodes"
	"log"
)

/** Interface for defining execution sequence of nodes. */
type InterfaceWorkflow interface {
	AddNode(node nodes.InterfaceNode)
	Run(interfaceWorkflow InterfaceWorkflow, workflowContext *nodes.WorkflowContext)
	runVirt(workflowContext *nodes.WorkflowContext)
	GetNodes() []nodes.InterfaceNode
}

/** Concrete type for defining execution sequence of nodes. */
type Workflow struct {
	workflowContext *nodes.WorkflowContext
	Logger *log.Logger
}

func (workflow *Workflow) Run(interfaceWorkflow InterfaceWorkflow, workflowContext *nodes.WorkflowContext) {
	workflow.Logger.Println("Running run method.")

	// Set the nodes in the workflow context to the nodes in this workflow.
	workflowContext.SetNodes(interfaceWorkflow.GetNodes())

	interfaceWorkflow.runVirt(workflowContext)
}
