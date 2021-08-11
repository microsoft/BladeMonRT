package main

import (
	"github.com/microsoft/BladeMonRT/nodes"
	"github.com/microsoft/BladeMonRT/nodes/dummy_node_a"
	//"github.com/microsoft/BladeMonRT/nodes/dummy_node_b"
	workflows "github.com/microsoft/BladeMonRT/workflows"
	gomock "github.com/golang/mock/gomock"
	"testing"
	"github.com/microsoft/BladeMonRT/utils"
	wevtapi "github.com/0xrawsec/golang-win32/win32/wevtapi"
	"unsafe"
	win32 "github.com/0xrawsec/golang-win32/win32"
	"log"
	"github.com/microsoft/BladeMonRT/logging"
	"time"
	"io/ioutil"
	"github.com/microsoft/BladeMonRT/test_configs"
	 "github.com/stretchr/testify/assert"
	 "fmt"
)

type UtilsForTest struct {
}

func (utilsForTest UtilsForTest) ParseEventXML(eventXML string) utils.EventFromXML {
	return utils.EventFromXML{}
}

func TestNewWorkflowScheduler(t *testing.T) {
	workflowsJson, err := ioutil.ReadFile(test_configs.TEST_WORKFLOW_FILE)
	if err != nil {
		log.Fatal(err)
	}
	var workflowFactory WorkflowFactory = newWorkflowFactory(workflowsJson, NodeFactory{})

	schedulesJson, err := ioutil.ReadFile(test_configs.TEST_SINGLE_SCHEDULE_FILE)
	if err != nil {
		log.Fatal(err)
	}
	var workflowScheduler *WorkflowScheduler = newWorkflowScheduler(schedulesJson, workflowFactory)
	var firstWorkflowNodes []nodes.InterfaceNode
	fmt.Println(workflowScheduler.guidToContext)
	//var secondWorkflowNodes []nodes.InterfaceNode
	var index int = 0
	for _, element := range workflowScheduler.guidToContext {
		if (index == 0) {
			firstWorkflowNodes = element.workflow.GetNodes()
		}
		index++
    }
	
	assert.EqualValues(t, *(firstWorkflowNodes[0]).(*dummy_node_a.DummyNodeA), *dummy_node_a.NewDummyNodeA())
	assert.EqualValues(t, *(firstWorkflowNodes[1]).(*dummy_node_a.DummyNodeA), *dummy_node_a.NewDummyNodeA())
/*
assert.EqualValues(t, *(firstWorkflowNodes[0]).(*dummy_node_b.DummyNodeB), *dummy_node_b.NewDummyNodeB())
	assert.EqualValues(t, *(firstWorkflowNodes[1]).(*dummy_node_c.DummyNodeC), *dummy_node_c.NewDummyNodeC())
	*/
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
	workflowScheduler.SubscriptionCallback(wevtapi.EvtSubscribeActionDeliver, win32.PVOID(unsafe.Pointer(GUIDForTest())), wevtapi.EVT_HANDLE(uintptr(0)))

	// Wait for 5 seconds since the current thread has to switch to the goroutine to run the workflow before Run() is called on mockWorkflow. 
	// If we do not wait, the assertion that Run() was called on mockWorkflow will fail.
	time.Sleep(5 * time.Second)
}