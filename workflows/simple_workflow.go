package workflows

import (
	"example.com/nodes"
)

// Workflow for executing nodes sequentially
type SimpleWorkflow struct {
	Workflow
}

func (simpleWorkflow *SimpleWorkflow) AddNode(node nodes.InterfaceNode) {
	simpleWorkflow.Nodes = append(simpleWorkflow.Nodes, node)
}

func (simpleWorkflow *SimpleWorkflow) RunVirt(workflowContextResult map[string]string) {
	for _, node := range simpleWorkflow.Nodes {
		node.ProcessVirt(workflowContextResult)
	}
}