package main

import (
	"encoding/json"
	"github.com/microsoft/BladeMonRT/nodes"
	"github.com/microsoft/BladeMonRT/workflows"
)

/** Class for parsing workflow definitions. */
type WorkflowFactory struct {
	nameToWorkflow map[string]WorkflowDescription
}

/** Class for the workflow description in the JSON. */
type WorkflowDescription struct {
	Type string `json:"type"`
	Nodes []NodeDescription `json:"nodes"`
}

/** Class for the node description in the JSON. */
type NodeDescription struct {
	Name string `json:"name"`
	Metadata interface{} `json:"metadata"`
	Args interface{} `json:"args"`
}

func newWorkflowFactory(workflowsJson []byte) WorkflowFactory {
	var workflows map[string]map[string]WorkflowDescription
	json.Unmarshal([]byte(workflowsJson), &workflows)

	return WorkflowFactory{nameToWorkflow : workflows["workflows"]}
}

func (WorkflowFactory *WorkflowFactory) constructWorkflow(workflowName string) workflows.InterfaceWorkflow {
	var currWorkflowDescription WorkflowDescription
	currWorkflowDescription = WorkflowFactory.nameToWorkflow[workflowName]

	// Create a collection of nodes from their names.
	var workflowNodes []nodes.InterfaceNode
	var nodeFactory NodeFactory = NodeFactory{}
	for _, nodeDescription := range currWorkflowDescription.Nodes {
		var node nodes.InterfaceNode = nodeFactory.constructNode(nodeDescription.Name)
		workflowNodes = append(workflowNodes, node)
	}

	// Add nodes to the workflow using the worfklow type
	switch currWorkflowDescription.Type {	
		case "simple":
			var workflow workflows.InterfaceWorkflow
			workflow = &workflows.SimpleWorkflow{}
			for _, node := range workflowNodes {
				workflow.AddNode(node)
			}
			return workflow
		default:
			panic("Workflow types other than simple are not implemented.")
	}

	return nil
}