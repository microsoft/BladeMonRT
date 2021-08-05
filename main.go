package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"
	"os"
	"syscall"
	"os/signal"
	"github.com/microsoft/BladeMonRT/logging"
	"github.com/microsoft/BladeMonRT/configs"
)

type Main struct {
	workflowFactory *WorkflowFactory
	logger          *log.Logger
}

func main() {
	var mainObj Main = NewMain()
	mainObj.logger.Println("Initialized main.")

	// Setup main such that main does not exit unless there is a keyboard interrupt.
	go forever()
	quitChannel := make(chan os.Signal, 1)
    signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM) 
	<-quitChannel
}

func NewMain() Main {
	workflowsJson, err := ioutil.ReadFile(configs.WORKFLOW_FILE)
	if err != nil {
		log.Fatal(err)
	}
	var workflowFactory WorkflowFactory = newWorkflowFactory(workflowsJson, NodeFactory{})

	schedulesJson, err := ioutil.ReadFile(configs.SCHEDULE_FILE)
	if err != nil {
		log.Fatal(err)
	}
	var workflowScheduler WorkflowScheduler = newWorkflowScheduler(schedulesJson, workflowFactory)
	fmt.Println(workflowScheduler) // TODO: Remove print statement.

	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("Main")
	return Main{workflowFactory: &workflowFactory, logger: logger}
}


func forever() {
    for {
        time.Sleep(time.Second)
    }
}