package graph

import "fmt"

// Looks for the cycles in DAG that involve particular Node;
// in case if cycle was found, returns true and stack of nodes containing the cycle
func cycleDiscoveryFromNode(n Node) (bool, NodeList, error) {

	var err error

	// Used for discovered cycle printing
	var stack NodeList
	var cycleDiscovered = false

	// Push stack. Paint node in gray when entering it
	stack.Push(n)

	err = n.setColorGray()
	if err != nil {
		return false, nil, err
	}

	// Traverse graph and check if the cycle was discovered in loop
	for _, successor := range n.Successors() {
		err = traverseGraphCS(successor, &stack, &cycleDiscovered)
		if err != nil {
			return false, nil, err
		}
		if cycleDiscovered {
			//trim stack to get cycle only
			cycleStartNode := stack.Pop()
			for _, sn := range stack {
				if sn == cycleStartNode {
					break
				}
				_ = stack.PopFront()
			}
			return true, stack, nil
		}
	}

	// Pop stack. Paint node in black when leaving it
	_ = stack.Pop()
	err = n.setColorBlack()
	if err != nil {
		return false, nil, err
	}

	if !stack.IsEmpty() {
		return false, nil, fmt.Errorf("Algorithm error: stack is not empty after DFS")
	}

	return false, nil, nil
}

// Depth first graph traversing for the sake of cycle discovery
// Side effect: changes cycleDiscovered bool passed from cycleDiscoveryFromNode
func traverseGraphCS(n Node, stack *NodeList, cycleDiscovered *bool) error {

	var err error

	stack.Push(n)

	// Coloring gray node in gray again == cycle was found
	if n.getColor() == gray {
		*cycleDiscovered = true
		return nil
	}

	// Push stack. Paint node in gray when entering it
	err = n.setColorGray()
	if err != nil {
		return err
	}

	// Traverse graph recursivly
	for _, successor := range n.Successors() {
		err = traverseGraphCS(successor, stack, cycleDiscovered)
		if err != nil {
			return err
		}
		if *cycleDiscovered {
			return nil
		}
	}

	// Pop stack. Paint node in black when leaving it
	_ = stack.Pop()

	err = n.setColorBlack()
	return err
}
