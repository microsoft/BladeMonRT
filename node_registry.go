package main

import (
	"github.com/microsoft/BladeMonRT/nodes"
)

func makeInstance(typeName string) nodes.InterfaceNode {
	switch typeName {
		case "DummyNode":
			return &nodes.DummyNode{}
		case "DummyNodeB":
			return &nodes.DummyNodeB{}
		case "DummyNodeC":
			return &nodes.DummyNodeC{}
	} 
	return nil
}

