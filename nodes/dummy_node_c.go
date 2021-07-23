
package nodes

/** Node that has a hard-coded value for its result. */
type DummyNodeC struct {
	Node
}

func (dummyNode *DummyNodeC) InitializeFields() {
	dummyNode.SetName("DummyNodeC")
}

func (dummyNode *DummyNodeC) ProcessVirt(predecessorNodeResults []interface{}) {
  dummyNode.SaveResult("dummy-result-c")
}