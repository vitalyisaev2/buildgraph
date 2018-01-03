package graph

import (
	//"fmt"
	"io/ioutil"
	"testing"

	yaml "gopkg.in/yaml.v2"

	"github.com/stretchr/testify/assert"
)

// newGraphFromYAMLFile builds a graph for a given YAML formatted file
// (used only for testing purposes)
func newGraphFromYAMLFile(path string) (Graph, error) {

	var err error

	// Parse yaml file
	rawData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	parsedData := make(map[string][]string)
	err = yaml.Unmarshal(rawData, parsedData)
	if err != nil {
		return nil, err
	}

	// Fill up new graph
	return NewGraphFromAdjacencyMap(parsedData)
}

// utility for fast checks
func NodeSeqNames(ns []Node) []string {
	ss := make([]string, len(ns))
	for _, n := range ns {
		ss = append(ss, n.Name())
	}
	return ss
}

func TestGraphInternalAPI(t *testing.T) {

	var err error
	var g Graph
	var n Node
	var successors []Node

	// simple1.yaml
	g, err = newGraphFromYAMLFile("test/simple1.yml")
	assert.NotNil(t, g)
	assert.NoError(t, err)

	// Get existing node
	n, err = g.GetNode("A")
	assert.NotNil(t, n)
	assert.NoError(t, err)

	// pay attention to closed channel when iterating over successors
	successors = n.Successors()
	assert.Len(t, successors, 1)
	assert.NotNil(t, successors[0])
	assert.Equal(t, "D", successors[0].Name())

	// Check that nonexistant node doesn't exist
	n, err = g.GetNode("Z")
	assert.Nil(t, n)
	assert.Error(t, err)
}

func TestPhasicTopologicalSort(t *testing.T) {

	var err error
	var g Graph
	var nodes []Node
	var pts PhasicTopologicalSort
	var siblingNodes [][]Node
	var nodeNames []string

	// simple1.yml - simple chain with one bifurcation
	g, _ = newGraphFromYAMLFile("test/simple1.yml")
	pts, err = g.PhasicTopologicalSortFromNode("A")
	assert.NotNil(t, pts)
	assert.NoError(t, err)

	siblingNodes = pts.SiblingNodes()
	assert.NotNil(t, siblingNodes)
	assert.Len(t, siblingNodes, 3)

	nodes = siblingNodes[0]
	assert.Len(t, nodes, 1)
	assert.Equal(t, "A", nodes[0].Name())

	nodes = siblingNodes[1]
	assert.Len(t, nodes, 1)
	assert.Equal(t, "D", nodes[0].Name())

	nodes = siblingNodes[2]
	assert.Len(t, nodes, 2)
	nodeNames = NodeSeqNames(nodes)
	assert.Contains(t, nodeNames, "F")
	assert.Contains(t, nodeNames, "G")

	// simple1.yml - root node has no successors (edge case)
	pts, err = g.PhasicTopologicalSortFromNode("G")
	assert.NotNil(t, pts)
	assert.NoError(t, err)

	siblingNodes = pts.SiblingNodes()
	assert.Len(t, siblingNodes, 1)
	nodes = siblingNodes[0]
	assert.Len(t, nodes, 1)
	assert.Equal(t, "G", nodes[0].Name())

	// simple2.yml - diverged at A, merged at E
	g, _ = newGraphFromYAMLFile("test/simple2.yml")
	pts, err = g.PhasicTopologicalSortFromNode("A")
	assert.NotNil(t, pts)
	assert.NoError(t, err)

	siblingNodes = pts.SiblingNodes()
	assert.Len(t, siblingNodes, 5)

	nodes = siblingNodes[0]
	assert.Len(t, nodes, 1)
	assert.Equal(t, "A", nodes[0].Name())

	nodes = siblingNodes[1]
	assert.Len(t, nodes, 1)
	assert.Equal(t, "B", nodes[0].Name())

	nodes = siblingNodes[2]
	assert.Len(t, nodes, 1)
	assert.Equal(t, "C", nodes[0].Name())

	nodes = siblingNodes[3]
	assert.Len(t, nodes, 1)
	assert.Equal(t, "D", nodes[0].Name())

	nodes = siblingNodes[4]
	assert.Len(t, nodes, 1)
	assert.Equal(t, "E", nodes[0].Name())
}

func TestCyclicGraph(t *testing.T) {
	var err error
	var g Graph
	var cyclic bool
	var cycleNodes []Node
	var cycleNodeNames []string

	// Acyclic
	g, err = newGraphFromYAMLFile("test/simple1.yml")
	assert.NotNil(t, g)
	assert.NoError(t, err)

	cyclic, cycleNodes, err = g.Cyclic()
	assert.False(t, cyclic)
	assert.Nil(t, cycleNodes)
	assert.Nil(t, err)

	// Cyclic 1
	g, err = newGraphFromYAMLFile("test/cyclic1.yml")
	assert.NotNil(t, g)
	assert.NoError(t, err)

	cyclic, cycleNodes, err = g.Cyclic()
	cycleNodeNames = NodeSeqNames(cycleNodes)
	assert.True(t, cyclic)
	assert.Len(t, cycleNodes, 3)
	assert.Contains(t, cycleNodeNames, "A")
	assert.Contains(t, cycleNodeNames, "B")
	assert.Contains(t, cycleNodeNames, "C")
	assert.Nil(t, err)
}
