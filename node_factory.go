package root

import (
	"github.com/microsoft/BladeMonRT/nodes"
)

/** Interface used to convert node names to node instances. */
type InterfaceNodeFactory interface {
	ConstructNode(typeName string) nodes.InterfaceNode
}

/** Concrete utility class used to convert node names to node instances. */
type NodeFactory struct {}

func (nodeFactory NodeFactory) ConstructNode(typeName string) nodes.InterfaceNode {
	switch typeName {
		case "DummyNode":
			return nodes.NewDummyNode()
		case "DummyNodeB":
			return nodes.NewDummyNodeB()
		case "DummyNodeC":
			return nodes.NewDummyNodeC()
		default:
			panic("Node for given name not found.")
	} 
	return nil
}

