package inttree

import "github.com/phf/go-queue/queue"

// Tree implements tree for integer values
type Tree interface {
	// Find finds data by value
	Find(val int) (foundData int, err error)
	// Append appends value to tree
	Append(val int) (err error)
	// Delete deletes value from tree
	Delete(val int) (err error)
}

// tree represents sorted v2.Card
type tree struct {
	count int
	root  *node
}

func (t *tree) Find(val int) (foundData int, err error) {
	cur := t.root
	for {
		if cur == nil {
			err = notFoundErr
			return
		}

		switch {
		case cur.val == val:
			foundData = val
			return
		case cur.val > val:
			cur = cur.lNode
		case cur.val < val:
			cur = cur.rNode
		}
	}
}

func (t *tree) Append(val int) (err error) {
	t.append(val)
	return
}

func (t *tree) Delete(val int) (err error) {
	var parent *node
	cur := t.root
	var toFixCount []*node
	isRightNode := false
	for {
		if cur == nil {
			err = notFoundErr
			return
		}

		switch {
		case cur.val == val:
			var replacement *node
			var cutLeafs *node
			switch {
			case cur.lNode == nil:
				replacement = cur.rNode
				cutLeafs = cur.lNode
			case cur.rNode == nil:
				replacement = cur.lNode
				cutLeafs = cur.rNode
			case cur.lNode.count > cur.rNode.count:
				replacement = cur.lNode
				cutLeafs = cur.rNode
			default:
				replacement = cur.rNode
				cutLeafs = cur.lNode
			}

			t.count--
			for _, n := range toFixCount {
				n.count--
			}
			if cutLeafs != nil {
				t.count = t.count - cutLeafs.count

				q := queue.New()
				q.Init()
				q.PushBack(cutLeafs)
				for q.Len() > 0 {
					elem := q.PopFront().(*node)
					t.appendToNode(replacement, elem.val)

					if elem.lNode != nil {
						q.PushBack(elem.lNode)
					}
					if elem.rNode != nil {
						q.PushBack(elem.rNode)
					}
				}
			}

			if parent == nil {
				t.root = replacement

				return
			}
			if isRightNode {
				parent.rNode = replacement
			} else {
				parent.lNode = replacement
			}
			return
		case cur.val > val:
			parent = cur
			cur = cur.lNode
			isRightNode = false
		case cur.val < val:
			parent = cur
			cur = cur.rNode
			isRightNode = true
		}
		toFixCount = append(toFixCount, parent)
	}
}

// append adds element to tree
func (t *tree) append(value int) {
	if t.root != nil {
		t.count++
		currentNode := t.root
		for {
			currentNode.count++
			if currentNode.val > value {
				// Left branch
				if currentNode.lNode != nil {
					currentNode = currentNode.lNode
					continue
				}
				currentNode.lNode = &node{
					count: 1,
					val:   value,
				}
				break
			}
			// Right branch
			if currentNode.rNode != nil {
				currentNode = currentNode.rNode
				continue
			}
			currentNode.rNode = &node{
				count: 1,
				val:   value,
			}
			break
		}
		return
	}
	t.root = &node{
		count: 1,
		val:   value,
	}
	t.count = 1
}

// append adds element to tree
func (t *tree) appendToNode(currentNode *node, val int) {
	t.count++
	n := &node{
		count: 1,
		val:   val,
	}
	for {
		currentNode.count++
		if currentNode.val > n.val {
			// Left branch
			if currentNode.lNode != nil {
				currentNode = currentNode.lNode
				continue
			}
			currentNode.lNode = n
			break
		}
		// Right branch
		if currentNode.rNode != nil {
			currentNode = currentNode.rNode
			continue
		}
		currentNode.rNode = n
		break
	}
	return
}

// NewTree creates new integer tree
func NewTree(
	values []int,
) Tree {
	t := tree{
		count: 0,
		root:  nil,
	}
	for _, value := range values {
		t.append(value)
	}
	return &t
}
