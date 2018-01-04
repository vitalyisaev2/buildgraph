package graph

import "fmt"

// Graph represents directed acyclic graph consisting of nodes of arbitrary types;
// every graph node (or vertex) has string identifier.
type Graph interface {
	fmt.Stringer

	phasicTopologicalSortBuilder

	// public API
	GetNode(string) (Node, error)
	CreateNode(string, interface{}) (Node, error)
	Link(parent string, child string) error
	Cyclic() (bool, NodeList, error)
}

// Node represents a node (vertex) of a directed acyclic graph.
// Every node is a container that can store value of arbitrary data types.
type Node interface {
	// Inherited interfaces
	fmt.Stringer

	// Public API
	Name() string
	Color() nodeColor
	Value() interface{}
	Successors() []Node

	// Node construction API
	// (not used outside the package)
	link(successor Node, keepSorted bool) error

	// Node coloring API is involved in various graph algorithms
	// (not used outside the package)
	getColor() nodeColor
	setColorWhite() error
	setColorGray() error
	setColorBlack() error
}
