package nodes

import (
	"fmt"
)

/** Node that has the concatenation of its predecessors' results and a hard-coded value for its result. */
type DummyNode struct {
	Node
}

func (dummyNode *DummyNode) InitializeFields() {
	dummyNode.SetName("DummyNode")
}



func (dummyNode *DummyNode) ProcessVirt(predecessorNodeResults []interface{}) {
  fmt.Println("Running ProcessVirt method.")
  var result string

  // Add the predecessor results.
  for _, predecessorResult := range predecessorNodeResults {
	  result += predecessorResult.(string) + "|"
  }

  // Add the result at the current node.
  result += "dummy-node-result"

  dummyNode.SaveResult(result)
}