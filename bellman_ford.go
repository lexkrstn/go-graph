package graph

// The Bellman-Ford algorithm Use-Case (aka Command) object.
// It reuses the shared vertex data to limit the number of allocations during runtime,
// but the consequence is that the algorithm is not thread-safe. You need a
// separate instance of the algorithm for each thread, but the graph itself can
// be shared safely and can be used by multiple algorithms at the same time.
type BellmanFord[I Id, C Cost, V any, E any] struct {
	graph *Graph[I, C, V, E]
	// The data that is attached to the vertices by the algorithms.
	// This is a speed optimization to avoid allocating memory for the
	// vertex data on each call.
	// It stores all the Bellman-Ford algorithm state and can access it with O(1)
	// time complexity during runtime.
	// To find the index of the associated data for a vertex, use the vertex's
	// GetCustomDataIndex() method.
	vertexData []bellmanFordVertexData[I, C]
	maxCost    C
	Amplifier  CostFunc[I, C, V, E]
}

// Creates a new Bellman-Ford instance for the given graph.
// This function is thread-safe and can be called concurrently as long as the
// graph doesn't change.
func NewBellmanFord[I Id, C Cost, V any, E any](graph *Graph[I, C, V, E]) *BellmanFord[I, C, V, E] {
	vertexData := make([]bellmanFordVertexData[I, C], len(graph.vertices))
	algorithm := &BellmanFord[I, C, V, E]{
		graph:      graph,
		vertexData: vertexData,
	}
	assignMaxNumber(&algorithm.maxCost)
	return algorithm
}

// Finds the shortest path between two vertices in the graph.
// Returns a slice of vertex IDs representing the shortest path.
// Returns nil if no path is found or if a negative cycle is detected.
// Time complexity: O(VE) where E is the number of edges and V is the number of vertices.
// Space complexity: O(V) where V is the number of vertices.
// WARNING: This function is not thread-safe and should not be called concurrently.
func (bf *BellmanFord[I, C, V, E]) FindShortestPath(start I, end I) []I {
	// Check if start and end vertices exist
	startVertex, err := bf.graph.GetVertexById(start)
	if err != nil {
		return nil // Start vertex not found
	}

	endVertex, err := bf.graph.GetVertexById(end)
	if err != nil {
		return nil // End vertex not found
	}

	// If start and end are the same, return the start vertex
	if start == end {
		return []I{start}
	}

	// Initialize vertex data for all vertices
	for i := range bf.vertexData {
		bf.vertexData[i].previous = nil
		bf.vertexData[i].cost = bf.maxCost
	}

	// Set start vertex distance to 0
	startIdx := startVertex.GetCustomDataIndex()
	bf.vertexData[startIdx].cost = 0

	// Relax all edges V-1 times
	for i := 0; i < len(bf.graph.vertices)-1; i++ {
		bf.relaxAllEdges()
	}

	// Check for negative cycles by trying to relax edges one more time
	if bf.hasNegativeCycle() {
		return nil // Negative cycle detected
	}

	// Check if end vertex is reachable
	endIdx := endVertex.GetCustomDataIndex()
	if bf.vertexData[endIdx].cost == bf.maxCost {
		return nil // No path found
	}

	// Reconstruct path by following previous pointers
	path := []I{}
	current := endVertex
	for current != nil {
		path = append(path, current.id)
		currentIdx := current.GetCustomDataIndex()
		current = bf.vertexData[currentIdx].previous
	}

	// Reverse the path to get start-to-end order
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}

// Relaxes all edges in the graph once.
// This is the core operation of the Bellman-Ford algorithm.
func (bf *BellmanFord[I, C, V, E]) relaxAllEdges() {
	for i := range bf.graph.vertices {
		current := &bf.graph.vertices[i]
		currentIdx := current.GetCustomDataIndex()
		currentData := &bf.vertexData[currentIdx]

		// Skip if current vertex is not reachable
		if currentData.cost == bf.maxCost {
			continue
		}

		// Process all outgoing edges
		for _, edge := range current.edges {
			neighbor := edge.targetVertex
			neighborIdx := neighbor.GetCustomDataIndex()
			neighborData := &bf.vertexData[neighborIdx]

			edgeCost := edge.cost

			if bf.Amplifier != nil {
				cost, enabled := bf.Amplifier(current, &edge)
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
			}
		}
	}
}

// Checks if there's a negative cycle in the graph.
// Returns true if a negative cycle is detected, false otherwise.
func (bf *BellmanFord[I, C, V, E]) hasNegativeCycle() bool {
	for i := range bf.graph.vertices {
		current := &bf.graph.vertices[i]
		currentIdx := current.GetCustomDataIndex()
		currentData := &bf.vertexData[currentIdx]

		// Skip if current vertex is not reachable
		if currentData.cost == bf.maxCost {
			continue
		}

		// Process all outgoing edges
		for _, edge := range current.edges {
			neighbor := edge.targetVertex
			neighborIdx := neighbor.GetCustomDataIndex()
			neighborData := &bf.vertexData[neighborIdx]

			edgeCost := edge.cost

			if bf.Amplifier != nil {
				cost, enabled := bf.Amplifier(current, &edge)
				if !enabled {
					continue
				}
				edgeCost = cost
			}

			// Calculate tentative distance
			tentativeDistance := currentData.cost + edgeCost

			// If we can still improve the distance, there's a negative cycle
			if tentativeDistance < neighborData.cost {
				return true
			}
		}
	}
	return false
}

// Detects if there's a negative cycle reachable from the given start vertex.
// Returns true if a negative cycle is detected, false otherwise.
// Time complexity: O(VE) where E is the number of edges and V is the number of vertices.
// WARNING: This function is not thread-safe and should not be called concurrently.
func (bf *BellmanFord[I, C, V, E]) HasNegativeCycle(start I) bool {
	// Check if start vertex exists
	startVertex, err := bf.graph.GetVertexById(start)
	if err != nil {
		return false // Start vertex not found
	}

	// Initialize vertex data for all vertices
	for i := range bf.vertexData {
		bf.vertexData[i].previous = nil
		bf.vertexData[i].cost = bf.maxCost
	}

	// Set start vertex distance to 0
	startIdx := startVertex.GetCustomDataIndex()
	bf.vertexData[startIdx].cost = 0

	// Relax all edges V-1 times
	for i := 0; i < len(bf.graph.vertices)-1; i++ {
		bf.relaxAllEdges()
	}

	// Check for negative cycles by trying to relax edges one more time
	return bf.hasNegativeCycle()
}
