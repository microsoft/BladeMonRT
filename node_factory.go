package main

import (
	"github.com/microsoft/BladeMonRT/nodes"
	"github.com/microsoft/BladeMonRT/nodes/dummy_node_a"
	"github.com/microsoft/BladeMonRT/nodes/dummy_node_b"
	"github.com/microsoft/BladeMonRT/nodes/dummy_node_c"
)

/** Interface used to convert node names to node instances. */
type InterfaceNodeFactory interface {
	constructNode(typeName string) nodes.InterfaceNode
}

/** Concrete utility class used to convert node names to node instances. */
type NodeFactory struct {}

func (nodeFactory NodeFactory) constructNode(typeName string) nodes.InterfaceNode {
	switch typeName {
		case "DummyNodeA":
			return dummy_node_a.NewDummyNodeA()
		case "DummyNodeB":
			return dummy_node_b.NewDummyNodeB()
		case "DummyNodeC":
			return dummy_node_c.NewDummyNodeC(dummy_node_c.ResultProvider{})
		default:
			panic("Node for given name not found.")
	} 
	return nil
}

