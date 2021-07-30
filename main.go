package main

import (
	"io/ioutil"
	"github.com/microsoft/BladeMonRT/workflows"
	"github.com/microsoft/BladeMonRT/nodes"
	"log"
	"github.com/microsoft/BladeMonRT/logging"
	"fmt"
)

type Main struct {
    workflowFactory *WorkflowFactory
	logger *log.Logger
}

func main() {
	var mainObj Main = NewMain()

	var workflow workflows.InterfaceWorkflow = mainObj.workflowFactory.constructWorkflow("dummy_workflow")
	var workflowContext *nodes.WorkflowContext = nodes.NewWorkflowContext()
	workflow.Run(workflow, workflowContext)
	
	for index, node := range workflow.GetNodes() {
		mainObj.logger.Println(fmt.Sprintf("Result for node index %d=%s", index, node.GetResult(node, workflowContext).(string)))
	}
}	

func NewMain() Main {
	workflowsJson, err := ioutil.ReadFile(workflow_file)
	if err != nil {
		log.Fatal(err)
	}
	var workflowFactory WorkflowFactory = newWorkflowFactory(workflowsJson, NodeFactory{})
    var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("Main")
	
	return Main{workflowFactory : &workflowFactory, logger : logger}
}