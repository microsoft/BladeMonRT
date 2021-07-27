package workflows

import (
	"github.com/microsoft/BladeMonRT/nodes"
)

/** Workflow for executing nodes sequentially. */
type SimpleWorkflow struct {
	Workflow
	nodes []nodes.InterfaceNode
}

func NewSimpleWorkflow() *SimpleWorkflow{
	return &SimpleWorkflow{Workflow: Workflow{}}
}

func (simpleWorkflow *SimpleWorkflow) AddNode(node nodes.InterfaceNode) {
	simpleWorkflow.nodes = append(simpleWorkflow.nodes, node)
}

func (simpleWorkflow *SimpleWorkflow) GetNodes() []nodes.InterfaceNode {
	return simpleWorkflow.nodes
}

func (simpleWorkflow *SimpleWorkflow) runVirt(workflowContext *nodes.WorkflowContext) {
	for _, node := range simpleWorkflow.GetNodes() {
		node.Process(node, workflowContext)
	}
}