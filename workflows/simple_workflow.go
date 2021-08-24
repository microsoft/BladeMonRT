package workflows

import (
	"fmt"
	"github.com/microsoft/BladeMonRT/logging"
	"github.com/microsoft/BladeMonRT/nodes"
	"log"
)

/** Workflow for executing nodes sequentially. */
type SimpleWorkflow struct {
	Workflow
	nodes []nodes.InterfaceNode
}

func NewSimpleWorkflow() *SimpleWorkflow {
	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("SimpleWorkflow")
	return &SimpleWorkflow{Workflow: Workflow{Logger: logger}}
}

func (simpleWorkflow *SimpleWorkflow) AddNode(node nodes.InterfaceNode) {
	simpleWorkflow.nodes = append(simpleWorkflow.nodes, node)
}

func (simpleWorkflow *SimpleWorkflow) GetNodes() []nodes.InterfaceNode {
	return simpleWorkflow.nodes
}

func (simpleWorkflow *SimpleWorkflow) runVirt(workflowContext *nodes.WorkflowContext) error {
	simpleWorkflow.Logger.Println("Running runVirt method.")
	for _, node := range simpleWorkflow.GetNodes() {
		var err error = simpleWorkflow.processNode(node, workflowContext)
		if err != nil {
			simpleWorkflow.Logger.Println(fmt.Sprintf("Workflow aborted."))
			return err
		}
	}
	return nil
}
