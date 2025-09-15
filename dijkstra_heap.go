package graph

// The data that is attached to the vertices by the Dijkstra algorithms.
type dijkstraVertexData[I Id, C Cost] struct {
	previous *Vertex[I, C]
	visited  bool
	cost     C
}

// dijkstraHeap implements heap.Interface for the priority queue
type dijkstraHeap[I Id, C Cost, V any, E any] struct {
	pq        []*Vertex[I, C]
	graph     *Graph[I, C, V, E]
	algorithm *Dijkstra[I, C, V, E]
}

func (h *dijkstraHeap[I, C, V, E]) Len() int { return len(h.pq) }

func (h *dijkstraHeap[I, C, V, E]) Less(i, j int) bool {
	vertexI := h.pq[i].GetCustomDataIndex()
	vertexJ := h.pq[j].GetCustomDataIndex()

	dataI := h.algorithm.vertexData[vertexI]
	dataJ := h.algorithm.vertexData[vertexJ]

	return dataI.cost < dataJ.cost
}

func (h *dijkstraHeap[I, C, V, E]) Swap(i, j int) {
	h.pq[i], h.pq[j] = h.pq[j], h.pq[i]
}

func (h *dijkstraHeap[I, C, V, E]) Push(x any) {
	h.pq = append(h.pq, x.(*Vertex[I, C]))
}

func (h *dijkstraHeap[I, C, V, E]) Pop() any {
	n := len(h.pq)
	node := h.pq[n-1]
	h.pq[n-1] = nil // avoid memory leak
	h.pq = h.pq[0 : n-1]
	return node
}
