package graph

import (
	"container/heap"
)

// HeuristicFunc represents a function that estimates the cost from a vertex to the goal.
// It should return an admissible heuristic (never overestimate the actual cost).
// The function takes the current vertex ID and the goal vertex ID and returns the estimated cost.
type HeuristicFunc[I Id, C Cost] func(current I, goal I) C

// The A* algorithm Use-Case (aka Command) object.
// It reuses the shared heap to limit the number of allocations during runtime,
// but the consequence is that the algorithm is not thread-safe. You need a
// separate instance of the algorithm for each thread, but the graph itself can
// be shared safely and can be used by multiple algorithms at the same time.
type AStar[I Id, C Cost, V any, E any] struct {
	graph     *Graph[I, C, V, E]
	heap      *astarHeap[I, C, V, E]
	heuristic HeuristicFunc[I, C]
	// The data that is attached to the vertices by the algorithms.
	// This is a speed optimization to avoid allocating memory for the heap and
	// vertex data on each call.
	// It stores all the A* algorithm state and can access it with O(1)
	// time complexity during runtime.
	// To find the index of the associated data for a vertex, use the vertex's
	// GetCustomDataIndex() method.
	vertexData []astarVertexData[I, C]
	maxCost    C
}

// Creates a new A* instance for the given graph with a heuristic function.
// This function is thread-safe and can be called concurrently as long as the
// graph doesn't change.
func NewAStar[I Id, C Cost, V any, E any](graph *Graph[I, C, V, E], heuristic HeuristicFunc[I, C]) *AStar[I, C, V, E] {
	vertexData := make([]astarVertexData[I, C], len(graph.vertices))
	algorithm := &AStar[I, C, V, E]{
		graph:      graph,
		heap:       &astarHeap[I, C, V, E]{},
		heuristic:  heuristic,
		vertexData: vertexData,
	}
	assignMaxNumber(&algorithm.maxCost)
	algorithm.heap.algorithm = algorithm
	return algorithm
}

// Finds the shortest path between two vertices in the graph using A* algorithm.
// Returns a slice of vertex IDs representing the shortest path.
// Returns nil if no path is found.
// Time complexity: O(E log V) where E is the number of edges and V is the number of vertices.
// Space complexity: O(V) where V is the number of vertices.
// WARNING: This function is not thread-safe and should not be called concurrently.
func (a *AStar[I, C, V, E]) FindShortestPath(start I, end I) []I {
	// Check if start and end vertices exist
	startVertex, err := a.graph.GetVertexById(start)
	if err != nil {
		return nil // Start vertex not found
	}

	endVertex, err := a.graph.GetVertexById(end)
	if err != nil {
		return nil // End vertex not found
	}

	// If start and end are the same, return the start vertex
	if start == end {
		return []I{start}
	}

	// Initialize vertex data for all vertices
	for i := range a.vertexData {
		a.vertexData[i].visited = false
		a.vertexData[i].previous = nil
		a.vertexData[i].gScore = a.maxCost
		a.vertexData[i].fScore = a.maxCost
	}

	// Initialize priority queue
	heap.Init(a.heap)

	// Set start vertex g-score to 0 and calculate f-score
	startIdx := startVertex.GetCustomDataIndex()
	a.vertexData[startIdx].gScore = 0
	a.vertexData[startIdx].fScore = a.heuristic(start, end)
	heap.Push(a.heap, startVertex)

	// Main A* loop
	for a.heap.Len() > 0 {
		// Get vertex with minimum f-score
		current := heap.Pop(a.heap).(*Vertex[I, C])
		currentIdx := current.GetCustomDataIndex()
		currentData := &a.vertexData[currentIdx]

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
			neighborData := &a.vertexData[neighborIdx]

			// Skip if neighbor already visited
			if neighborData.visited {
				continue
			}

			// Calculate tentative g-score (cost from start to neighbor)
			tentativeGScore := currentData.gScore + edge.cost

			// If this is a better path to the neighbor
			if tentativeGScore < neighborData.gScore {
				neighborData.gScore = tentativeGScore
				neighborData.fScore = tentativeGScore + a.heuristic(neighbor.id, end)
				neighborData.previous = current
				heap.Push(a.heap, neighbor)
			}
		}
	}

	// Reconstruct path by following previous pointers
	endIdx := endVertex.GetCustomDataIndex()
	if !a.vertexData[endIdx].visited {
		return nil // No path found
	}

	path := []I{}
	current := endVertex
	for current != nil {
		path = append(path, current.id)
		currentIdx := current.GetCustomDataIndex()
		current = a.vertexData[currentIdx].previous
	}

	// Reverse the path to get start-to-end order
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}
