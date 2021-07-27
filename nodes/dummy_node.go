package nodes

import (
	"fmt"
)

/** Node that has the concatenation of its predecessors' results and a hard-coded value for its result. */
type DummyNode struct {
	Node
}

func NewDummyNode() *DummyNode {
	return &DummyNode{}
}

func (dummyNode *DummyNode) processVirt(workflowContext *WorkflowContext) {
  fmt.Println("Running ProcessVirt method.")
  var result string

  // Add the predecessor results.
  for _, predecessorNodeResult := range dummyNode.GetPredecessorResults(dummyNode, workflowContext) {
	  result += predecessorNodeResult.(string) + "|"
  }

  // Add the result at the current node.
  result += "dummy-node-result"

  dummyNode.SaveResult(dummyNode, workflowContext, result)
}