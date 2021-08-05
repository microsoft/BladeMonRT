package windows_events

import "C"
import (
	"github.com/microsoft/BladeMonRT/workflows"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
	"log"
)

const (
	// EvtSubscribeToFutureEvents instructs the
	// subscription to only receive events that occur
	// after the subscription has been made
	EvtSubscribeToFutureEvents = 1

	// EvtSubscribeStartAtOldestRecord instructs the	
	// subscription to receive all events (past and future)
	// that match the query
	EvtSubscribeStartAtOldestRecord = 2

	// evtSubscribeActionError defines a action
	// code that may be received by the winAPICallback.
	// ActionError defines that an internal logger.Println occurred
	// while obtaining an event for the callback
	evtSubscribeActionError = 0

	// evtSubscribeActionDeliver defines a action
	// code that may be received by the winAPICallback.
	// ActionDeliver defines that the internal API was
	// successful in obtaining an event that matched
	// the subscription query
	evtSubscribeActionDeliver = 1

	// evtRenderEventXML instructs procEvtRender
	// to render the event details as a XML string
	evtRenderEventXML = 1
)

var (
	modwevtapi = windows.NewLazySystemDLL("wevtapi.dll")

	procEvtSubscribe = modwevtapi.NewProc("EvtSubscribe")
	procEvtRender    = modwevtapi.NewProc("EvtRender")
	procEvtClose     = modwevtapi.NewProc("EvtClose")
)

type CallbackContext struct {
	Workflow workflows.InterfaceWorkflow
}

type EventCallback func(context CallbackContext)


// EventSubscription is a subscription to
// Windows Events, it defines details about the
// subscription including the channel and query
type EventSubscription struct {
	Logger *log.Logger
	Channel         string
	Query           string
	SubscribeMethod int
	Callback        EventCallback
	Context			CallbackContext
	winAPIHandle windows.Handle
}

// Create will setup an event subscription in the
// windows kernel with the provided channel and
// event query
func (evtSub *EventSubscription) CreateSubscription() {
	evtSub.Logger.Println("Creating subscription")
	if evtSub.winAPIHandle != 0 {
		evtSub.Logger.Println("subscription already created in kernel")
		return
	}

	winChannel, err := windows.UTF16PtrFromString(evtSub.Channel)
	if err != nil {
		evtSub.Logger.Println("bad channel name: %s", err)
		return
	}

	winQuery, err := windows.UTF16PtrFromString(evtSub.Query)
	if err != nil {
		evtSub.Logger.Println("bad query string: %s", err)
		return
	}

	handle, _, err := procEvtSubscribe.Call(
		0,
		0,
		uintptr(unsafe.Pointer(winChannel)),
		uintptr(unsafe.Pointer(winQuery)),
		0,
		uintptr(unsafe.Pointer(&evtSub.Context)),
		syscall.NewCallback(evtSub.winAPICallback),
		uintptr(evtSub.SubscribeMethod),
	)

	if handle == 0 {
		evtSub.Logger.Println("failed to subscribe to events: %s", err)
		return
	}

	evtSub.winAPIHandle = windows.Handle(handle)
}

// winAPICallback receives the callback from the windows
// kernel when an event matching the query and channel is
// received. It will query the kernel to get the event rendered
// as a XML string, the XML string is then unmarshaled to an
// `Event` and the custom callback invoked
func (evtSub *EventSubscription) winAPICallback(action, userContext, event uintptr) uintptr {
	switch action {
	case evtSubscribeActionError:
		evtSub.Logger.Println("Win event subscription failed.")
		return 0

	case evtSubscribeActionDeliver:
		renderSpace := make([]uint16, 4096)
		bufferUsed := uint16(0)
		propertyCount := uint16(0)

		returnCode, _, err := procEvtRender.Call(0, event, evtRenderEventXML, 4096, uintptr(unsafe.Pointer(&renderSpace[0])), uintptr(unsafe.Pointer(&bufferUsed)), uintptr(unsafe.Pointer(&propertyCount)))
		// TODO: use renderSpace to get the XML of the event and pass it to the callback

		if returnCode == 0 {
			evtSub.Logger.Println("failed to render event data: %s", err)
		} else {
			var callbackContext CallbackContext = *(*CallbackContext)(unsafe.Pointer(userContext))
			go evtSub.Callback(callbackContext)
		}

	default:
		evtSub.Logger.Println("encountered error during callback: unsupported action code %x", uint16(action))
	}

	return 0
}