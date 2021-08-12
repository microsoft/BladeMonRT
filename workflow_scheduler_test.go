package main

import (
	workflows "github.com/microsoft/BladeMonRT/workflows"
	gomock "github.com/golang/mock/gomock"
	"testing"
	"github.com/microsoft/BladeMonRT/utils"
	"github.com/microsoft/BladeMonRT/test_utils"
	wevtapi "github.com/0xrawsec/golang-win32/win32/wevtapi"
	"unsafe"
	win32 "github.com/0xrawsec/golang-win32/win32"
	"log"
	"github.com/microsoft/BladeMonRT/logging"
	"time"
)

type UtilsForTest struct {
}

func (utilsForTest UtilsForTest) ParseEventXML(eventXML string) utils.EtwEvent {
	return utils.EtwEvent{}
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