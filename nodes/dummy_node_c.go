
package nodes

/** Node that has a hard-coded value for its result. */
type DummyNodeC struct {
	Node
	result string
}

func (dummyNode *DummyNodeC) InitializeFields() {
	dummyNode.result = "dummy-result-c"
}

func (dummyNode *DummyNodeC) processVirt(predecessorNodesResults []interface{}) {
  dummyNode.SaveResult(dummyNode.result)
}