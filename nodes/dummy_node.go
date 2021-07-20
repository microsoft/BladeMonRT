package nodes

import (
	"fmt"
)

type DummyNode struct {
	Node
}

func (dummyNode *DummyNode) ProcessVirt(workflowContextResult map[string]string) {
  fmt.Println("Running ProcessVirt method.")
  dummyNode.NodeResult = "dummy-node-result"
  dummyNode.SaveResult(workflowContextResult)
}
