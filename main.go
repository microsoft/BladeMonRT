package main

import (
	"io/ioutil"
	"log"
	"os"
	"github.com/microsoft/BladeMonRT/workflows"
	"github.com/microsoft/BladeMonRT/nodes"
	"fmt"
)

type Main struct {
    WorkflowFactory *WorkflowFactory
	Logger *log.Logger
}

func main() {
	var mainObj Main = NewMain()
	var workflow workflows.InterfaceWorkflow = mainObj.WorkflowFactory.constructWorkflow("dummy_workflow")
	workflow.Run(workflow)
	
	var workflowContext *nodes.WorkflowContext = workflow.GetWorkflowContext()
	for index, node := range workflow.GetNodes() {
		mainObj.Logger.Println(fmt.Sprintf("Result for node index %d=%s", index, node.GetResult(node, workflowContext).(string)))
	}
}	

func NewMain() Main {
	const (
		workflow_file = "configs/workflows.json"
		logging_file = "log"
	)

	workflowsJson, err := ioutil.ReadFile(workflow_file)
	if err != nil {
		log.Fatal(err)
	}
	var WorkflowFactory WorkflowFactory = newWorkflowFactory(workflowsJson)

	file, err := os.OpenFile(logging_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatal(err)
    }
    var logger *log.Logger = log.New(file, "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
	
	return Main{WorkflowFactory : &WorkflowFactory, Logger : logger}
}