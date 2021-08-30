package frequency

// CommitCountHeap basic implementation so that CommitCount can be used with container/heap
type CommitCountHeap struct {
	Heap []CommitCount
}

func (h *CommitCountHeap) Len() int {
	return len(h.Heap)
}

func (h *CommitCountHeap) Less(i, j int) bool {
	return h.Heap[i].CommitCount > h.Heap[j].CommitCount
}

func (h *CommitCountHeap) Swap(i, j int) {
	h.Heap[i], h.Heap[j] = h.Heap[j], h.Heap[i]
}

func (h *CommitCountHeap) Push(x interface{}) {
	h.Heap = append(h.Heap, x.(CommitCount))
}

func (h *CommitCountHeap) Pop() interface{} {
	old := h.Heap
	n := len(old)
	x := old[n-1]
	h.Heap = old[0 : n-1]
	return x
}
