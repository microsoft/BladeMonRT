package main

import (
	"encoding/json"
	"github.com/microsoft/BladeMonRT/workflows"
	"github.com/microsoft/BladeMonRT/logging"
	"github.com/microsoft/BladeMonRT/nodes"
	"log"
	winEvents "github.com/microsoft/BladeMonRT/windows_events"
)

/** Class for scheduling workflows. */
type WorkflowScheduler struct {
	schedules []interface{}
	logger *log.Logger
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

/** Class that represents a query for subscribing to a windows event. */
type WinEventSubscribeQuery struct {
	channel string
	query string
}

func workflowCallback(context winEvents.CallbackContext) {
	var workflowContext *nodes.WorkflowContext = nodes.NewWorkflowContext()
	var workflow workflows.InterfaceWorkflow = context.Workflow
	workflow.Run(workflow, workflowContext)
}

func (workflowScheduler *WorkflowScheduler) addWinEventBasedSchedule(workflow workflows.InterfaceWorkflow, eventQueries []WinEventSubscribeQuery) {
	workflowScheduler.logger.Println("Workflow:", workflow)

	// Subscribe to the events that match the event queries specified.
	for _, eventQuery := range eventQueries {
		workflowScheduler.logger.Println("Channel:", eventQuery.channel)
		workflowScheduler.logger.Println("Query:", eventQuery.query)
		// TO DO: Subscribe to an event using the gowinlog library

		var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("EventSubscription")
		var eventSubscription *winEvents.EventSubscription = &winEvents.EventSubscription{
			Logger : logger,
			Channel:        eventQuery.channel,
			Query:          eventQuery.query,
			SubscribeMethod: winEvents.EvtSubscribeToFutureEvents,
			Callback:        workflowCallback,
			Context:         winEvents.CallbackContext{Workflow : workflow},
		}

		eventSubscription.CreateSubscription()
	}
}

func newWorkflowScheduler(schedulesJson []byte, workflowFactory WorkflowFactory) WorkflowScheduler {
	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("WorkflowScheduler")
	var workflowScheduler *WorkflowScheduler = &WorkflowScheduler{logger: logger}

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
				panic("Given schedule type not supported.")
		}
	}
	return WorkflowScheduler{}
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