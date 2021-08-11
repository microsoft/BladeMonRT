package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"os"
	"syscall"
	"os/signal"
	"github.com/microsoft/BladeMonRT/test_configs"
	"testing"
)

/** Instructions on how to run this test are in the README. */
func TestEndToEnd(t *testing.T) {
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
	var workflowScheduler *WorkflowScheduler = newWorkflowScheduler(schedulesJson, workflowFactory)
	fmt.Println(workflowScheduler)

	// Setup main such that main does not exit unless there is a keyboard interrupt.
	quitChannel := make(chan os.Signal, 1)
    signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM) 
	<-quitChannel
}