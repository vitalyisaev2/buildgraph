package graph

import (
	"bytes"
	"fmt"
	"sort"
)

// Graph represents directed acyclic graph of project dependencies;
// Graph is a container and it can store arbitrary values in its Nodes
type Graph interface {
	// Inherited interfaces
	fmt.Stringer

	cycleSearcher
	phasicTopologicalSortBuilder

	// Graph construction API
	GetNode(string) (Node, error)
	CreateNode(string, interface{}) (Node, error)
	GetOrCreateNode(string, interface{}) Node
	Link(parent string, child string) error

	//SortedKeys() []string
	//Items() map[string]Node
}

type graphImpl struct {
	storage map[string]Node
}

func (g *graphImpl) GetNode(nodeName string) (Node, error) {

	if n, ok := g.storage[nodeName]; ok {
		return n, nil
	}
	return nil, fmt.Errorf("Node %s does not exist", nodeName)
}

func (g *graphImpl) CreateNode(nodeName string, nodeValue interface{}) (Node, error) {

	existingNode, _ := g.GetNode(nodeName)
	if existingNode != nil {
		return nil, fmt.Errorf("Trying to rewrite existing nodeName: %s", nodeName)
	}

	newNode := NewNode(nodeName, nodeValue)
	g.storage[nodeName] = newNode
	return newNode, nil
}

func (g *graphImpl) GetOrCreateNode(nodeName string, nodeValue interface{}) Node {

	if n, ok := g.storage[nodeName]; ok {
		return n
	}

	newNode := NewNode(nodeName, nodeValue)
	g.storage[nodeName] = newNode
	return newNode
}

func (g *graphImpl) Link(nodeName string, successorName string) error {

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

	return n.link(s)
}

// PhasicTopologicalSortFromRoot returns new PhasicTopologicalSort for a given root
func (g *graphImpl) PhasicTopologicalSortFromNode(rootName string) (PhasicTopologicalSort, error) {

	// Get root node by name
	root, err := g.GetNode(rootName)
	if err != nil {
		return nil, err
	}

	return phasicTopologicalSortFromNode(root)
}

// Cyclic property performs cycle discovery in the given Directed Graph
func (g *graphImpl) Cyclic() (bool, []Node, error) {

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
func (g *graphImpl) SortedKeys() []string {
	keys := make([]string, 0, len(g.storage))
	for key := range g.storage {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (g *graphImpl) Items() map[string]Node { return g.storage }

// String returns string representation of graph
func (g *graphImpl) String() string {

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
	return &graphImpl{make(map[string]Node)}
}
