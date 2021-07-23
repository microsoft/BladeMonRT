package main

import (
	"github.com/microsoft/BladeMonRT/nodes"
)

/** Utility class used to convert node names to node instances. */
type NodeFactory struct {}

func (nodeFactory *NodeFactory) constructNode(typeName string) nodes.InterfaceNode {
	switch typeName {
		case "DummyNode":
			return &nodes.DummyNode{}
		case "DummyNodeB":
			return &nodes.DummyNodeB{}
		case "DummyNodeC":
			return &nodes.DummyNodeC{}
		default:
			panic("Node for given name not found.")

	} 
	return nil
}

