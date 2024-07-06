package fibonacci

import (
	"cmp"
	"errors"

	"github.com/ssenerg/heaps"
)

// Heap represents a fibonacci heap
type Heap[I comparable, K cmp.Ordered] struct {
	root *Node[I, K]
	size int
}

// NewHeap creates a new fibonacci heap
func NewHeap[I comparable, K cmp.Ordered]() *Heap[I, K] {
	return &Heap[I, K]{}
}

// Size returns the number of nodes in the heap
func (h *Heap[I, K]) Size() int {
	return h.size
}

// Insert inserts a new node into the heap, it returns an error if the node is not of
// the fibonacci heap node
func (h *Heap[I, K]) Insert(node heaps.Node[I, K]) error {
	n, ok := node.(*Node[I, K])
	if !ok {
		return errors.New("invalid node type for insertion into fibonacci heap")
	}
	h.size++
	h.addToRoot(n)
	return nil
}

// Meld melds the heap with another heap, it returns an error if the other heap is not
// of the fibonacci heap type
func (h *Heap[I, K]) Meld(other heaps.Heap[I, K]) error {
	o, ok := other.(*Heap[I, K])
	if !ok {
		return errors.New("invalid heap type for melding with a fibonacci heap")
	}

	switch {
	case h.root == nil:
		h.root = o.root
		h.size = o.size
		return nil
	case o.root == nil:
		return nil
	}

	h.size += o.size

	last := o.root.left
	o.root.left = h.root.left
	h.root.left.right = o.root
	h.root.left = last
	h.root.left.right = h.root

	if cmp.Compare(o.root.key, h.root.key) < 0 {
		h.root = o.root
	}

	return nil
}

// FindMin returns the minimum node from the heap, it returns an error if the heap is
// empty
func (h *Heap[I, K]) FindMin() (heaps.Node[I, K], error) {
	if h.size == 0 {
		return nil, errors.New("can't find min from an empty heap")
	}
	return h.root, nil
}

// PopMin removes and returns the minimum node from the heap, it returns the minimum
// node and an error if the heap is empty
func (h *Heap[I, K]) PopMin() (heaps.Node[I, K], error) {

	minNode := h.root

	switch h.size {
	case 0:
		return nil, errors.New("can't pop min from an empty heap")
	case 1:
		h.root = nil
		h.size = 0
		return minNode, nil
	}

	// add children of minNode to the root and remove parent of them
	if minNode.child != nil {
		child := minNode.child
		for child.parent != nil {
			child.parent = nil
			child = child.right
		}
		last := child.left
		child.left = minNode.left
		minNode.left.right = child
		minNode.left = last
		minNode.left.right = minNode
	}

	// remove minNode from the root
	h.removeFromRoot(minNode)
	h.consolidate()

	h.size--

	return minNode, nil
}

// DecreaseKey decreases the key of the given node to the given key, it returns an error
// if the given key is greater than the current key
func (h *Heap[I, K]) DecreaseKey(node heaps.Node[I, K], key K) error {
	n, ok := node.(*Node[I, K])
	if !ok {
		return errors.New("invalid node type for decrease key in fibonacci heap")
	}

	if cmp.Compare(n.key, key) < 0 {
		return errors.New("the given key is greater than the current key")
	}
	if n.key == key {
		return nil
	}
	n.key = key
	parent := n.parent
	if parent != nil && cmp.Compare(n.key, parent.key) < 0 {
		h.cut(parent, n)
		h.cascadingCut(parent)
	}
	if cmp.Compare(n.key, h.root.key) < 0 {
		h.root = n
	}
	return nil
}

// addToRoot adds the given node to the root
func (h *Heap[I, K]) addToRoot(node *Node[I, K]) {
	if h.root == nil {
		h.root = node
		node.left = node
		node.right = node
		return
	}
	node.left = h.root.left
	node.right = h.root
	h.root.left.right = node
	h.root.left = node

	// update the root if the given node has smaller key, because the root should have
	// the minimum key in this implementation
	if cmp.Compare(node.key, h.root.key) < 0 {
		h.root = node
	}
}

// removeFromRoot removes the given node from the root
func (h *Heap[I, K]) removeFromRoot(node *Node[I, K]) {
	if h.root == node {
		h.root = node.right
	}
	node.left.right = node.right
	node.right.left = node.left
}

// link links the given node to the given parent node
func (h *Heap[I, K]) link(parent, node *Node[I, K]) {
	h.removeFromRoot(node)
	if parent.child == nil {
		parent.child = node
		node.left = node
		node.right = node
	} else {
		node.right = parent.child.right
		node.left = parent.child
		parent.child.right.left = node
		parent.child.right = node
	}
	node.parent = parent
	parent.degree++
	node.marked = false
}

// consolidate consolidates the heap
func (h *Heap[I, K]) consolidate() {
	// degreeToRoot's length is bounded by log_phi(size) where phi is the golden ratio
	degreeToRoot := make(map[int]*Node[I, K])

	current, last := h.root, h.root.left
	for {
		right := current.right
		x, degree := current, current.degree
		for {
			y, ok := degreeToRoot[degree]
			if !ok {
				break
			}
			if cmp.Compare(y.key, x.key) < 0 {
				y, x = x, y
			}
			h.link(x, y)
			delete(degreeToRoot, degree)
			degree++
		}
		degreeToRoot[degree] = x
		if current == last {
			break
		}
		current = right
	}

	h.root = nil
	for _, v := range degreeToRoot {
		h.addToRoot(v)
	}

}

// cut cuts the given node from the given parent node
func (h *Heap[I, K]) cut(parent, node *Node[I, K]) {
	// remove node from parent's children
	switch parent.child {
	case parent.child.right:
		parent.child = nil
	case node:
		parent.child = node.right
	}
	node.left.right = node.right
	node.right.left = node.left

	// add node to the root
	h.addToRoot(node)

	parent.degree--
	node.parent = nil
	node.marked = false
}

// cascadingCut performs cascading cut on the given node
func (h *Heap[I, K]) cascadingCut(node *Node[I, K]) {
	parent := node.parent
	if parent == nil {
		return
	}
	switch node.marked {
	case false:
		node.marked = true
	case true:
		h.cut(parent, node)
		h.cascadingCut(parent)
	}
}
