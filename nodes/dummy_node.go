package nodes

import (
	"fmt"
)

/** Node that has the concatenation of its predecessors' results and a hard-coded value for its result. */
type DummyNode struct {
	Node
}

func (dummyNode *DummyNode) InitializeFields() {
	// No fields to initialize.
}

func (dummyNode *DummyNode) processVirt(predecessorNodesResults []interface{}) {
  fmt.Println("Running ProcessVirt method.")
  var result string

  // Add the predecessor results.
  for _, predecessorNodeResult := range predecessorNodesResults {
	  result += predecessorNodeResult.(string) + "|"
  }

  // Add the result at the current node.
  result += "dummy-node-result"

  dummyNode.SaveResult(result)
}