package main

import (
	"encoding/json"
	"github.com/microsoft/BladeMonRT/nodes"
	"github.com/microsoft/BladeMonRT/workflows"
)

/** Class for parsing workflow definitions. */
type WorkflowManager struct {
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

func newWorkflowManager(workflowsJson []byte) WorkflowManager {
	var workflows map[string]map[string]WorkflowDescription
	json.Unmarshal([]byte(workflowsJson), &workflows)

	// var nodeRegistry NodeRegistry = newNodeRegistrySingleton()

	return WorkflowManager{nameToWorkflow : workflows["workflows"]} //, nodeRegistry : nodeRegistry}
}

func (workflowManager *WorkflowManager) constructWorkflow(workflowName string) workflows.InterfaceWorkflow {
	var currWorkflowDescription WorkflowDescription
	currWorkflowDescription = workflowManager.nameToWorkflow[workflowName]
	if currWorkflowDescription.Type != "simple" {
		panic("Workflow types other than simple are not implemented.") 
	}

	// Create a collection of nodes from their names.
	var workflowNodes []nodes.InterfaceNode
	for _, nodeDescription := range currWorkflowDescription.Nodes {
		var node nodes.InterfaceNode = makeInstance(nodeDescription.Name)
		node.InitializeFields();
		workflowNodes = append(workflowNodes, node)
	}

	// Create a simple workflow and add nodes to the workflow.
	var workflow workflows.InterfaceWorkflow
	workflow = &workflows.SimpleWorkflow{}
	for _, node := range workflowNodes {
		workflow.AddNode(node)
	}

	return workflow
}