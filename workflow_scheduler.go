package main

import (
	"encoding/json"
	"github.com/microsoft/BladeMonRT/workflows"
	"github.com/microsoft/BladeMonRT/logging"
	"github.com/microsoft/BladeMonRT/configs"
	"log"
	winEvents "github.com/microsoft/BladeMonRT/windows_events"
	"golang.org/x/sys/windows"
	"regexp"
	"strings"
	"fmt"
	"strconv"
	win32 "github.com/0xrawsec/golang-win32/win32"
	wevtapi "github.com/0xrawsec/golang-win32/win32/wevtapi"
	"unsafe"
)

/** Class for scheduling workflows. */
type WorkflowScheduler struct {
	schedules []interface{}
	logger *log.Logger
	subscriber winEvents.EventSubscriber
	eventSubscriptionHandles []windows.Handle
	queryToEventRecordIdBookmark map[string]int
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

func (workflowScheduler *WorkflowScheduler) getEventRecordIdBookmark(query string) int {
	if (!configs.ENABLE_BOOKMARK_FEATURE) {
		return 0
	}

	if eventRecordIdBookmark, ok := workflowScheduler.queryToEventRecordIdBookmark[query]; ok {
		return eventRecordIdBookmark
	} else {
		return 0
	}
}

func (workflowScheduler *WorkflowScheduler) updateEventRecordIdBookmark(query string, newEventRecordId int) {
	workflowScheduler.queryToEventRecordIdBookmark[query] = newEventRecordId
}


func (workflowScheduler *WorkflowScheduler) addWinEventBasedSchedule(workflow workflows.InterfaceWorkflow, eventQueries []WinEventSubscribeQuery) {
	workflowScheduler.logger.Println("Workflow:", workflow)

	// Subscribe to the events that match the event queries specified.
	for _, eventQuery := range eventQueries {
		// Decide whether to subscribe to future events or start at the oldest record.
		
		var subscribeToFutureEvents bool = true
		var queryText = eventQuery.query
		queryIncludesCondition, err := regexp.MatchString(".*{condition}.*", eventQuery.query)
		if (err != nil) {
			return
		}
		if (queryIncludesCondition) {
			var eventRecordIdBookmark int = workflowScheduler.getEventRecordIdBookmark(eventQuery.query)
			if (eventRecordIdBookmark != 0) {
				subscribeToFutureEvents = false
			}
			workflowScheduler.logger.Println(eventRecordIdBookmark)
			queryText = strings.Replace(eventQuery.query, "{condition}", strconv.Itoa(eventRecordIdBookmark), -1)
		}
		workflowScheduler.logger.Println(fmt.Sprintf("Constructed queryText: %s; subscribeToFutureEvents: %t", queryText, subscribeToFutureEvents))
		
		var subscribeMethod win32.DWORD = wevtapi.EvtSubscribeToFutureEvents
		if (subscribeToFutureEvents) {
			subscribeMethod = wevtapi.EvtSubscribeStartAtOldestRecord
		}
		
		ctx := &winEvents.CallbackContext{Workflow : workflow}
		workflowScheduler.logger.Println(win32.NULL)
		_, err = wevtapi.EvtSubscribe(
		wevtapi.EVT_HANDLE(win32.NULL),

		win32.HANDLE(win32.NULL),
		eventQuery.channel,
		queryText,
		wevtapi.EVT_HANDLE(win32.NULL),
		win32.PVOID(unsafe.Pointer(ctx)),
		workflowScheduler.subscriber.SubscriptionCallback,
		subscribeMethod)

		if err != nil {
			workflowScheduler.logger.Println(err)
		}

		// Add the handle for the current subscription to the workflow scheduler.
		// var subscriptionEventHandle []windows.Handle = workflowScheduler.subscriber.CreateSubscription(eventSubscription)
		//workflowScheduler.eventSubscriptionHandles = append(workflowScheduler.eventSubscriptionHandles, subscriptionEventHandle...)
	}
	workflowScheduler.logger.Println("Workflow:", workflowScheduler.eventSubscriptionHandles)
}

func newWorkflowScheduler(schedulesJson []byte, workflowFactory WorkflowFactory) *WorkflowScheduler {
	var subscriber winEvents.EventSubscriber = winEvents.NewEventSubscriber()
	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("WorkflowScheduler")
	var workflowScheduler *WorkflowScheduler = &WorkflowScheduler{subscriber : subscriber, logger: logger}

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