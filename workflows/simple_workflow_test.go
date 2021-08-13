package workflows

import (
	"github.com/microsoft/BladeMonRT/nodes"
	"github.com/microsoft/BladeMonRT/nodes/dummy_node_a"
	"testing"
	"gotest.tools/assert"
	"errors"
	gomock "github.com/golang/mock/gomock"
)

func TestWorkflow(t *testing.T) {
	var dummyNodeA nodes.InterfaceNode = dummy_node_a.NewDummyNodeA()
	var dummyNodeB nodes.InterfaceNode = dummy_node_a.NewDummyNodeA()
	var dummyNodeC nodes.InterfaceNode = dummy_node_a.NewDummyNodeA()

	var workflow *SimpleWorkflow = NewSimpleWorkflow()
	workflow.AddNode(dummyNodeA)
	workflow.AddNode(dummyNodeB)
	workflow.AddNode(dummyNodeC)

	var workflowContext *nodes.WorkflowContext = nodes.NewWorkflowContext()
	workflow.Run(workflow, workflowContext)

	// Check that the result at each node includes the predecessor results and the expected hard-coded value.
	resultA := dummyNodeA.GetResult(dummyNodeA, workflowContext)
	assert.Equal(t, resultA, "dummy-node-result");
	resultB := dummyNodeB.GetResult(dummyNodeB, workflowContext)
	assert.Equal(t, resultB, "dummy-node-result|dummy-node-result");
	resultC := dummyNodeC.GetResult(dummyNodeC, workflowContext)
	assert.Equal(t, resultC, "dummy-node-result|dummy-node-result|dummy-node-result|dummy-node-result");
}

func TestAbortWorkflowOnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Assume
	mockNodeA := nodes.NewMockInterfaceNode(ctrl)
	mockNodeB := nodes.NewMockInterfaceNode(ctrl)
	mockNodeC := nodes.NewMockInterfaceNode(ctrl)

	var workflow *SimpleWorkflow = NewSimpleWorkflow()
	workflow.AddNode(mockNodeA)
	workflow.AddNode(mockNodeB)
	workflow.AddNode(mockNodeC)

	// Assert
	// Set up assertion
	mockNodeA.EXPECT().Process(gomock.Any(), gomock.Any()).Return(errors.New("Unable to execute process function."))
	// Expect process is not called on mockNodeB and mockNodeC.

	// Action
	var workflowContext *nodes.WorkflowContext = nodes.NewWorkflowContext()
	workflow.Run(workflow, workflowContext)
}
