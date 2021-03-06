package dummy_node_b

import (
	gomock "github.com/golang/mock/gomock"
	"github.com/microsoft/BladeMonRT/nodes"
	"gotest.tools/assert"
	"testing"
)

func TestMockFunctionsInDummyNodeBExample(t *testing.T) {
	// This test gives an example where we mock DummyNodeB. We need to specify the return value of the result method which is in
	// the DummyNodeB class but not in the InterfaceNode class.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNodeB := NewMockInterfaceDummyNodeB(ctrl)

	mockNodeB.EXPECT().GetResult(gomock.Any(), gomock.Any()).Return("result-in-test").AnyTimes()
	mockNodeB.EXPECT().result().Return("result-in-test")

	// Check the return value when result is called.
	returnValResult := mockNodeB.result()
	assert.Equal(t, returnValResult, "result-in-test")

	// Check the return value when GetResult is called.
	returnValGetResult := mockNodeB.GetResult(mockNodeB, nodes.NewWorkflowContext())
	assert.Equal(t, returnValGetResult, "result-in-test")
}
