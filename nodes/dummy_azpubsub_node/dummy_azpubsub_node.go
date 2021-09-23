package dummy_azpubsub_node

import (
	"log"

	"github.com/microsoft/BladeMonRT/logging"
	"github.com/microsoft/BladeMonRT/nodes"
	"github.com/microsoft/BladeMonRT/AzPubSub"
	"github.com/microsoft/BladeMonRT/configs"
)

/** Node that has a hard-coded value for its result. */
type DummyAzPubSubNode struct {
	nodes.Node
}

func NewDummyAzPubSubNode() *DummyAzPubSubNode {
	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("DummyAzPubSubNode")
	var dummyNode DummyAzPubSubNode = DummyAzPubSubNode{Node: nodes.Node{Logger: logger}}
	return &dummyNode
}

func (dummyNode *DummyAzPubSubNode) ProcessVirt(workflowContext *nodes.WorkflowContext) error {
	dummyNode.Logger.Println("Running process virt for the DummyAzPubSubNode.")

	vip, err := azpubsub.NewUtils().FetchAzPubSubPfVIP(configs.PF_AZ_PUB_SUB_VIP_FILE)
	if (err != nil) {
		return err
	}
	client := azpubsub.NewAzPubSubSimpleClient(false, vip)

	response, err := client.SendMessage("AzureCompute.Anvil.Request","test_message_1")
	dummyNode.Logger.Println("The response from AzPubSub is", response)
	if (err != nil) {
		return err
	}
	return err
}