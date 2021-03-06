package main

import (
	win32 "github.com/0xrawsec/golang-win32/win32"
	wevtapi "github.com/0xrawsec/golang-win32/win32/wevtapi"
	gomock "github.com/golang/mock/gomock"
	"github.com/microsoft/BladeMonRT/logging"
	"github.com/microsoft/BladeMonRT/nodes"
	"github.com/microsoft/BladeMonRT/nodes/dummy_node_a"
	"github.com/microsoft/BladeMonRT/store"
	"github.com/microsoft/BladeMonRT/test_configs"
	"github.com/microsoft/BladeMonRT/test_utils"
	"github.com/microsoft/BladeMonRT/utils"
	workflows "github.com/microsoft/BladeMonRT/workflows"
	"gotest.tools/assert"
	"io/ioutil"
	"log"
	"testing"
	"time"
	"unsafe"
)

type UtilsForTest struct {
}

func (utilsForTest UtilsForTest) ParseEventXML(eventXML string) utils.EtwEvent {
	return utils.EtwEvent{Provider: "disk", EventID: 7, TimeCreated: time.Now(), EventRecordID: 6}
}

func (utilsForTest UtilsForTest) GetEventRecordIdBookmark(bookmarkStore store.PersistentKeyValueStoreInterface, query string) int {
	return 0
}

type UtilsForTestWithOldEvent struct {
}

func (utilsForTest UtilsForTestWithOldEvent) ParseEventXML(eventXML string) utils.EtwEvent {
	twentyFiveHoursAgo := time.Now().Add(time.Hour * -25)
	return utils.EtwEvent{Provider: "disk", EventID: 7, TimeCreated: twentyFiveHoursAgo, EventRecordID: 6}
}

func (utilsForTest UtilsForTestWithOldEvent) GetEventRecordIdBookmark(bookmarkStore store.PersistentKeyValueStoreInterface, query string) int {
	return 0
}

func TestSetupWorkflowsBasic(t *testing.T) {
	// Assume
	config := test_configs.NewTestConfig()
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
	var mainObj *Main = newMain(config)
	mainObj.setupWorkflows(schedulesJson, workflowFactory)

	// Assert
	// Assert that the GUID to context map has only 1 context.
	var workflowScheduler *WorkflowScheduler = mainObj.WorkflowScheduler.(*WorkflowScheduler)
	assert.Equal(t, len(workflowScheduler.guidToContext), 1)

	// Check that the first and second nodes in the context's workflow are of type DummyNodeA by checking the value of result field in the objects.
	var firstWorkflowNodes []nodes.InterfaceNode
	var context *CallbackContext
	for _, context = range workflowScheduler.guidToContext {
		firstWorkflowNodes = context.workflow.GetNodes()
		break
	}
	actualFirstNode := *(firstWorkflowNodes[0]).(*dummy_node_a.DummyNodeA)
	actualSecondNode := *(firstWorkflowNodes[1]).(*dummy_node_a.DummyNodeA)
	assert.Equal(t, actualFirstNode.Result, "dummy-node-result")
	assert.Equal(t, actualSecondNode.Result, "dummy-node-result")
}

func TestAddWinEventBasedSchedule_Basic(t *testing.T) {
	// Case 1: Call the addWinEventBasedSchedule method with a schedule containing a query that does not contain a condition.

	// Assume
	config := test_configs.NewTestConfig()
	workflowScheduler := newWorkflowScheduler(config)
	workflow := workflows.NewSimpleWorkflow()
	cpuSpeedMonitoringQuery := WinEventSubscribeQuery{channel: "System", query: `*[System[Provider[@Name='CpuSpeedMonitoring']]]`}
	var eventQueries []WinEventSubscribeQuery = []WinEventSubscribeQuery{cpuSpeedMonitoringQuery}

	// Action
	workflowScheduler.addWinEventBasedSchedule(workflow, eventQueries)

	// Assert
	// Assert that the GUID to context map has only 1 context.
	assert.Equal(t, len(workflowScheduler.guidToContext), 1)
	var context *CallbackContext
	for _, context = range workflowScheduler.guidToContext {
		break
	}

	// Check that context's contents related to the bookmark feature.
	assert.Equal(t, context.queryIncludesCondition, false)
	assert.Equal(t, context.bookmarkStore, nil)
	assert.Equal(t, context.query, "*[System[Provider[@Name='CpuSpeedMonitoring']]]")
}

func TestAddWinEventBasedSchedule_QueryWithCondition(t *testing.T) {
	// Case 2: Call the addWinEventBasedSchedule method with a schedule containing a query that contains a condition.

	// Assume
	config := test_configs.NewTestConfig()
	workflowScheduler := newWorkflowScheduler(config)
	workflow := workflows.NewSimpleWorkflow()
	diskQuery := WinEventSubscribeQuery{channel: "System", query: `["System", "*[System[Provider[@Name='disk'] and EventID=7 and EventRecordID > {condition}]]"]`}
	var eventQueries []WinEventSubscribeQuery = []WinEventSubscribeQuery{diskQuery}

	// Action
	workflowScheduler.addWinEventBasedSchedule(workflow, eventQueries)

	// Assert
	// Assert that the GUID to context map has only 1 context.
	assert.Equal(t, len(workflowScheduler.guidToContext), 1)
	var context *CallbackContext
	for _, context = range workflowScheduler.guidToContext {
		break
	}

	// Check that context's contents related to the bookmark feature.
	assert.Equal(t, context.queryIncludesCondition, true)
	assert.Assert(t, context.bookmarkStore != nil)
	assert.Equal(t, context.query, `["System", "*[System[Provider[@Name='disk'] and EventID=7 and EventRecordID > {condition}]]"]`)
}

func TestSubscriptionCallback_Basic(t *testing.T) {
	// Case 1: Call the SubscriptionCallback method with a query that does not contain a condition.

	// Assume
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("WorkflowScheduler")
	var guidToContext map[string]*CallbackContext = make(map[string]*CallbackContext)
	config := test_configs.NewTestConfig()
	mockBookmarkStore := store.NewMockPersistentKeyValueStoreInterface(ctrl)
	var workflowScheduler *WorkflowScheduler = &WorkflowScheduler{config: config, logger: logger, guidToContext: guidToContext, bookmarkStore: mockBookmarkStore, utils: UtilsForTest{}}
	mockWorkflow := workflows.NewMockInterfaceWorkflow(ctrl)
	var callbackContext *CallbackContext = &CallbackContext{workflow: mockWorkflow, bookmarkStore: mockBookmarkStore, queryIncludesCondition: false}
	workflowScheduler.guidToContext["50bd065e-f3e9-4887-8093-b171f1b01372"] = callbackContext

	// Set up assertion
	mockWorkflow.EXPECT().Run(gomock.Any(), gomock.Any())
	// Expect no calls made on mockBookmarkStore.

	// Action
	workflowScheduler.SubscriptionCallback(wevtapi.EvtSubscribeActionDeliver, win32.PVOID(unsafe.Pointer(test_utils.ToCString("50bd065e-f3e9-4887-8093-b171f1b01372"))), wevtapi.EVT_HANDLE(uintptr(0)))

	// Wait for 5 seconds since the main thread has to switch to the goroutine to run the workflow before Run() is called on mockWorkflow.
	// If we do not wait, the assertion that Run() was called on mockWorkflow will fail.
	time.Sleep(5 * time.Second)
}

func TestSubscriptionCallback_QueryWithCondition_1(t *testing.T) {
	// Case 2: Call the SubscriptionCallback method with a query that contains a condition. The bookmark store contains the eventRecordID for a query less than the eventRecordID in the current event.

	// Assume
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("WorkflowScheduler")
	var guidToContext map[string]*CallbackContext = make(map[string]*CallbackContext)
	config := test_configs.NewTestConfig()
	mockBookmarkStore := store.NewMockPersistentKeyValueStoreInterface(ctrl)
	var workflowScheduler *WorkflowScheduler = &WorkflowScheduler{config: config, logger: logger, guidToContext: guidToContext, bookmarkStore: mockBookmarkStore, utils: UtilsForTest{}}
	var queryWithCondition string = "*[System[Provider[@Name='disk'] and EventID=7 and EventRecordID > {condition}]]"
	workflow := workflows.NewSimpleWorkflow()

	var callbackContext *CallbackContext = &CallbackContext{workflow: workflow, query: queryWithCondition, bookmarkStore: mockBookmarkStore, queryIncludesCondition: true}
	workflowScheduler.guidToContext["50bd065e-f3e9-4887-8093-b171f1b01372"] = callbackContext

	// Set up assertions
	// Event has EventRecordID=6 as specified in the return value of 'ParseEventXML'.
	mockBookmarkStore.EXPECT().GetValue(gomock.Any()).Return("5", nil)
	mockBookmarkStore.EXPECT().SetValue(queryWithCondition, "6")

	// Action
	workflowScheduler.SubscriptionCallback(wevtapi.EvtSubscribeActionDeliver, win32.PVOID(unsafe.Pointer(test_utils.ToCString("50bd065e-f3e9-4887-8093-b171f1b01372"))), wevtapi.EVT_HANDLE(uintptr(0)))

	// Wait for 5 seconds since the main thread has to switch to the goroutine to run the workflow before Run() is called on workflow.
	// If we do not wait, the assertion that SetValue was called may fail. (Run() calls SetValue on the bookmark store.)
	time.Sleep(5 * time.Second)
}

func TestSubscriptionCallback_QueryWithCondition_2(t *testing.T) {
	// Case 3: Call the SubscriptionCallback method with a query that contains a condition. The bookmark store contains a eventRecordID for the query greater than the eventRecordID in the current event.

	// Assume
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("WorkflowScheduler")
	var guidToContext map[string]*CallbackContext = make(map[string]*CallbackContext)
	mockBookmarkStore := store.NewMockPersistentKeyValueStoreInterface(ctrl)
	mockUtils := utils.NewMockUtilsInterface(ctrl)
	config := test_configs.NewTestConfig()
	var workflowScheduler *WorkflowScheduler = &WorkflowScheduler{config: config, logger: logger, guidToContext: guidToContext, bookmarkStore: mockBookmarkStore, utils: mockUtils}
	var queryWithCondition string = "*[System[Provider[@Name='disk'] and EventID=7 and EventRecordID > {condition}]]"
	workflow := workflows.NewSimpleWorkflow()

	var callbackContext *CallbackContext = &CallbackContext{workflow: workflow, query: queryWithCondition, bookmarkStore: mockBookmarkStore, queryIncludesCondition: true}
	workflowScheduler.guidToContext["50bd065e-f3e9-4887-8093-b171f1b01372"] = callbackContext

	// Set up assertions
	mockUtils.EXPECT().ParseEventXML(gomock.Any()).Return(utils.EtwEvent{Provider: "disk", EventID: 7, TimeCreated: time.Now(), EventRecordID: 6})
	mockBookmarkStore.EXPECT().GetValue(gomock.Any()).Return("7", nil)
	// Expect SetValue is not called because EventRecordID is 6 which is less than the EventRecordID in the bookmark store which is 7.

	// Action
	workflowScheduler.SubscriptionCallback(wevtapi.EvtSubscribeActionDeliver, win32.PVOID(unsafe.Pointer(test_utils.ToCString("50bd065e-f3e9-4887-8093-b171f1b01372"))), wevtapi.EVT_HANDLE(uintptr(0)))

	// Wait for 5 seconds since the main thread has to switch to the goroutine to run the workflow before Run() is called on workflow.
	time.Sleep(5 * time.Second)
}

func TestSubscriptionCallback_OldEvent(t *testing.T) {
	// Case 4: Call the SubscriptionCallback method with an event older than MaxAgeToProcessWinEvtsInDays in the config.
	// Requires MaxAgeToProcessWinEvtsInDays in the test config is 1.

	// Assume
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("WorkflowScheduler")
	var guidToContext map[string]*CallbackContext = make(map[string]*CallbackContext)

	// 'ParseEventXML' from UtilsForTestWithOldEvent returns an event from 25 hours ago.
	config := test_configs.NewTestConfig()
	var workflowScheduler *WorkflowScheduler = &WorkflowScheduler{config: config, logger: logger, guidToContext: guidToContext, utils: UtilsForTestWithOldEvent{}}
	mockWorkflow := workflows.NewMockInterfaceWorkflow(ctrl)
	var callbackContext *CallbackContext = &CallbackContext{workflow: mockWorkflow}
	workflowScheduler.guidToContext["50bd065e-f3e9-4887-8093-b171f1b01372"] = callbackContext

	// Set up assertions
	// Expect that mockWorkflow is not run because the event is 25 hours old. If any function is called on the mock workflow, the test will throw an error.

	// Action
	workflowScheduler.SubscriptionCallback(wevtapi.EvtSubscribeActionDeliver, win32.PVOID(unsafe.Pointer(test_utils.ToCString("50bd065e-f3e9-4887-8093-b171f1b01372"))), wevtapi.EVT_HANDLE(uintptr(0)))
	// Wait for 5 seconds since the main thread has to switch to the goroutine to run the workflow before Run() may be called on mockWorkflow.
	time.Sleep(5 * time.Second)
}

func TestDecideSubscriptionType_QueryExistsInBookmarkStore(t *testing.T) {
	// Assume
	context := &CallbackContext{}
	context.query = "*[System[Provider[@Name='disk'] and EventID=7 and EventRecordID > {condition}]]"
	context.queryIncludesCondition = true
	config := test_configs.NewTestConfig()
	workflowScheduler := newWorkflowScheduler(config)
	workflowScheduler.bookmarkStore.Clear()
	workflowScheduler.bookmarkStore.SetValue(context.query, "7")

	// Action
	queryText, subscribeMethod := workflowScheduler.decideSubscriptionType(context)

	// Assert
	assert.Equal(t, queryText, "*[System[Provider[@Name='disk'] and EventID=7 and EventRecordID > 7]]")
	assert.Equal(t, subscribeMethod, win32.DWORD(wevtapi.EvtSubscribeStartAtOldestRecord))
}

func TestDecideSubscriptionType_QueryDoesNotExistInBookmarkStore(t *testing.T) {
	// Assume
	context := &CallbackContext{}
	context.query = "*[System[Provider[@Name='disk'] and EventID=7 and EventRecordID > {condition}]]"
	context.queryIncludesCondition = true
	config := test_configs.NewTestConfig()
	workflowScheduler := newWorkflowScheduler(config)
	workflowScheduler.bookmarkStore.Clear()

	// Action
	queryText, subscribeMethod := workflowScheduler.decideSubscriptionType(context)

	// Assert
	assert.Equal(t, queryText, "*[System[Provider[@Name='disk'] and EventID=7 and EventRecordID > 0]]")
	assert.Equal(t, subscribeMethod, win32.DWORD(wevtapi.EvtSubscribeToFutureEvents))
}
