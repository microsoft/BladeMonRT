package windows_events

import (
	"github.com/microsoft/BladeMonRT/workflows"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
	"log"
	"github.com/microsoft/BladeMonRT/logging"
	"github.com/microsoft/BladeMonRT/nodes"
)

const (
	EVT_SUBSCRIBE_TO_FUTURE_EVENTS = 1
	EVT_SUBSCRIBE_START_AT_OLDEST_RECORD = 2
	EVT_SUBSCRIBE_ACTION_ERROR = 0
	EVT_SUBSCRIBE_ACTION_DELIVER = 1
	EVT_RENDER_EVENT_XML = 1
)

var (
	modwevtapi = windows.NewLazySystemDLL("wevtapi.dll")

	procEvtSubscribe = modwevtapi.NewProc("EvtSubscribe")
	procEvtRender    = modwevtapi.NewProc("EvtRender")
)

/** Class that holds information used in the callback function. */
type CallbackContext struct {
	Workflow workflows.InterfaceWorkflow
}

/** Type for a callback that is run when an event subscribed to is triggered. */
type SubscriptionCallback func(uintptr, uintptr, uintptr) uintptr

/** Utility class used to create subscriptions to windows events. */
type EventSubscriber struct {
	Logger *log.Logger
}

/** Class that defines the parameters used to subscribe to a windows event. */
type EventSubscription struct {
	Channel         string
	Query           string
	SubscribeMethod int
	Callback        SubscriptionCallback
	Context			CallbackContext
}

func NewEventSubscriber() EventSubscriber {
	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("EventSubscriber")
	return EventSubscriber{Logger : logger}
}

func (subscriber *EventSubscriber) CreateSubscription(evtSub *EventSubscription) {
	subscriber.Logger.Println("Creating subscription")

	winChannel, err := windows.UTF16PtrFromString(evtSub.Channel)
	if err != nil {
		subscriber.Logger.Println("bad channel name: %s", err)
		return
	}

	winQuery, err := windows.UTF16PtrFromString(evtSub.Query)
	if err != nil {
		subscriber.Logger.Println("bad query string: %s", err)
		return
	}

	// Subscribe to the windows event using the subscription object.
	handle, _, err := procEvtSubscribe.Call(
		0,
		0,
		uintptr(unsafe.Pointer(winChannel)),
		uintptr(unsafe.Pointer(winQuery)),
		0,
		uintptr(unsafe.Pointer(&evtSub.Context)),
		syscall.NewCallback(evtSub.Callback),
		uintptr(evtSub.SubscribeMethod),
	)

	if handle == 0 {
		subscriber.Logger.Println("failed to subscribe to events: %s", err)
		return
	}

	// TO DO: Use windows.Handle(handle) to add the windows handle to a list of handles
	// As is done in Python version.
}

func runWorkflow(context CallbackContext) {
	var workflowContext *nodes.WorkflowContext = nodes.NewWorkflowContext()
	var workflow workflows.InterfaceWorkflow = context.Workflow
	workflow.Run(workflow, workflowContext)
}

func (subscriber *EventSubscriber) SubscriptionCallback(action, context, event uintptr) uintptr {
	switch action {
		case EVT_SUBSCRIBE_ACTION_ERROR:
			subscriber.Logger.Println("Win event subscription failed.")
			return 0

		case EVT_SUBSCRIBE_ACTION_DELIVER:
			renderSpace := make([]uint16, 4096)
			bufferUsed := uint16(0)
			propertyCount := uint16(0)

			returnCode, _, err := procEvtRender.Call(0, event, EVT_RENDER_EVENT_XML, 4096, uintptr(unsafe.Pointer(&renderSpace[0])), uintptr(unsafe.Pointer(&bufferUsed)), uintptr(unsafe.Pointer(&propertyCount)))
			// TODO: use renderSpace to get the XML of the event and pass it to the callback

			if returnCode == 0 {
				subscriber.Logger.Println("failed to render event data: %s", err)
			} else {
				var callbackContext CallbackContext = *(*CallbackContext)(unsafe.Pointer(context))
				// Create a light-weight thread (goroutine) to run the workflow included in the callback context.
				go runWorkflow(callbackContext)
			}

		default:
			subscriber.Logger.Println("encountered error during callback: unsupported action code %x", uint16(action))
	}
	return 0
}