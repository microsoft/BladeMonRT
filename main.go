package main

import (
	"io/ioutil"
	"log"
	"os"
	"github.com/microsoft/BladeMonRT/workflows"
)

type Main struct {
    WorkflowManager *WorkflowManager
	Logger *log.Logger
}

func main() {
	var mainObj Main = NewMain()

	var workflow workflows.InterfaceWorkflow = mainObj.WorkflowManager.constructWorkflow("dummy_workflow")
	workflow.RunVirt()
	mainObj.Logger.Println("The result is: ", workflow.GetResult())
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
	var workflowManager WorkflowManager = newWorkflowManager(workflowsJson)

	file, err := os.OpenFile(logging_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatal(err)
    }
    var logger *log.Logger = log.New(file, "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
	
	return Main{WorkflowManager : &workflowManager, Logger : logger}
}