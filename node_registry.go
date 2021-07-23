package main

import (
	"github.com/microsoft/BladeMonRT/nodes"
	// "reflect"
	//"fmt"
)


/** Class for parsing workflow definitions. */
/*
type NodeRegistry struct {
	nameToType map[string]reflect.Type
}

func newNodeRegistrySingleton() NodeRegistry {
	var nameToType map[string]reflect.Type = make(map[string]reflect.Type)
	var nodeRegistry NodeRegistry = NodeRegistry{nameToType : nameToType}
	nodeRegistry.registerType(&nodes.DummyNode{})
	// Add all new types of nodes here.

	return nodeRegistry
}

func (nodeRegistry *NodeRegistry) registerType(typeInstance interface{}) {
	var newTypeInstance reflect.Type = reflect.TypeOf(typeInstance)
	var typeClassName string = reflect.Indirect(reflect.ValueOf(typeInstance)).Type().Name()
	nodeRegistry.nameToType[typeClassName] = newTypeInstance
}
*/

// func (nodeRegistry *NodeRegistry) makeInstance(typeName string) nodes.InterfaceNode {
	
	// var nodeType reflect.Type = nodeRegistry.nameToType[typeName]
	// var node nodes.InterfaceNode = reflect.New(nodeType).Elem().Interface().(nodes.InterfaceNode)
	// fmt.Println(node)
	// node.InitializeFields()

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

