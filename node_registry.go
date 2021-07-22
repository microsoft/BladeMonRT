package main

import (
	"github.com/microsoft/BladeMonRT/nodes"
	"reflect"
)

/** Class for parsing workflow definitions. */
type NodeRegistry struct {
	nameToType map[string]reflect.Type
}

func newNodeRegistrySingleton() NodeRegistry {
	var nameToType map[string]reflect.Type = make(map[string]reflect.Type)
	var nodeRegistry NodeRegistry = NodeRegistry{nameToType : nameToType}
	nodeRegistry.registerType(&nodes.DummyNode{})
	return nodeRegistry
}

func (nodeRegistry *NodeRegistry) registerType(typeInstance nodes.InterfaceNode) {
	var newTypeInstance reflect.Type = reflect.TypeOf(typeInstance).Elem()
	nodeRegistry.nameToType[newTypeInstance.Name()] = newTypeInstance
}

func (nodeRegistry *NodeRegistry) makeInstance(typeName string) nodes.InterfaceNode {
	return reflect.New(nodeRegistry.nameToType[typeName]).Elem().Interface().(nodes.InterfaceNode)
}

