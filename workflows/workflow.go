package workflows

import (
	"github.com/microsoft/BladeMonRT/nodes"
)

/** Interface for defining execution sequence of nodes. */
type InterfaceWorkflow interface {
	AddNode(node nodes.InterfaceNode)
	Run(interfaceWorkflow InterfaceWorkflow)
	runVirt(workflowContext *nodes.WorkflowContext)
	GetNodes() []nodes.InterfaceNode
	GetWorkflowContext() *nodes.WorkflowContext
}

/** Concrete type for defining execution sequence of nodes. */
type Workflow struct {
	workflowContext *nodes.WorkflowContext
}

func (workflow *Workflow) Run(interfaceWorkflow InterfaceWorkflow) {
	// TODO: add logging

	// Create a workflow context using the nodes in this workflow.
	var workflowNodes []nodes.InterfaceNode = interfaceWorkflow.GetNodes()
	var workflowContext *nodes.WorkflowContext = nodes.NewWorkflowContext(workflowNodes)
	workflow.workflowContext = workflowContext

	interfaceWorkflow.runVirt(workflowContext)
}

func (workflow *Workflow) GetWorkflowContext() *nodes.WorkflowContext {
	return workflow.workflowContext
}
