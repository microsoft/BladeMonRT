package workflows

import (
	"errors"
	"fmt"
	"github.com/microsoft/BladeMonRT/nodes"
	"log"
)

// InterfaceWorkflow mock generation.
//go:generate mockgen -source=./workflow.go -destination=./mock_workflow.go -package=workflows

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
	Logger          *log.Logger
}

func (workflow *Workflow) Run(interfaceWorkflow InterfaceWorkflow, workflowContext *nodes.WorkflowContext) {
	workflow.Logger.Println("Running run method.")

	// Set the nodes in the workflow context to the nodes in this workflow.
	workflowContext.SetNodes(interfaceWorkflow.GetNodes())

	interfaceWorkflow.runVirt(workflowContext)
}

/** Processes a single node. Classes that implement InterfaceWorkflow should call this to process a node instead of node.process. */
func (workflow *Workflow) processNode(node nodes.InterfaceNode, workflowContext *nodes.WorkflowContext) (err error) {
	// Recover from panic during the processing of a node.
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("panic: %s", r))
		}
	}()
	// Return error returned by the processing of a node to the caller function.
	err = node.Process(node, workflowContext)
	return err
}
