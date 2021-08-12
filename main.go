package main

import (
	"io/ioutil"
	"log"
	"runtime"
	"os"
	"syscall"
	"os/signal"
	"github.com/microsoft/BladeMonRT/logging"
	"github.com/microsoft/BladeMonRT/test_configs"
	"encoding/json"
	"github.com/microsoft/BladeMonRT/workflows"
)

type Main struct {
	logger          *log.Logger
	WorkflowScheduler	*WorkflowScheduler
}

func main() {
	// Set GOMAXPROCS such that all operations execute on a single thread.
	runtime.GOMAXPROCS(1)

	workflowsJson, err := ioutil.ReadFile(test_configs.TEST_WORKFLOW_FILE)
	if err != nil {
		log.Fatal(err)
	}
	var workflowFactory WorkflowFactory = newWorkflowFactory(workflowsJson, NodeFactory{})

	schedulesJson, err := ioutil.ReadFile(test_configs.TEST_SCHEDULE_FILE)
	if err != nil {
		log.Fatal(err)
	}

	var mainObj *Main = newMain()
	mainObj.setupWorkflows(schedulesJson,workflowFactory)
	mainObj.logger.Println("Initialized main.")

	// Setup main such that main does not exit unless there is a keyboard interrupt.
	quitChannel := make(chan os.Signal, 1)
    signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM) 
	<-quitChannel
}

func newMain() *Main{
	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("Main")
	return &Main{WorkflowScheduler: newWorkflowScheduler(), logger: logger}
}

func (main *Main) setupWorkflows(schedulesJson []byte, workflowFactory WorkflowFactory) {
	// Parse the schedules JSON and add the schedules to the workflow scheduler.
	var schedules map[string][]ScheduleDescription
	json.Unmarshal([]byte(schedulesJson), &schedules)
	for _, schedule := range schedules["schedules"] {
		switch schedule.ScheduleType {
			case "on_win_event":
				var workflow workflows.InterfaceWorkflow = workflowFactory.constructWorkflow(schedule.WorkflowName)	
				var eventQueries []WinEventSubscribeQuery = parseEventSubscribeQueries(schedule.WinEventSubscribeQueries)			
				main.WorkflowScheduler.addWinEventBasedSchedule(workflow, eventQueries) 
			default:
				main.WorkflowScheduler.logger.Println("Given schedule type not supported.")
		}
	}
}