package root

import (
	"github.com/microsoft/BladeMonRT/nodes"
	"github.com/microsoft/BladeMonRT/workflows"
	"github.com/microsoft/BladeMonRT"
	"testing"
	gomock "github.com/golang/mock/gomock"
	"log"
	"io/ioutil"
	"fmt"
)

func TestWorkflowFactory(t *testing.T) {
	const (
		workflow_file = "test_workflows.json"
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNodeA := nodes.NewMockInterfaceNode(ctrl)
	mockNodeB := nodes.NewMockInterfaceNode(ctrl)
	mockNodeC := nodes.NewMockInterfaceNode(ctrl)

	mockNodeA.EXPECT().Process(gomock.Any(), gomock.Any()).AnyTimes()
	mockNodeB.EXPECT().Process(gomock.Any(), gomock.Any()).AnyTimes()
	mockNodeC.EXPECT().Process(gomock.Any(), gomock.Any()).AnyTimes()

	mockNodeA.EXPECT().GetResult(gomock.Any(), gomock.Any()).Return("mock-node-a-result").AnyTimes()
	mockNodeB.EXPECT().GetResult(gomock.Any(), gomock.Any()).Return("mock-node-b-result").AnyTimes()
	mockNodeC.EXPECT().GetResult(gomock.Any(), gomock.Any()).Return("mock-node-c-result").AnyTimes()

	mockNodeFactory := root.NewMockInterfaceNodeFactory(ctrl)

	workflowsJson, err := ioutil.ReadFile(workflow_file)
	if err != nil {
		log.Fatal(err)
	}
	var workflowFactory root.WorkflowFactory = root.NewWorkflowFactory(workflowsJson, mockNodeFactory)
	mockNodeFactory.EXPECT().ConstructNode("DummyNodeA").Return(mockNodeA)
	mockNodeFactory.EXPECT().ConstructNode("DummyNodeA").Return(mockNodeA)
	mockNodeFactory.EXPECT().ConstructNode("DummyNodeB").Return(mockNodeB)
	mockNodeFactory.EXPECT().ConstructNode("DummyNodeC").Return(mockNodeC)
	

	var workflow workflows.InterfaceWorkflow = workflowFactory.ConstructWorkflow("dummy_workflow")
	workflow.Run(workflow)
	
	var workflowContext *nodes.WorkflowContext = workflow.GetWorkflowContext()
	fmt.Println("hello")
	fmt.Println(workflow.GetNodes())
	for index, node := range workflow.GetNodes() {
		fmt.Println(fmt.Sprintf("Result for node index %d=%s", index, node.GetResult(node, workflowContext).(string)))
	}
}




