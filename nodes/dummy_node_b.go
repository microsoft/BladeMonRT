package nodes

/** Node that has a hard-coded value for its result. */
type DummyNodeB struct {
	Node
}

func (dummyNode *DummyNodeB) InitializeFields() {
	dummyNode.SetName("DummyNodeB")
}

func (dummyNode *DummyNodeB) ProcessVirt(predecessorNodeResults []interface{}) {
  dummyNode.SaveResult("dummy-result-B")
}