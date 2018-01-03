package graph

import (
	"bytes"
	"fmt"
	"sort"
)

type defaultGraph struct {
	storage map[string]Node
}

func (g *defaultGraph) GetNode(nodeName string) (Node, error) {

	if n, ok := g.storage[nodeName]; ok {
		return n, nil
	}
	return nil, fmt.Errorf("Node %s does not exist", nodeName)
}

func (g *defaultGraph) CreateNode(nodeName string, nodeValue interface{}) (Node, error) {

	existingNode, _ := g.GetNode(nodeName)
	if existingNode != nil {
		return nil, fmt.Errorf("Node %s already exists", nodeName)
	}

	newNode := NewNode(nodeName, nodeValue)
	g.storage[nodeName] = newNode
	return newNode, nil
}

func (g *defaultGraph) Link(nodeName string, successorName string) error {

	var (
		n, s Node
		ok   bool
	)

	if n, ok = g.storage[nodeName]; !ok {
		return fmt.Errorf("Link: adding successor %s to nonexistant node %s", successorName, nodeName)
	}

	if s, ok = g.storage[successorName]; !ok {
		return fmt.Errorf("Link: adding nonexistant successor %s to node %s", successorName, nodeName)
	}

	return n.link(s, true)
}

// PhasicTopologicalSortFromRoot returns new PhasicTopologicalSort for a given root
func (g *defaultGraph) PhasicTopologicalSortFromNode(rootName string) (PhasicTopologicalSort, error) {

	// Get root node by name
	root, err := g.GetNode(rootName)
	if err != nil {
		return nil, err
	}

	return phasicTopologicalSortFromNode(root)
}

// Cyclic property performs cycle discovery in the given Directed Graph
func (g *defaultGraph) Cyclic() (bool, NodeList, error) {

	// Need to run discovery for the every node (until first cycle occurrence at least)
	for _, n := range g.storage {

		cycleDiscovered, stack, err := cycleDiscoveryFromNode(n)
		if err != nil {
			return false, nil, err
		}
		if cycleDiscovered {
			return cycleDiscovered, stack, err
		}

		// Cleanup after every run
		for _, nodeToClean := range g.storage {
			err := nodeToClean.setColorWhite()
			if err != nil {
				return false, nil, err
			}
		}
	}
	return false, nil, nil
}

//
func (g *defaultGraph) SortedKeys() []string {
	keys := make([]string, 0, len(g.storage))
	for key := range g.storage {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (g *defaultGraph) Items() map[string]Node { return g.storage }

// String returns string representation of graph
func (g *defaultGraph) String() string {

	// Sort node names for pretty printing
	var nodeNames []string
	for name := range g.storage {
		nodeNames = append(nodeNames, name)
	}
	sort.Strings(nodeNames)

	// Fill buffer with pretty printed strings
	var buffer bytes.Buffer
	_, _ = buffer.WriteString("\n")
	for _, nodeName := range nodeNames {
		_, _ = buffer.WriteString(g.storage[nodeName].String())
	}

	return buffer.String()
}

// NewGraph returns new empty Graph
func NewGraph() Graph {
	return &defaultGraph{make(map[string]Node)}
}

// NewGraphFromAdjacencyMap builds Graph for a given Adjacency Map
// (which is a sort of adjacency list)
func NewGraphFromAdjacencyMap(dependencies map[string][]string) (Graph, error) {
	g := &defaultGraph{make(map[string]Node)}
	for parent, children := range dependencies {
		if _, exists := g.storage[parent]; !exists {
			g.storage[parent] = NewNode(parent, nil)
		}
		for _, child := range children {
			if _, exists := g.storage[child]; !exists {
				g.storage[child] = NewNode(child, nil)
			}
			if err := g.Link(parent, child); err != nil {
				return nil, err
			}
		}
	}
	return g, nil
}
