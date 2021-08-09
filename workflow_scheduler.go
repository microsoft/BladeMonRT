package main

import (
	"C"
	"encoding/json"
	"github.com/microsoft/BladeMonRT/workflows"
	"github.com/microsoft/BladeMonRT/nodes"
	"github.com/microsoft/BladeMonRT/logging"
	"github.com/microsoft/BladeMonRT/configs"
	"log"
	"strings"
	"fmt"
	win32 "github.com/0xrawsec/golang-win32/win32"
	wevtapi "github.com/0xrawsec/golang-win32/win32/wevtapi"
	"unsafe"
	"github.com/microsoft/BladeMonRT/utils"
	"time"
	"github.com/google/uuid"
)

/** Class for scheduling workflows. */
type WorkflowScheduler struct {
	schedules []interface{}
	logger *log.Logger
	eventSubscriptionHandles []wevtapi.EVT_HANDLE
	queryToEventRecordIdBookmark map[string]int
	guidToContext map[string]*CallbackContext
}

/** Class that represents a query for subscribing to a windows event. */
type WinEventSubscribeQuery struct {
	channel string
	query string
}

/** Class for the schedule description in the JSON. */
type ScheduleDescription struct {
    Name string `json:"name"`
    ScheduleType string `json:"schedule_type"`
    Workflow string `json:"workflow"`
    Enable bool `json:"enable"`

	// Attributes specific to events of type 'on_win_event'.
	WinEventSubscribeQueries [][]string `json:"win_event_subscribe_queries"`
}

/** Class that holds information used in the subscription callback function. */
type CallbackContext struct {
	workflow workflows.InterfaceWorkflow
	provider string
	eventID int
	timeCreated time.Time
	eventRecordID int
}

func runWorkflow(context *CallbackContext) {
	fmt.Println("run workflow")
	var workflowContext *nodes.WorkflowContext = nodes.NewWorkflowContext()
	fmt.Println(context.workflow)
	var workflow workflows.InterfaceWorkflow = context.workflow
	workflow.Run(workflow, workflowContext)
}

func (workflowScheduler *WorkflowScheduler) SubscriptionCallback(Action wevtapi.EVT_SUBSCRIBE_NOTIFY_ACTION, UserContext win32.PVOID, Event wevtapi.EVT_HANDLE) uintptr {
	var CStringGuid *C.char = (*C.char)(unsafe.Pointer(UserContext))
	var guid string = C.GoStringN(CStringGuid, 32)
	var callbackContext *CallbackContext = workflowScheduler.guidToContext[guid]

	switch Action {
		case wevtapi.EvtSubscribeActionError:
			workflowScheduler.logger.Println("Win event subscription failed.")
			return 0
		case wevtapi.EvtSubscribeActionDeliver:
			UTF16EventXML, err := wevtapi.EvtRenderXML(Event)
			if err != nil {
				workflowScheduler.logger.Println("Error converting event to XML: %s", err)
			}
			eventXML := win32.UTF16BytesToString(UTF16EventXML)

			provider, eventID, timeCreated, eventRecordID := utils.NewUtils().ParseEventXML(eventXML)
			var nowTime Time.time = time.Now()

			// We use the start of today because the time defaults to 00:00 in timeCreated.
			var startOfToday time.Time = time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 0, 0, 0, 0, timeCreated.Location())

			// Only hours and not days is available in the API.
			if (startOfToday.Sub(timeCreated).Hours() / 24 /* hours */ > configs.MAX_AGE_TO_PROCESS_WIN_EVTS_IN_DAYS) {
				workflowScheduler.logger.Println("Event flagged as too old.")
				return uintptr(0)
			}
			
			callbackContext.provider = provider
			callbackContext.eventID = eventID
			callbackContext.eventRecordID = eventRecordID
			go runWorkflow(callbackContext)
		default:
			workflowScheduler.logger.Println("encountered error during callback: unsupported action code %x", uint16(Action))
	}
	return uintptr(0)
}

func (workflowScheduler *WorkflowScheduler) storeCallbackContext(context *CallbackContext) string {
	var uuidWithHyphen uuid.UUID = uuid.New()
    var uuid string = strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	workflowScheduler.guidToContext[uuid] = context
	fmt.Println(uuid)
	return uuid
}

func (workflowScheduler *WorkflowScheduler) addWinEventBasedSchedule(workflow workflows.InterfaceWorkflow, eventQueries []WinEventSubscribeQuery) {
	workflowScheduler.logger.Println("Workflow:", workflow)

	// Subscribe to the events that match the event queries specified.
	for _, eventQuery := range eventQueries {
		// Create the callback context for the subscription.
		var ctx *CallbackContext = &CallbackContext{workflow : workflow}
		var callbackContextUID string = workflowScheduler.storeCallbackContext(ctx)
		var CStringCallbackContextUID *C.char = C.CString(callbackContextUID)

		// Create a subscription for the event.
		subscriptionEventHandle, err := wevtapi.EvtSubscribe(
		wevtapi.EVT_HANDLE(win32.NULL),
		win32.HANDLE(win32.NULL),
		eventQuery.channel,
		eventQuery.query,
		wevtapi.EVT_HANDLE(win32.NULL),
		win32.PVOID(unsafe.Pointer(CStringCallbackContextUID)),
		workflowScheduler.SubscriptionCallback,
		wevtapi.EvtSubscribeToFutureEvents)
		if err != nil {
			workflowScheduler.logger.Println(err)
		}

		// Add the handle for the current subscription to the workflow scheduler.
		workflowScheduler.eventSubscriptionHandles = append(workflowScheduler.eventSubscriptionHandles, subscriptionEventHandle)
	}
	workflowScheduler.logger.Println("Workflow:", workflowScheduler.eventSubscriptionHandles)
}

func newWorkflowScheduler(schedulesJson []byte, workflowFactory WorkflowFactory) *WorkflowScheduler {
	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("WorkflowScheduler")
	var guidToContext map[string]*CallbackContext = make(map[string]*CallbackContext)
	var workflowScheduler *WorkflowScheduler = &WorkflowScheduler{logger: logger, guidToContext : guidToContext}

	// Parse the schedules JSON and add the schedules to the workflow scheduler.
	var schedules map[string][]ScheduleDescription
	json.Unmarshal([]byte(schedulesJson), &schedules)
	for _, schedule := range schedules["schedules"] {
		switch schedule.ScheduleType {
			case "on_win_event":
				var workflow workflows.InterfaceWorkflow = workflowFactory.constructWorkflow(schedule.Workflow)	
				var eventQueries []WinEventSubscribeQuery = parseEventSubscribeQueries(schedule.WinEventSubscribeQueries)			
				workflowScheduler.addWinEventBasedSchedule(workflow, eventQueries) 
			default:
				workflowScheduler.logger.Println("Given schedule type not supported.")
		}
	}
	return workflowScheduler
}

func parseEventSubscribeQueries(eventQueries [][]string) []WinEventSubscribeQuery {
	var parsedEventQueries []WinEventSubscribeQuery
	// Parse each of the event queries into the 'WinEventSubscribeQuery' type.
	for _, eventQuery := range eventQueries {
		var parsedEventQuery = eventQuery
		var channel string = parsedEventQuery[0]
		var query string = parsedEventQuery[1]
		parsedEventQueries = append(parsedEventQueries, WinEventSubscribeQuery{channel : channel, query : query})
	}
	return parsedEventQueries
}
