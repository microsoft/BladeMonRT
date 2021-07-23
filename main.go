package main

import (
	"io/ioutil"
	"log"
	"fmt"
)

func main() {
	const (
		workflow_file = "configs/workflows.json"
	)

	workflowsJson, err := ioutil.ReadFile(workflow_file)
	if err != nil {
		log.Fatal(err)
	}

	var workflowManager WorkflowManager = newWorkflowManager(workflowsJson)
	workflow := workflowManager.constructWorkflow("dummy_workflow")
	workflow.RunVirt()
	fmt.Println("The result is" + workflow.GetResult())
}