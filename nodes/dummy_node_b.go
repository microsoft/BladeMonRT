package nodes

/** Node that has a hard-coded value for its result. */
type DummyNodeB struct {
	Node
}

func (dummyNode *DummyNodeB) InitializeFields() {
	// No fields to initialize.
}

func (dummyNode *DummyNodeB) processVirt(predecessorNodesResults []interface{}) {
  dummyNode.SaveResult("dummy-result-B")
}