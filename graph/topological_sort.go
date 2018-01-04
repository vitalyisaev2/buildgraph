package graph

import (
	"bytes"
	"fmt"
)

// PhasicTopologicalSort performs sorting of directed acyclic graph for the given root,
// and returns sequence of sibling Nodes
type PhasicTopologicalSort interface {
	fmt.Stringer
	SiblingNodes() [][]Node
}

type phasicTopologicalSortBuilder interface {
	PhasicTopologicalSortFromNode(string) (PhasicTopologicalSort, error)
}

type phasicTopologicalSort struct {
	siblingNodes [][]Node
}

// SiblingNodes returns iterator over sibling (or same level) nodes
func (pts *phasicTopologicalSort) SiblingNodes() [][]Node { return pts.siblingNodes }

// String returns string representation of phasicTopologicalSort
func (pts *phasicTopologicalSort) String() string {
	var err error
	var buffer bytes.Buffer

	for phaseID, phaseNodes := range pts.siblingNodes {
		_, err = buffer.WriteString(fmt.Sprintf("%d -> ", phaseID+1))
		if err != nil {
			return err.Error()
		}

		for _, node := range phaseNodes {
			_, err = buffer.WriteString(fmt.Sprintf(" %s ", node.Name()))
			if err != nil {
				return err.Error()
			}
		}

		_, err = buffer.WriteString("\n")
		if err != nil {
			return err.Error()
		}
	}
	return buffer.String()
}

// Visits all the nodes than belong to subgraph built from a given root node and
// stores the maximal distance from the root for the every node
func phasicTopologicalSortFromNode(n Node) (PhasicTopologicalSort, error) {

	var stack NodeList
	nodeLevels := make(map[Node]int)

	// Traverse graph and store nodeLevels in map
	stack.Push(n)
	nodeLevels[n] = stack.Len()
	for _, successor := range n.Successors() {
		traverseGraphPTS(successor, stack, nodeLevels)
	}
	_ = stack.Pop()
	if !stack.IsEmpty() {
		return nil, fmt.Errorf("Algorithm error: stack is not empty after DFS")
	}

	// Invert nodeLevels to sequence of node slices
	siblingNodesMap := make(map[int][]Node)
	for n, level := range nodeLevels {
		siblingNodesMap[level] = append(siblingNodesMap[level], n)
	}
	var siblingNodesSeq [][]Node
	for i := 0; i < len(siblingNodesMap); i++ {
		siblingNodesSeq = append(siblingNodesSeq, siblingNodesMap[i+1])
	}

	pts := &phasicTopologicalSort{siblingNodesSeq}
	return pts, nil
}

// Depth first graph traversing for the sake of Phasic Topological Search
func traverseGraphPTS(n Node, stack NodeList, nodeLevels map[Node]int) {
	// Append particular node to stack
	stack.Push(n)

	// Need to compare and may be upsert node level
	nodeCurrentLevel, ok := nodeLevels[n]
	if !ok {
		nodeCurrentLevel = 0
	}
	if stack.Len() > nodeCurrentLevel {
		nodeLevels[n] = len(stack)
	}

	// Map traverse onto successors
	for _, successor := range n.Successors() {
		traverseGraphPTS(successor, stack, nodeLevels)
	}

	// Pop particular node from stack
	_ = stack.Pop()
}
