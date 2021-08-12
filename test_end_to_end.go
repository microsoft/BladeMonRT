package main

import (
	"io/ioutil"
	"log"
	"runtime"
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

	schedulesJson, err := ioutil.ReadFile(test_configs.TEST_SCHEDULE_FILE)
	if err != nil {
		log.Fatal(err)
	}
	var mainObj *Main = newMain()
	var workflowFactory WorkflowFactory = newWorkflowFactory(workflowsJson, NodeFactory{})
	mainObj.setupWorkflows(schedulesJson, workflowFactory)
	mainObj.logger.Println("Initialized main for test.")
}