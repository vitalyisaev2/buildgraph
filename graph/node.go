package graph

import (
	"bytes"
	"fmt"
	"sort"
)

// Node represents an item of a directed acyclic graph.
// Every node is a container that can store arbitrary data types.
// WARNING: no goroutine-safety
type Node interface {
	// Inherited interfaces
	fmt.Stringer

	// Public API
	Name() string
	Color() nodeColor
	Value() interface{}
	Successors() []Node

	// Node construction API (not used outside the package)
	link(successor Node, keepSorted bool) error

	// Node coloring API is involved in various graph algorithms (not used outside the package)
	getColor() nodeColor
	setColorWhite() error
	setColorGray() error
	setColorBlack() error
}

// nodeColor abstraction is used in graph algorithms
type nodeColor uint8

const (
	white nodeColor = 1 << iota
	gray            = 1 << iota
	black           = 1 << iota
)

// nodeImpl implements Node interface
type nodeImpl struct {
	name       string
	value      interface{}
	successors []Node
	color      nodeColor
}

// Name getter
func (n *nodeImpl) Name() string {
	return n.name
}

// Color getter
func (n *nodeImpl) Color() nodeColor {
	return n.color
}

// Value getter
func (n *nodeImpl) Value() interface{} {
	return n.value
}

// String returns string representation of node
func (n *nodeImpl) String() string {
	var err error
	var buffer bytes.Buffer

	_, err = buffer.WriteString(fmt.Sprintf("%s ->", n.name))
	if err != nil {
		return err.Error()
	}

	for _, successor := range n.successors {
		_, err = buffer.WriteString(fmt.Sprintf(" %s ", successor.name))
		if err != nil {
			return err.Error()
		}
	}

	_, err = buffer.WriteString("\n")
	if err != nil {
		return err.Error()
	}

	return buffer.String()
}

// Successors returns iterator over node successors
func (n *nodeImpl) Successors() []Node {
	items := make([]Node, 0, len(n.successors))
	for _, successor := range n.successors {
		items = append(items, successor)
	}
	return items
}

func (n *nodeImpl) link(successor Node, keepSorted bool) error {
	if successor == nil {
		return fmt.Errorf("Trying to add nil successor to node")
	} else {
		n.successors = append(n.successors, successor)
		return nil
	}
	if keepSorted {
		// I prefer to keep the sequence of successors lexicographically sorted;
		sort.Slice(
			n.successors,
			func(i, j int) bool { return n.successors[i].Name() < n.successors[j].Name() },
		)
	}
}

func (n *nodeImpl) getColor() nodeColor {
	return n.color
}

func (n *nodeImpl) setColorWhite() error {
	n.color = white
	return nil
}

func (n *nodeImpl) setColorGray() error {
	var err error
	switch n.color {
	case white:
		n.color = gray
	case gray:
		err = fmt.Errorf("algorithm error: cannot change color from gray to gray")
	case black:
		n.color = gray
	default:
		err = fmt.Errorf("algorithm error: invalid node color value: %d", n.color)
	}
	return err
}

func (n *nodeImpl) setColorBlack() error {
	var err error
	switch n.color {
	case white:
		err = fmt.Errorf("algorithm error: cannot change color from white to black")
	case gray:
		n.color = black
	case black:
		err = fmt.Errorf("algorithm error: cannot change color from black to black")
	default:
		err = fmt.Errorf("algorithm error: invalid node color value: %d", n.color)
	}
	return err
}

// NewNode returns new Node interface instance
func NewNode(nodeName string, nodeValue interface{}) Node {
	return &nodeImpl{nodeName, nodeValue, make([]Node, 0), white}
}

// nodeStack - custom stack implementation. Copied from here: http://gitlab.srv.pv.km/id/Settings/blob/master/json/converter/stack.go (much thanks to Denis Shilkin)
type nodeStack []Node

func (s *nodeStack) Push(n Node) {
	*s = append(*s, n)
}

func (s *nodeStack) Pop() Node {
	ret := (*s)[len(*s)-1]
	*s = (*s)[0 : len(*s)-1]
	return ret
}

func (s *nodeStack) PopFront() Node {
	ret := (*s)[0]
	*s = (*s)[1:]
	return ret
}

func (s *nodeStack) Top() Node {
	return (*s)[len(*s)-1]
}

func (s *nodeStack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *nodeStack) Len() int {
	return len(*s)
}
