package nodes

import (
	"github.com/microsoft/BladeMonRT/utils"
	"github.com/microsoft/BladeMonRT/store"
	"log"
)

/** Class that stores information about the current state of a running workflow. */
type WorkflowContext struct {
	nodes        []InterfaceNode
	nodeToResult map[InterfaceNode]interface{}
	Seed         string
	EtwEvent     utils.EtwEvent
	Query 			string
	BookmarkStore	store.PersistentKeyValueStoreInterface
}

func NewWorkflowContext() *WorkflowContext {
	var nodeToResult map[InterfaceNode]interface{} = make(map[InterfaceNode]interface{})
	var workflowContext WorkflowContext = WorkflowContext{nodeToResult: nodeToResult}
	return &workflowContext
}

func (workflowContext *WorkflowContext) SetNodes(nodes []InterfaceNode) {
	workflowContext.nodes = nodes
}

func (workflowContext *WorkflowContext) GetNodes() []InterfaceNode {
	return workflowContext.nodes
}

/** Interface for defining unit of work to be processed by event loop. Classes that implement InterfaceNode can provide their own constructor. */
type InterfaceNode interface {
	Process(interfaceNode InterfaceNode, workflowContext *WorkflowContext) error
	ProcessVirt(workflowContext *WorkflowContext) error
	SaveResult(interfaceNode InterfaceNode, workflowContext *WorkflowContext, result interface{})
	GetResult(interfaceNode InterfaceNode, workflowContext *WorkflowContext) interface{}
	GetPredecessorNodes(interfaceNode InterfaceNode, workflowContext *WorkflowContext) []InterfaceNode
	GetPredecessorNodesResults(interfaceNode InterfaceNode, workflowContext *WorkflowContext) []interface{}
}

/** Concrete type for defining unit of work to be processed by event loop. */
type Node struct {
	Logger *log.Logger
}

func (node *Node) Process(interfaceNode InterfaceNode, workflowContext *WorkflowContext) error {
	// TODO: Add logging.
	return interfaceNode.ProcessVirt(workflowContext)
}

func (node *Node) SaveResult(interfaceNode InterfaceNode, workflowContext *WorkflowContext, result interface{}) {
	node.Logger.Println("Running SaveResult method.")
	workflowContext.nodeToResult[interfaceNode] = result
}

func (node *Node) GetResult(interfaceNode InterfaceNode, workflowContext *WorkflowContext) interface{} {
	node.Logger.Println("Running GetResult method.")
	return workflowContext.nodeToResult[interfaceNode]
}

func (node *Node) GetPredecessorNodes(interfaceNode InterfaceNode, workflowContext *WorkflowContext) []InterfaceNode {
	var predecessorNodes []InterfaceNode
	var nodes []InterfaceNode = workflowContext.GetNodes()
	for _, currNode := range nodes {
		if interfaceNode == currNode {
			break
		}
		predecessorNodes = append(predecessorNodes, currNode)
	}
	return predecessorNodes
}

func (node *Node) GetPredecessorNodesResults(interfaceNode InterfaceNode, workflowContext *WorkflowContext) []interface{} {
	var predecessorNodesResults []interface{}
	var predecessorNodes []InterfaceNode = interfaceNode.GetPredecessorNodes(interfaceNode, workflowContext)
	for _, predecessorNode := range predecessorNodes {
		predecessorNodesResults = append(predecessorNodesResults, workflowContext.nodeToResult[predecessorNode])
	}
	return predecessorNodesResults
}
