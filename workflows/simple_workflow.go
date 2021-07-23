package workflows

import (
	"github.com/microsoft/BladeMonRT/nodes"
)

/** Workflow for executing nodes sequentially. */
type SimpleWorkflow struct {
	Workflow
	Nodes []nodes.InterfaceNode
}

func (simpleWorkflow *SimpleWorkflow) AddNode(node nodes.InterfaceNode) {
	simpleWorkflow.Nodes = append(simpleWorkflow.Nodes, node)
}


func (simpleWorkflow *SimpleWorkflow) GetNodes() []nodes.InterfaceNode {
	return simpleWorkflow.Nodes
}

func (simpleWorkflow *SimpleWorkflow) runVirt() {
	var predecessorNodesResults []interface{}
	for _, node := range simpleWorkflow.Nodes {
		node.Process(node, predecessorNodesResults)
		predecessorNodesResults = append(predecessorNodesResults, node.GetResult()) 
	}
}


