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
	EvtSubscribeToFutureEvents = 1
	EvtSubscribeStartAtOldestRecord = 2
	evtSubscribeActionError = 0
	evtSubscribeActionDeliver = 1
	evtRenderEventXML = 1
)

var (
	modwevtapi = windows.NewLazySystemDLL("wevtapi.dll")

	procEvtSubscribe = modwevtapi.NewProc("EvtSubscribe")
	procEvtRender    = modwevtapi.NewProc("EvtRender")
	procEvtClose     = modwevtapi.NewProc("EvtClose")
)

/** Class that holds information used in the callback function. */
type CallbackContext struct {
	Workflow workflows.InterfaceWorkflow
}

/** Type for a callback that is run when an event subscribed to is triggered. */
type EventCallback func(context CallbackContext)


/** Class that defines the parameters used to subscribe to a windows event. */
type EventSubscription struct {
	Logger *log.Logger
	Channel         string
	Query           string
	SubscribeMethod int
	Callback        EventCallback
	Context			CallbackContext
	winAPIHandle windows.Handle
}

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

	// Subscribe to the windows event using the subscription object.
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