package inttree

// branchDirection is direction of element of stack. It's required by skip/take mechanics
type branchDirection int

const (
	bDirRoot  branchDirection = iota
	bDirLeft  branchDirection = iota
	bDirRight branchDirection = iota
)

// stackNode is element of stack. It's required by skip/take mechanics
type stackNode struct {
	node *node
	dir  branchDirection
}

// node is element of tree
type node struct {
	lNode *node
	rNode *node
	count int
	val   int
}

// skip skips certain elements of tree and mutates current stack. It's used by further take mechanics
func (n *node) skip(skipCount int, selfDir branchDirection, currentStack *[]stackNode) (skipped int) {
	if n.count <= skipCount {
		skipped = n.count
		return
	}

	*currentStack = append(*currentStack, stackNode{
		node: n,
		dir:  selfDir,
	})

	if n.lNode != nil {
		skip := n.lNode.skip(skipCount, bDirLeft, currentStack)
		skipped += skip
		skipCount -= skip
	}

	if skipCount <= 0 {
		return
	}

	// Skip current
	skipped++
	skipCount--

	if n.rNode != nil {
		skip := n.rNode.skip(skipCount, bDirRight, currentStack)
		skipped += skip
	}

	return
}

// take takes this element and all children
func (n *node) take(takeCount int, currentValues *[]int) (taken int) {
	if n.lNode != nil {
		take := n.lNode.take(takeCount, currentValues)
		taken += take
		takeCount -= take
		if takeCount <= 0 {
			return
		}
	}

	*currentValues = append(*currentValues, n.val)
	taken++
	takeCount--
	if takeCount <= 0 {
		return
	}

	if n.rNode != nil {
		take := n.rNode.take(takeCount, currentValues)
		taken += take
		takeCount -= take
		if takeCount <= 0 {
			return
		}
	}

	return
}
