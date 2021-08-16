package workflows

import (
	"log"
	"github.com/microsoft/BladeMonRT/nodes"
	"github.com/microsoft/BladeMonRT/logging"
	"fmt"
	"errors"
)

/** Workflow for executing nodes sequentially. */
type SimpleWorkflow struct {
	Workflow
	nodes []nodes.InterfaceNode
}

func NewSimpleWorkflow() *SimpleWorkflow{
	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("SimpleWorkflow")
	return &SimpleWorkflow{Workflow: Workflow{Logger : logger}}
}

func (simpleWorkflow *SimpleWorkflow) AddNode(node nodes.InterfaceNode) {
	simpleWorkflow.nodes = append(simpleWorkflow.nodes, node)
}

func (simpleWorkflow *SimpleWorkflow) GetNodes() []nodes.InterfaceNode {
	return simpleWorkflow.nodes
}

func (simpleWorkflow *SimpleWorkflow) runVirt(workflowContext *nodes.WorkflowContext) {
	simpleWorkflow.Logger.Println("Running runVirt method.")
	for _, node := range simpleWorkflow.GetNodes() {
		var err error = simpleWorkflow.processNode(node, workflowContext)
		if (err != nil) {
			simpleWorkflow.Logger.Println(fmt.Sprintf("Aborting the workflow due to error: %s", err))
			break
		}
	}
}

func (simpleWorkflow *SimpleWorkflow) processNode(node nodes.InterfaceNode, workflowContext *nodes.WorkflowContext) (processNodeError error) {
	// Recover from panic during the processing of a node.
	defer func() {
		if r := recover(); r != nil {
			processNodeError = errors.New("Panic during execution of processNode function.") 
		}
	}()
	// Return error returned by the processing of a node to the caller function.
	processNodeError = node.Process(node, workflowContext)
	return processNodeError
}