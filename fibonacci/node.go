package fibonacci

import "cmp"

// Node represents a node in a fibonacci heap
type Node[I comparable, K cmp.Ordered] struct {
	id     I
	key    K
	left   *Node[I, K]
	right  *Node[I, K]
	parent *Node[I, K]
	child  *Node[I, K]
	degree int
	marked bool
}

// NewNode creates a new node with the given id and key
func NewNode[I comparable, K cmp.Ordered](id I, key K) *Node[I, K] {
	return &Node[I, K]{
		id:  id,
		key: key,
	}
}

// GetKey returns the key of the node
func (n *Node[I, K]) GetKey() K {
	return n.key
}

// GetID returns the id of the node
func (n *Node[I, K]) GetID() I {
	return n.id
}
