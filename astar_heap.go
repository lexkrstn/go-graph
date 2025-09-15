package graph

// The data that is attached to the vertices by the A* algorithms.
type astarVertexData[I Id, C Cost] struct {
	previous *Vertex[I, C]
	visited  bool
	gScore   C // Cost from start to this vertex
	fScore   C // gScore + heuristic estimate to goal
}

// astarHeap implements heap.Interface for the priority queue
type astarHeap[I Id, C Cost, V any, E any] struct {
	pq        []*Vertex[I, C]
	algorithm *AStar[I, C, V, E]
}

func (h *astarHeap[I, C, V, E]) Len() int {
	return len(h.pq)
}

func (h *astarHeap[I, C, V, E]) Less(i, j int) bool {
	vertexI := h.pq[i].GetCustomDataIndex()
	vertexJ := h.pq[j].GetCustomDataIndex()

	dataI := h.algorithm.vertexData[vertexI]
	dataJ := h.algorithm.vertexData[vertexJ]

	return dataI.fScore < dataJ.fScore
}

func (h *astarHeap[I, C, V, E]) Swap(i, j int) {
	h.pq[i], h.pq[j] = h.pq[j], h.pq[i]
}

func (h *astarHeap[I, C, V, E]) Push(x any) {
	h.pq = append(h.pq, x.(*Vertex[I, C]))
}

func (h *astarHeap[I, C, V, E]) Pop() any {
	n := len(h.pq)
	node := h.pq[n-1]
	h.pq[n-1] = nil // avoid memory leak
	h.pq = h.pq[0 : n-1]
	return node
}
