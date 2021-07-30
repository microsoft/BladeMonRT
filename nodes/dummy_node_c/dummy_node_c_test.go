package dummy_node_c

import (
	"github.com/microsoft/BladeMonRT/nodes"
	"testing"
	"gotest.tools/assert"
)

/** Class that provides the result for the dummy node for testing. */
type TestResultProvider struct {
}

func (dummyNode TestResultProvider) result() string {
	return "test-provider-dummy-result-c"
}

func TestPatchFunctionsInDummyNodeCExample(t *testing.T) {
	// This test provides an example where we patch the result function using the ResultProvider interface but the production implementation of Process
	// and GetResult are used.
	var dummyNode *DummyNodeC = NewDummyNodeC(TestResultProvider{})
	var workflowContext *nodes.WorkflowContext = nodes.NewWorkflowContext()
	dummyNode.Process(dummyNode, workflowContext)
	result := dummyNode.GetResult(dummyNode, workflowContext)

	// Check that the result from the GetResult method includes the value in the TestResultProvider's result method.
	assert.Equal(t, result, "test-provider-dummy-result-c");
}
