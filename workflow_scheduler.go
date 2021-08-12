package main

import (
	"C"
	"github.com/microsoft/BladeMonRT/workflows"
	"github.com/microsoft/BladeMonRT/nodes"
	"github.com/microsoft/BladeMonRT/logging"
	"log"
	win32 "github.com/0xrawsec/golang-win32/win32"
	wevtapi "github.com/0xrawsec/golang-win32/win32/wevtapi"
	"unsafe"
	"github.com/microsoft/BladeMonRT/utils"
	"github.com/google/uuid"
	"fmt"
)

/** Class for scheduling workflows. */
type WorkflowScheduler struct {
	logger *log.Logger
	eventSubscriptionHandles []wevtapi.EVT_HANDLE
	guidToContext map[string]*CallbackContext
	utils utils.UtilsInterface
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
    WorkflowName string `json:"workflow"`
    Enable bool `json:"enable"`

	// Attributes specific to events of type 'on_win_event'.
	WinEventSubscribeQueries [][]string `json:"win_event_subscribe_queries"`
}

/** Class that holds information used in the subscription callback function. */
type CallbackContext struct {
	workflow workflows.InterfaceWorkflow
	workflowContext *nodes.WorkflowContext
}

func (workflowScheduler *WorkflowScheduler) SubscriptionCallback(Action wevtapi.EVT_SUBSCRIBE_NOTIFY_ACTION, UserContext win32.PVOID, Event wevtapi.EVT_HANDLE) uintptr {
	var CStringGuid *C.char = (*C.char)(unsafe.Pointer(UserContext))
	var guid string = C.GoString(CStringGuid)
	var callbackContext *CallbackContext = workflowScheduler.guidToContext[guid]

	switch Action {
		case wevtapi.EvtSubscribeActionError:
			workflowScheduler.logger.Println("Win event subscription failed.")
			return uintptr(0)
		case wevtapi.EvtSubscribeActionDeliver:
			Utf16EventXml, err := wevtapi.EvtRenderXML(Event)
			if err != nil {
				workflowScheduler.logger.Println("Error converting event to XML:", err)
			}
			var eventXML string = win32.UTF16BytesToString(Utf16EventXml)
			var event utils.EtwEvent = workflowScheduler.utils.ParseEventXML(eventXML)

			callbackContext.workflowContext = nodes.NewWorkflowContext()
			callbackContext.workflowContext.Seed = eventXML
			callbackContext.workflowContext.EtwEvent = event

			// Create a goroutine to run the workflow included in the callback context.
			go callbackContext.workflow.Run(callbackContext.workflow, callbackContext.workflowContext)
		default:
			workflowScheduler.logger.Println(fmt.Sprintf("encountered error during callback: unsupported action code %x", uint16(Action)))
	}
	return uintptr(0)
}

func (workflowScheduler *WorkflowScheduler) storeCallbackContext(context *CallbackContext) string {
	var uuid string = uuid.New().String()
	workflowScheduler.guidToContext[uuid] = context
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
			return
		}

		// Add the handle for the current subscription to the workflow scheduler.
		workflowScheduler.eventSubscriptionHandles = append(workflowScheduler.eventSubscriptionHandles, subscriptionEventHandle)
	}
}

func newWorkflowScheduler() *WorkflowScheduler {
	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("WorkflowScheduler")
	var guidToContext map[string]*CallbackContext = make(map[string]*CallbackContext)
	var workflowScheduler *WorkflowScheduler = &WorkflowScheduler{logger: logger, guidToContext : guidToContext, utils : utils.NewUtils()}

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
	var workflowScheduler *WorkflowScheduler = &WorkflowScheduler{logger: logger, guidToContext : guidToContext, utils: utils.NewUtils()}
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
