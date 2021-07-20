package nodes

import (
	"fmt"
)

// Node that has a hard-coded value for its result
type DummyNode struct {
	Node
}

func (dummyNode *DummyNode) ProcessVirt(workflowContextResult map[string]string) {
  fmt.Println("Running ProcessVirt method.")
  dummyNode.NodeResult = "dummy-node-result"
  dummyNode.SaveResult(workflowContextResult)
}
