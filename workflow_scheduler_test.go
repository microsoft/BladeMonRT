package main

import (
	"github.com/microsoft/BladeMonRT/nodes"
	"github.com/microsoft/BladeMonRT/nodes/dummy_node_a"
	"testing"
	"log"
	"io/ioutil"
	"github.com/microsoft/BladeMonRT/test_configs"
	"gotest.tools/assert"
	workflows "github.com/microsoft/BladeMonRT/workflows"
	gomock "github.com/golang/mock/gomock"
	"github.com/microsoft/BladeMonRT/utils"
	"github.com/microsoft/BladeMonRT/test_utils"
	wevtapi "github.com/0xrawsec/golang-win32/win32/wevtapi"
	"unsafe"
	win32 "github.com/0xrawsec/golang-win32/win32"
	"github.com/microsoft/BladeMonRT/logging"
	"time"
)

type UtilsForTest struct {
}

func (utilsForTest UtilsForTest) ParseEventXML(eventXML string) utils.EtwEvent {
	return utils.EtwEvent{}
}

func TestSetupWorkflows(t *testing.T) {
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
	var workflowScheduler *WorkflowScheduler = mainObj.WorkflowScheduler.(*WorkflowScheduler)
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

func TestSubscriptionCallback(t *testing.T) {
	// Assume
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("WorkflowScheduler")
	var guidToContext map[string]*CallbackContext = make(map[string]*CallbackContext)
	var workflowScheduler *WorkflowScheduler = &WorkflowScheduler{logger: logger, guidToContext : guidToContext, utils : UtilsForTest{}}

	mockWorkflow := workflows.NewMockInterfaceWorkflow(ctrl)
	// Set up assertion
	mockWorkflow.EXPECT().Run(gomock.Any(), gomock.Any())

	// Assume
	var callbackContext *CallbackContext = &CallbackContext{workflow: mockWorkflow}
	workflowScheduler.guidToContext["50bd065e-f3e9-4887-8093-b171f1b01372"] = callbackContext

	// Action
	workflowScheduler.SubscriptionCallback(wevtapi.EvtSubscribeActionDeliver, win32.PVOID(unsafe.Pointer(test_utils.ToCString("50bd065e-f3e9-4887-8093-b171f1b01372"))), wevtapi.EVT_HANDLE(uintptr(0)))

	// Wait for 5 seconds since the main thread has to switch to the goroutine to run the workflow before Run() is called on mockWorkflow. 
	// If we do not wait, the assertion that Run() was called on mockWorkflow will fail.
	time.Sleep(5 * time.Second)
}