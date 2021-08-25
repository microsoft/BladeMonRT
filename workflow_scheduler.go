package main

import (
	"C"
	"fmt"
	win32 "github.com/0xrawsec/golang-win32/win32"
	wevtapi "github.com/0xrawsec/golang-win32/win32/wevtapi"
	"github.com/google/uuid"
	"github.com/microsoft/BladeMonRT/logging"
	"github.com/microsoft/BladeMonRT/nodes"
	"github.com/microsoft/BladeMonRT/utils"
	"github.com/microsoft/BladeMonRT/workflows"
	"log"
	"unsafe"
	"regexp"
	"time"
	"github.com/microsoft/BladeMonRT/configs"
	"github.com/microsoft/BladeMonRT/store"
	"strings"
	"strconv"
	"errors"
)

/** Interface for scheduling workflows. */
type WorkflowSchedulerInterface interface {
	addWinEventBasedSchedule(workflow workflows.InterfaceWorkflow, eventQueries []WinEventSubscribeQuery)
}

/** Class for scheduling workflows. */
type WorkflowScheduler struct {
	logger                   *log.Logger
	eventSubscriptionHandles []wevtapi.EVT_HANDLE
	guidToContext            map[string]*CallbackContext
	bookmarkStore 			 store.PersistentKeyValueStoreInterface
	utils                    utils.UtilsInterface
}

/** Class that represents a query for subscribing to a windows event. */
type WinEventSubscribeQuery struct {
	channel string
	query   string
}

/** Class for the schedule description in the JSON. */
type ScheduleDescription struct {
	Name         string `json:"name"`
	ScheduleType string `json:"schedule_type"`
	WorkflowName string `json:"workflow"`
	Enable       bool   `json:"enable"`

	// Attributes specific to events of type 'on_win_event'.
	WinEventSubscribeQueries [][]string `json:"win_event_subscribe_queries"`
}

/** Class that holds information used in the subscription callback function. */
type CallbackContext struct {
	query 			string
	queryIncludesCondition bool
	bookmarkStore	store.PersistentKeyValueStoreInterface
	workflow        workflows.InterfaceWorkflow
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
		var nowTime time.Time = time.Now()

		// We use the start of today because the time defaults to 00:00 in event.TimeCreated.
		var startOfToday time.Time = time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 0, 0, 0, 0, event.TimeCreated.Location())

		// Check if the event is too old to process.
		var age float64 = startOfToday.Sub(event.TimeCreated).Hours() / float64(24)
		if (age > configs.MAX_AGE_TO_PROCESS_WIN_EVTS_IN_DAYS) {
			workflowScheduler.logger.Println("Event flagged as too old. Age:", age)
			return uintptr(0)
		}

		if (configs.ENABLE_BOOKMARK_FEATURE && callbackContext.queryIncludesCondition) {
			workflowScheduler.logger.Println(fmt.Sprintf("Updating event record ID bookmark: %s to %d.", callbackContext.query, event.EventRecordID))
			workflowScheduler.updateEventRecordIdBookmark(callbackContext.query, event.EventRecordID)
	
		}

		callbackContext.workflowContext = nodes.NewWorkflowContext()
		callbackContext.workflowContext.Query = callbackContext.query
		callbackContext.workflowContext.BookmarkStore = callbackContext.bookmarkStore
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
		ctx, err := workflowScheduler.decideCallbackContext(eventQuery, workflow)
		if (err != nil) {
			workflowScheduler.logger.Println("Error in deciding callback context:s", err)
			return
		}
		var callbackContextUID string = workflowScheduler.storeCallbackContext(ctx)
		var CStringCallbackContextUID *C.char = C.CString(callbackContextUID)

		// Decide whether to subscribe to future events or start at the oldest record. 
		var subscribeToFutureEvents bool = true
		var queryText string
        if (ctx.queryIncludesCondition) {
            var eventRecordIdBookmark int = workflowScheduler.getEventRecordIdBookmark(eventQuery.query)
            if (eventRecordIdBookmark != 0) {
                subscribeToFutureEvents = false
            }
            queryText = strings.Replace(eventQuery.query, "{condition}", strconv.Itoa(eventRecordIdBookmark), -1)
        }

        workflowScheduler.logger.Println(fmt.Sprintf("Constructed queryText: %s; subscribeToFutureEvents: %t", queryText, subscribeToFutureEvents))
        var subscribeMethod win32.DWORD = wevtapi.EvtSubscribeStartAtOldestRecord
        if (subscribeToFutureEvents) {
             subscribeMethod = wevtapi.EvtSubscribeToFutureEvents
        }

		// Create a subscription for the event.
		subscriptionEventHandle, err := wevtapi.EvtSubscribe(
			wevtapi.EVT_HANDLE(win32.NULL),
			win32.HANDLE(win32.NULL),
			eventQuery.channel,
			queryText,
			wevtapi.EVT_HANDLE(win32.NULL),
			win32.PVOID(unsafe.Pointer(CStringCallbackContextUID)),
			workflowScheduler.SubscriptionCallback,
			subscribeMethod)
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
	var utils utils.UtilsInterface = utils.NewUtils()
	
	var bookmarkStore *store.PersistentKeyValueStore
	bookmarkStore, err := store.NewPersistentKeyValueStore(configs.BOOKMARK_DATABASE_FILE, configs.BOOKMARK_DATABASE_TABLE_NAME)
	bookmarkStore.InitTable()

	if err != nil {
		panic("Unable to create bookmark store.")
	}
	
	var workflowScheduler *WorkflowScheduler = &WorkflowScheduler{logger: logger, guidToContext: guidToContext, bookmarkStore: bookmarkStore, utils: utils}
	return workflowScheduler
}

func parseEventSubscribeQueries(eventQueries [][]string) []WinEventSubscribeQuery {
	var parsedEventQueries []WinEventSubscribeQuery
	// Parse each of the event queries into the 'WinEventSubscribeQuery' type.
	for _, eventQuery := range eventQueries {
		var parsedEventQuery = eventQuery
		var channel string = parsedEventQuery[0]
		var query string = parsedEventQuery[1]
		parsedEventQueries = append(parsedEventQueries, WinEventSubscribeQuery{channel: channel, query: query})
	}
	return parsedEventQueries
}

func (workflowScheduler *WorkflowScheduler) getEventRecordIdBookmark(query string) int {
	if (!configs.ENABLE_BOOKMARK_FEATURE) {
		return 0
	}

	stringEventRecordId, err := workflowScheduler.bookmarkStore.GetValue(query)
	if err != nil {
		return 0
	}
	eventRecordId, err := strconv.Atoi(stringEventRecordId)
	return eventRecordId
}

func (workflowScheduler *WorkflowScheduler) decideCallbackContext(eventQuery WinEventSubscribeQuery, workflow workflows.InterfaceWorkflow) (*CallbackContext, error) {	
	var query string = eventQuery.query

	// Check if the query contains a condition.
	queryIncludesCondition, err := regexp.MatchString(".*{condition}.*", query)
	if (err != nil) {
		return nil, errors.New(fmt.Sprintf("Failed to check if query %s contains condition.", query))
	}
	var bookmarkStore store.PersistentKeyValueStoreInterface = nil
	if (queryIncludesCondition) {
		bookmarkStore = workflowScheduler.bookmarkStore
	}

	var ctx *CallbackContext = &CallbackContext{query: query, workflow: workflow, queryIncludesCondition: queryIncludesCondition, bookmarkStore: bookmarkStore}
	return ctx, nil
}

func (workflowScheduler *WorkflowScheduler) updateEventRecordIdBookmark(query string, newEventRecordId int) {
	workflowScheduler.bookmarkStore.SetValue(query, strconv.Itoa(newEventRecordId))
}
