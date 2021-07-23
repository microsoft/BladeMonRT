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

func (simpleWorkflow *SimpleWorkflow) RunVirt() {
	var predecessorNodeResults []interface{}
	for _, node := range simpleWorkflow.Nodes {
		node.ProcessVirt(predecessorNodeResults)
		predecessorNodeResults = append(predecessorNodeResults, node.GetResult()) 
	}
}

func (simpleWorkflow *SimpleWorkflow) GetResult() map[string]interface{} {
	var nodeToResult map[string]interface{} = make(map[string]interface{})
	for _, node := range simpleWorkflow.Nodes {
		nodeToResult[node.GetName()] = node.GetResult()
	}
	return nodeToResult
}

