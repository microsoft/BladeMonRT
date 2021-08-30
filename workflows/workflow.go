package workflows

import (
	"errors"
	"fmt"
	"github.com/microsoft/BladeMonRT/nodes"
	"github.com/microsoft/BladeMonRT/store"
	"github.com/microsoft/BladeMonRT/utils"
	"github.com/microsoft/BladeMonRT/configs"
	"log"
	"strconv"
)

/** Interface for defining execution sequence of nodes. */
type InterfaceWorkflow interface {
	AddNode(node nodes.InterfaceNode)
	Run(interfaceWorkflow InterfaceWorkflow, workflowContext *nodes.WorkflowContext)
	runVirt(workflowContext *nodes.WorkflowContext) error
	GetNodes() []nodes.InterfaceNode
}

/** Concrete type for defining execution sequence of nodes. */
type Workflow struct {
	Logger *log.Logger
	utils *utils.Utils
}

func (workflow *Workflow) Run(interfaceWorkflow InterfaceWorkflow, workflowContext *nodes.WorkflowContext) {
	workflow.Logger.Println("Running run method.")

	// Set the nodes in the workflow context to the nodes in this workflow.
	workflowContext.SetNodes(interfaceWorkflow.GetNodes())

	var err error = interfaceWorkflow.runVirt(workflowContext)
	// Handle errors thrown when running the workflow.
	if err != nil {
		workflow.Logger.Println(fmt.Sprintf("Workflow error: %s", err))
		return
	}

	if configs.ENABLE_BOOKMARK_FEATURE && workflowContext.QueryIncludesCondition {
		workflow.Logger.Println(fmt.Sprintf("Updating event record ID bookmark: %s to %d.", workflowContext.Query, workflowContext.EtwEvent.EventRecordID))
		workflow.updateEventRecordIdBookmark(workflowContext.BookmarkStore, workflowContext.Query, workflowContext.EtwEvent.EventRecordID)
	}
}

/** Processes a single node. Classes that implement InterfaceWorkflow should call this to process a node instead of node.process. */
func (workflow *Workflow) processNode(node nodes.InterfaceNode, workflowContext *nodes.WorkflowContext) (err error) {
	// Recover from panic during the processing of a node.
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("panic: %s", r))
		}
	}()
	// Return error returned by the processing of a node to the caller function.
	err = node.Process(node, workflowContext)
	return err
}

/** Updates the event record ID for the query if the current event record ID for the query is less than the given event record ID. */
func (workflow *Workflow) updateEventRecordIdBookmark(bookmarkStore store.PersistentKeyValueStoreInterface, query string, newEventRecordId int) {
	var currEventRecordID int = workflow.utils.GetEventRecordIdBookmark(bookmarkStore, query)
	if (newEventRecordId > currEventRecordID) {
		err := bookmarkStore.SetValue(query, strconv.Itoa(newEventRecordId))
		if (err != nil) {
			workflow.Logger.Println("Unable to update event record ID bookmark for query:", query)
		}
	}	
}
