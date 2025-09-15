package graph

import (
	"container/heap"
)

// The Dijkstra algorithm Use-Case (aka Command) object.
// It reuses the shared heap to limit the number of allocations during runtime,
// but the consequence is that the algorithm is not thread-safe. You need a
// separate instance of the algorithm for each thread, but the graph itself can
// be shared safely and can be used by multiple algorithms at the same time.
type Dijkstra[I Id, C Cost, V any, E any] struct {
	graph *Graph[I, C, V, E]
	heap  *dijkstraHeap[I, C, V, E]
	// The data that is attached to the vertices by the algorithms.
	// This is a speed optimization to avoid allocating memory for the heap and
	// vertex data on each call.
	// It stores all the Dijkstra algorithm state and can access it with O(1)
	// time complexity during runtime.
	// To find the index of the associated data for a vertex, use the vertex's
	// GetCustomDataIndex() method.
	vertexData []dijkstraVertexData[I, C]
	maxCost    C
	Amplifier  CostFunc[I, C, V, E]
}

// Creates a new Dijkstra instance for the given graph.
// This function is thread-safe and can be called concurrently as long as the
// graph doesn't change.
func NewDijkstra[I Id, C Cost, V any, E any](graph *Graph[I, C, V, E]) *Dijkstra[I, C, V, E] {
	vertexData := make([]dijkstraVertexData[I, C], len(graph.vertices))
	algorithm := &Dijkstra[I, C, V, E]{
		graph: graph,
		heap: &dijkstraHeap[I, C, V, E]{
			algorithm: nil,
		},
		vertexData: vertexData,
	}
	assignMaxNumber(&algorithm.maxCost)
	algorithm.heap.algorithm = algorithm
	return algorithm
}

// Finds the shortest path between two vertices in the graph.
// Returns a slice of vertex IDs representing the shortest path.
// Returns nil if no path is found.
// Time complexity: O(E log V) where E is the number of edges and V is the number of vertices.
// Space complexity: O(V) where V is the number of vertices.
// WARNING: This function is not thread-safe and should not be called concurrently.
func (d *Dijkstra[I, C, V, E]) FindShortestPath(start I, end I) []I {
	// Check if start and end vertices exist
	startVertex, err := d.graph.GetVertexById(start)
	if err != nil {
		return nil // Start vertex not found
	}

	endVertex, err := d.graph.GetVertexById(end)
	if err != nil {
		return nil // End vertex not found
	}

	// If start and end are the same, return the start vertex
	if start == end {
		return []I{start}
	}

	// Initialize vertex data for all vertices
	for i := range d.vertexData {
		d.vertexData[i].visited = false
		d.vertexData[i].previous = nil
		d.vertexData[i].cost = d.maxCost
	}

	// Initialize priority queue
	heap.Init(d.heap)

	// Set start vertex distance to 0 and add to queue
	startIdx := startVertex.GetCustomDataIndex()
	d.vertexData[startIdx].cost = 0
	heap.Push(d.heap, startVertex)

	// Main Dijkstra loop
	for d.heap.Len() > 0 {
		// Get vertex with minimum distance
		current := heap.Pop(d.heap).(*Vertex[I, C])
		currentIdx := current.GetCustomDataIndex()
		currentData := &d.vertexData[currentIdx]

		// Skip if already visited
		if currentData.visited {
			continue
		}

		// Mark as visited
		currentData.visited = true

		// If we reached the target, we can stop
		if current.id == end {
			break
		}

		// Process all neighbors
		for _, edge := range current.edges {
			neighbor := edge.targetVertex
			neighborIdx := neighbor.GetCustomDataIndex()
			neighborData := &d.vertexData[neighborIdx]

			// Skip if neighbor already visited
			if neighborData.visited {
				continue
			}

			edgeCost := edge.cost

			if d.Amplifier != nil {
				cost, enabled := d.Amplifier(current, &edge)
				if !enabled {
					continue
				}
				edgeCost = cost
			}

			// Calculate tentative distance
			tentativeDistance := currentData.cost + edgeCost

			// If this is a better path to the neighbor
			if tentativeDistance < neighborData.cost {
				neighborData.cost = tentativeDistance
				neighborData.previous = current
				heap.Push(d.heap, neighbor)
			}
		}
	}

	// Reconstruct path by following previous pointers
	endIdx := endVertex.GetCustomDataIndex()
	if !d.vertexData[endIdx].visited {
		return nil // No path found
	}

	path := []I{}
	current := endVertex
	for current != nil {
		path = append(path, current.id)
		currentIdx := current.GetCustomDataIndex()
		current = d.vertexData[currentIdx].previous
	}

	// Reverse the path to get start-to-end order
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}
