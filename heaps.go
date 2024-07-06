package heaps

import "cmp"

type Heap[I comparable, K cmp.Ordered] interface {
	Size() int
	Insert(Node[I, K]) error
	Meld(Heap[I, K]) error
	FindMin() (Node[I, K], error)
	PopMin() (Node[I, K], error)
	DecreaseKey(Node[I, K], K) error
}

type Node[I comparable, K cmp.Ordered] interface {
	GetKey() K
	GetID() I
}
