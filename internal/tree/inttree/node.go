package inttree

// node is element of tree
type node struct {
	lNode *node
	rNode *node
	count int
	val   int
}
