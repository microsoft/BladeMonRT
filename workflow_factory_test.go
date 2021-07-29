package main

import (
	nodes "github.com/microsoft/BladeMonRT/nodes"
	workflows "github.com/microsoft/BladeMonRT/workflows"
	gomock "github.com/golang/mock/gomock"
	"testing"
	"log"
	"io/ioutil"
	"gotest.tools/assert"
)

func TestWorkflowFactory(t *testing.T) {
	const (
		workflow_file = "test_configs/test_workflows.json"
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNodeA := nodes.NewMockInterfaceNode(ctrl)
	mockNodeB := nodes.NewMockInterfaceNode(ctrl)
	mockNodeC := nodes.NewMockInterfaceNode(ctrl)

	mockNodeA.EXPECT().Process(gomock.Any(), gomock.Any()).AnyTimes()
	mockNodeB.EXPECT().Process(gomock.Any(), gomock.Any()).AnyTimes()
	mockNodeC.EXPECT().Process(gomock.Any(), gomock.Any()).AnyTimes()

	mockNodeA.EXPECT().GetResult(gomock.Any(), gomock.Any()).Return("node-a-result").AnyTimes()
	mockNodeB.EXPECT().GetResult(gomock.Any(), gomock.Any()).Return("node-b-result").AnyTimes()
	mockNodeC.EXPECT().GetResult(gomock.Any(), gomock.Any()).Return("node-c-result").AnyTimes()

	mockNodeFactory := NewMockInterfaceNodeFactory(ctrl)
	workflowsJson, err := ioutil.ReadFile(workflow_file)
	if err != nil {
		log.Fatal(err)
	}
	var workflowFactory WorkflowFactory = newWorkflowFactory(workflowsJson, mockNodeFactory)
	mockNodeFactory.EXPECT().constructNode("TestNodeA").Return(mockNodeA)
	mockNodeFactory.EXPECT().constructNode("TestNodeA").Return(mockNodeA)
	mockNodeFactory.EXPECT().constructNode("TestNodeB").Return(mockNodeB)
	mockNodeFactory.EXPECT().constructNode("TestNodeC").Return(mockNodeC)

	var workflow workflows.InterfaceWorkflow = workflowFactory.constructWorkflow("dummy_workflow")
	var workflowContext *nodes.WorkflowContext = nodes.NewWorkflowContext()
	workflow.Run(workflow, workflowContext)
	
	// Check the results of nodes in the workflow.
	var workflowNodes []nodes.InterfaceNode = workflow.GetNodes()
	assert.Equal(t, workflowNodes[0].GetResult(workflowNodes[0], workflowContext), "node-a-result");
	assert.Equal(t, workflowNodes[1].GetResult(workflowNodes[1], workflowContext), "node-a-result");
	assert.Equal(t, workflowNodes[2].GetResult(workflowNodes[2], workflowContext), "node-b-result");
	assert.Equal(t, workflowNodes[3].GetResult(workflowNodes[3], workflowContext), "node-c-result");
}