package main

import (
	"github.com/microsoft/BladeMonRT/nodes"
	"github.com/microsoft/BladeMonRT/nodes/dummy_node_a"
	"testing"
	"log"
	"io/ioutil"
	"github.com/microsoft/BladeMonRT/test_configs"
	"gotest.tools/assert"
)

func TestSetupWorkflowScheduler(t *testing.T) {
	// Assume
	workflowsJson, err := ioutil.ReadFile(test_configs.TEST_WORKFLOW_FILE)
	if err != nil {
		log.Fatal(err)
	}
	var workflowFactory WorkflowFactory = newWorkflowFactory(workflowsJson, NodeFactory{})

	// Create the schedules JSON from a file with a single schedule.
	schedulesJson, err := ioutil.ReadFile(test_configs.TEST_SINGLE_SCHEDULE_FILE)
	if err != nil {
		log.Fatal(err)
	}

	// Action
	var mainObj *Main = newMain()
	mainObj.setupWorkflows(schedulesJson, workflowFactory)

	// Assert
	// Assert that the GUID to context map has only 1 context.
	var workflowScheduler *WorkflowScheduler = mainObj.WorkflowSchedule
	assert.Equal(t, len(workflowScheduler.guidToContext), 1)
	
	// Check that the first and second nodes in the context's workflow are of type DummyNodeA by checking the value of result field in the objects.
	var firstWorkflowNodes []nodes.InterfaceNode
	for _, element := range workflowScheduler.guidToContext {
		firstWorkflowNodes = element.workflow.GetNodes()
		break
    }
	actualFirstNode :=  *(firstWorkflowNodes[0]).(*dummy_node_a.DummyNodeA)
	actualSecondNode := *(firstWorkflowNodes[1]).(*dummy_node_a.DummyNodeA)
	assert.Equal(t, actualFirstNode.Result, "dummy-node-result")
	assert.Equal(t, actualSecondNode.Result, "dummy-node-result")
}