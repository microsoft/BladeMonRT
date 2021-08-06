package windows_events

import (
	"C"
	"github.com/microsoft/BladeMonRT/workflows"
	"github.com/microsoft/BladeMonRT/configs"
    win32 "github.com/0xrawsec/golang-win32/win32"
	wevtapi "github.com/0xrawsec/golang-win32/win32/wevtapi"
	"time"
	"log"
	"github.com/microsoft/BladeMonRT/logging"
	"github.com/microsoft/BladeMonRT/nodes"
	"github.com/microsoft/BladeMonRT/utils"
	"fmt"
)

/** Class that holds information used in the callback function. */
type CallbackContext struct {
	// QueryToEventRecordIdBookmark map[string]int
	Workflow workflows.InterfaceWorkflow
}

/** Utility class used to create subscriptions to windows events. */
type EventSubscriber struct {
	logger *log.Logger
}

func NewEventSubscriber() EventSubscriber {
	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("EventSubscriber")
	return EventSubscriber{logger : logger}
}

func runWorkflow(context CallbackContext) {
	fmt.Println("run workflow")
	var workflowContext *nodes.WorkflowContext = nodes.NewWorkflowContext()
	var workflow workflows.InterfaceWorkflow = context.Workflow
	workflow.Run(workflow, workflowContext)
}

func (eventSubscriber *EventSubscriber) SubscriptionCallback(Action wevtapi.EVT_SUBSCRIBE_NOTIFY_ACTION, UserContext win32.PVOID, Event wevtapi.EVT_HANDLE) uintptr {
	// ctx := (*CallbackContext)(unsafe.Pointer(UserContext))
	// fmt.Println(ctx)

	switch Action {
		case wevtapi.EvtSubscribeActionError:
			eventSubscriber.logger.Println("Win event subscription failed.")
			return 0
		case wevtapi.EvtSubscribeActionDeliver:
			data, err := wevtapi.EvtRenderXML(Event)
			if err != nil {
				eventSubscriber.logger.Println("Error converting event to XML: %s", err)
			}
			dataUTF8 := win32.UTF16BytesToString(data)

			provider, eventID, timeCreated, eventRecordID := utils.NewUtils().ParseEventXML(dataUTF8)
			// TODO: Remove print statements.
			eventSubscriber.logger.Println("Provider:", provider)
			eventSubscriber.logger.Println("eventID:", eventID)
			eventSubscriber.logger.Println("eventRecordID:", eventRecordID)
			nowTime := time.Now()

			// We use the start of today because the time defaults to 00:00 in timeCreated.
			startOfToday := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 0, 0, 0, 0, timeCreated.Location())
			eventSubscriber.logger.Println(startOfToday)

			// Only hours and not days is available in the API.
			if (startOfToday.Sub(timeCreated).Hours() / 24 /* hours */ > configs.MAX_AGE_TO_PROCESS_WIN_EVTS_IN_DAYS) {
				eventSubscriber.logger.Println("Event flagged as too old.")
				return uintptr(0)
			}
			
			// go runWorkflow(*ctx)
		default:
			eventSubscriber.logger.Println("encountered error during callback: unsupported action code %x", uint16(Action))
	}
	return uintptr(0)
}