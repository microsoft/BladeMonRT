package main

import (
	"encoding/json"
	"github.com/microsoft/BladeMonRT/nodes"
	"github.com/microsoft/BladeMonRT/workflows"
	"github.com/microsoft/BladeMonRT/logging"
	"log"
)

/** Class for parsing workflow definitions. */
type WorkflowFactory struct {
	nameToWorkflow map[string]WorkflowDescription
	nodeFactory InterfaceNodeFactory
	logger *log.Logger
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

func newWorkflowFactory(workflowsJson []byte, nodeFactory InterfaceNodeFactory) WorkflowFactory {
	var workflows map[string]map[string]WorkflowDescription
	json.Unmarshal([]byte(workflowsJson), &workflows)

	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("WorkflowFactory")
	return WorkflowFactory{nameToWorkflow : workflows["workflows"], nodeFactory : nodeFactory, logger : logger}
}

func (workflowFactory *WorkflowFactory) constructWorkflow(workflowName string) workflows.InterfaceWorkflow {
	workflowFactory.logger.Println("Constructing workflow.")

	var currWorkflowDescription WorkflowDescription
	currWorkflowDescription = workflowFactory.nameToWorkflow[workflowName]

	// Create a collection of nodes from their names.
	var workflowNodes []nodes.InterfaceNode
	var nodeFactory InterfaceNodeFactory = workflowFactory.nodeFactory
	for _, nodeDescription := range currWorkflowDescription.Nodes {
		var node nodes.InterfaceNode = nodeFactory.constructNode(nodeDescription.Name)
		workflowNodes = append(workflowNodes, node)
	}

	// Add nodes to the workflow using the workflow type.
	switch currWorkflowDescription.Type {	
		case "simple":
			var workflow workflows.InterfaceWorkflow
			workflow = workflows.NewSimpleWorkflow()
			for _, node := range workflowNodes {
				workflow.AddNode(node)
			}
			return workflow
		default:
			panic("Workflow types other than simple are not implemented.")
	}

	return nil
}