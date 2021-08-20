package workflows

import (
	"github.com/microsoft/BladeMonRT/nodes"
	"log"
	"errors"
	"fmt"
)

/** Interface for defining execution sequence of nodes. */
type InterfaceWorkflow interface {
	AddNode(node nodes.InterfaceNode)
	Run(interfaceWorkflow InterfaceWorkflow, workflowContext *nodes.WorkflowContext)
	runVirt(workflowContext *nodes.WorkflowContext) error
	GetNodes() []nodes.InterfaceNode
}

/** Concrete type for defining execution sequence of nodes. */
type Workflow struct {
	Logger *log.Logger
}

func (workflow *Workflow) Run(interfaceWorkflow InterfaceWorkflow, workflowContext *nodes.WorkflowContext) {
	workflow.Logger.Println("Running run method.")

	// Set the nodes in the workflow context to the nodes in this workflow.
	workflowContext.SetNodes(interfaceWorkflow.GetNodes())

	var err error = interfaceWorkflow.runVirt(workflowContext)
	// Handle errors thrown when running the workflow.
	if (err != nil) {
		workflow.Logger.Println(fmt.Sprintf("Workflow error: %s", err))
	}
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