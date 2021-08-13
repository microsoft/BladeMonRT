package workflows

import (
	"log"
	"github.com/microsoft/BladeMonRT/nodes"
	"github.com/microsoft/BladeMonRT/logging"
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
		var err error = node.Process(node, workflowContext)
		if (err != nil) {
			break
		}
	}
}