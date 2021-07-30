package dummy_node_c

import (
	"log"
	"testing"

	"github.com/microsoft/BladeMonRT/logging"
	"github.com/microsoft/BladeMonRT/nodes"
	"gotest.tools/assert"
)

/** Class that provides the result for the dummy node for testing. */
type TestResultProvider struct {
}

func (dummyNode TestResultProvider) result() string {
	return "test-provider-dummy-result-c"
}

func NewDummyNodeCForTest() *DummyNodeC {
	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("DummyNodeB")
	var dummyNode DummyNodeC = DummyNodeC{Node: nodes.Node{Logger: logger}, resultProvider: TestResultProvider{}}
	return &dummyNode
}

func TestPatchFunctionsInDummyNodeCExample(t *testing.T) {
	// This test provides an example where we patch the result function using the ResultProvider interface. The production implementations of other
	// DummyNodeC methods like Process and GetResult are used.
	var dummyNode *DummyNodeC = NewDummyNodeCForTest()
	var workflowContext *nodes.WorkflowContext = nodes.NewWorkflowContext()
	dummyNode.Process(dummyNode, workflowContext)
	result := dummyNode.GetResult(dummyNode, workflowContext)

	// Check that the result from the GetResult method includes the value in the TestResultProvider's result method.
	assert.Equal(t, result, "test-provider-dummy-result-c")
}
