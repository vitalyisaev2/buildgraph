package graph

import (
	"bytes"
	"fmt"
	"sort"
)

// nodeColor abstraction is used in graph algorithms
type nodeColor uint8

const (
	white nodeColor = 1 << iota
	gray
	black
)

// defaultNode implements Node interface
type defaultNode struct {
	name       string
	value      interface{}
	successors []Node
	color      nodeColor
}

// Name getter
func (n *defaultNode) Name() string {
	return n.name
}

// Color getter
func (n *defaultNode) Color() nodeColor {
	return n.color
}

// Value getter
func (n *defaultNode) Value() interface{} {
	return n.value
}

// String returns string representation of node
func (n *defaultNode) String() string {
	var err error
	var buffer bytes.Buffer

	_, err = buffer.WriteString(fmt.Sprintf("%s ->", n.name))
	if err != nil {
		return err.Error()
	}

	for _, successor := range n.successors {
		_, err = buffer.WriteString(fmt.Sprintf(" %s ", successor.Name()))
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
func (n *defaultNode) Successors() []Node {
	items := make([]Node, 0, len(n.successors))
	for _, successor := range n.successors {
		items = append(items, successor)
	}
	return items
}

func (n *defaultNode) link(successor Node, keepSorted bool) error {
	if successor == nil {
		return fmt.Errorf("Trying to add nil successor to node")
	}

	n.successors = append(n.successors, successor)
	if keepSorted {
		// I prefer to keep the sequence of successors lexicographically sorted;
		sort.Slice(
			n.successors,
			func(i, j int) bool { return n.successors[i].Name() < n.successors[j].Name() },
		)
	}
	return nil
}

func (n *defaultNode) getColor() nodeColor {
	return n.color
}

func (n *defaultNode) setColorWhite() error {
	n.color = white
	return nil
}

func (n *defaultNode) setColorGray() error {
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

func (n *defaultNode) setColorBlack() error {
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
	return &defaultNode{nodeName, nodeValue, make([]Node, 0), white}
}

// NodeList - custom stack implementation. Copied from here: http://gitlab.srv.pv.km/id/Settings/blob/master/json/converter/stack.go (much thanks to Denis Shilkin)
type NodeList []Node

func (s *NodeList) Push(n Node) {
	*s = append(*s, n)
}

func (s *NodeList) Pop() Node {
	ret := (*s)[len(*s)-1]
	*s = (*s)[0 : len(*s)-1]
	return ret
}

func (s *NodeList) PopFront() Node {
	ret := (*s)[0]
	*s = (*s)[1:]
	return ret
}

func (s *NodeList) Top() Node {
	return (*s)[len(*s)-1]
}

func (s *NodeList) IsEmpty() bool {
	return len(*s) == 0
}

func (s *NodeList) Len() int {
	return len(*s)
}

func (s *NodeList) String() string {
	var b bytes.Buffer
	b.WriteString("\n[\n")
	for _, node := range *s {
		b.WriteString(node.String())
	}
	b.WriteString("]\n")
	return b.String()
}
