package main

import (
	"github.com/microsoft/BladeMonRT/test_configs"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"testing"
)

/** Instructions on how to run this test are in the README. */
func TestEndToEnd(t *testing.T) {
	// Skip the test if tests are run in short mode since this test will run until a keyboard interrupt is used.
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	runtime.GOMAXPROCS(1)

	workflowsJson, err := ioutil.ReadFile(test_configs.TEST_WORKFLOW_FILE)
	if err != nil {
		log.Fatal(err)
	}

	schedulesJson, err := ioutil.ReadFile(test_configs.TEST_SCHEDULE_FILE)
	if err != nil {
		log.Fatal(err)
	}

	var mainObj *Main = newMain()
	var workflowFactory WorkflowFactory = newWorkflowFactory(workflowsJson, NodeFactory{})
	mainObj.setupWorkflows(schedulesJson, workflowFactory)
	mainObj.logger.Println("Initialized main for test.")

	// Setup test such that test does not exit unless there is a keyboard interrupt.
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
}
