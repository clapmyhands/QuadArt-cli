package util

type HeapItem interface {
	Less(other HeapItem) bool
	More(other HeapItem) bool
}

type MaxHeap []HeapItem

func (h MaxHeap) Len() int { return len(h) }

func (h MaxHeap) Less(i, j int) bool {
	return h[i].More(h[j])
}

func (h MaxHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MaxHeap) Push(x any) {
	item := x.(HeapItem)
	*h = append(*h, item)
}

func (h *MaxHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*h = old[0 : n-1]
	return item
}
