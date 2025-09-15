package graph

// The data that is attached to the vertices by the DFS algorithm.
type dfsVertexData[I Id, C Cost] struct {
	visited bool
	parent  *Vertex[I, C]
	// For cycle detection: unvisited, visiting (in current path), visited
	visiting bool
}

// The DFS algorithm Use-Case (aka Command) object.
// It provides methods to perform depth-first search operations on the graph.
// The algorithm is not thread-safe and should not be called concurrently.
type DFS[I Id, C Cost, V any, E any] struct {
	graph      *Graph[I, C, V, E]
	vertexData []dfsVertexData[I, C]
}

// Creates a new DFS instance for the given graph.
// This function is thread-safe and can be called concurrently as long as the
// graph doesn't change.
func NewDFS[I Id, C Cost, V any, E any](graph *Graph[I, C, V, E]) *DFS[I, C, V, E] {
	vertexData := make([]dfsVertexData[I, C], len(graph.vertices))
	algorithm := &DFS[I, C, V, E]{
		graph:      graph,
		vertexData: vertexData,
	}
	return algorithm
}

// TraverseFrom performs a depth-first search starting from the given vertex,
// calling the provided callback function for each vertex and edge visited.
// The callback receives the current vertex and the edge that led to it (nil for the start vertex).
// Time complexity: O(V + E) where V is the number of vertices and E is the number of edges.
// Space complexity: O(V) where V is the number of vertices.
// WARNING: This function is not thread-safe and should not be called concurrently.
func (d *DFS[I, C, V, E]) TraverseFrom(start I, callback func(vertex *Vertex[I, C], edge *Edge[I, C])) {
	// Check if start vertex exists
	startVertex, err := d.graph.GetVertexById(start)
	if err != nil {
		return // Start vertex not found
	}

	// Initialize vertex data for all vertices
	for i := range d.vertexData {
		d.vertexData[i].visited = false
		d.vertexData[i].parent = nil
		d.vertexData[i].visiting = false
	}

	// Perform DFS traversal with callback
	d.dfsTraverseWithCallback(startVertex, nil, callback)
}

// FindPath finds a path from start to end vertex using DFS.
// Returns a slice of vertex IDs representing the path, or nil if no path exists.
// Time complexity: O(V + E) where V is the number of vertices and E is the number of edges.
// Space complexity: O(V) where V is the number of vertices.
// WARNING: This function is not thread-safe and should not be called concurrently.
func (d *DFS[I, C, V, E]) FindPath(start I, end I) []I {
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
		d.vertexData[i].parent = nil
		d.vertexData[i].visiting = false
	}

	// Perform DFS to find path
	found := d.dfsSearch(startVertex, endVertex)
	if !found {
		return nil // No path found
	}

	// Reconstruct path by following parent pointers
	path := []I{}
	current := endVertex
	for current != nil {
		path = append(path, current.id)
		currentIdx := current.GetCustomDataIndex()
		current = d.vertexData[currentIdx].parent
	}

	// Reverse the path to get start-to-end order
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}

// IsReachable checks if there is a path from start to end vertex.
// Returns true if a path exists, false otherwise.
// Time complexity: O(V + E) where V is the number of vertices and E is the number of edges.
// Space complexity: O(V) where V is the number of vertices.
// WARNING: This function is not thread-safe and should not be called concurrently.
func (d *DFS[I, C, V, E]) IsReachable(start I, end I) bool {
	// Check if start and end vertices exist
	startVertex, err := d.graph.GetVertexById(start)
	if err != nil {
		return false // Start vertex not found
	}

	_, err = d.graph.GetVertexById(end)
	if err != nil {
		return false // End vertex not found
	}

	// If start and end are the same, they are reachable
	if start == end {
		return true
	}

	// Initialize vertex data for all vertices
	for i := range d.vertexData {
		d.vertexData[i].visited = false
		d.vertexData[i].parent = nil
		d.vertexData[i].visiting = false
	}

	// Perform DFS to check reachability
	return d.dfsSearch(startVertex, end)
}

// GetAllReachable returns all vertices reachable from the start vertex.
// Returns a slice of vertex IDs that can be reached from the start vertex.
// Time complexity: O(V + E) where V is the number of vertices and E is the number of edges.
// Space complexity: O(V) where V is the number of vertices.
// WARNING: This function is not thread-safe and should not be called concurrently.
func (d *DFS[I, C, V, E]) GetAllReachable(start I) []I {
	// Check if start vertex exists
	startVertex, err := d.graph.GetVertexById(start)
	if err != nil {
		return nil // Start vertex not found
	}

	// Initialize vertex data for all vertices
	for i := range d.vertexData {
		d.vertexData[i].visited = false
		d.vertexData[i].parent = nil
		d.vertexData[i].visiting = false
	}

	// Perform DFS to find all reachable vertices
	var result []I
	d.dfsTraverse(startVertex, &result)
	return result
}

// dfsTraverse performs the actual DFS traversal starting from the given vertex.
// It marks all reachable vertices as visited and adds them to the result slice.
// Uses an iterative approach with an explicit stack to avoid recursion.
func (d *DFS[I, C, V, E]) dfsTraverse(vertex *Vertex[I, C], result *[]I) {
	// Use a stack to store vertices to visit
	stack := []*Vertex[I, C]{vertex}

	for len(stack) > 0 {
		// Pop vertex from stack
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		currentIdx := current.GetCustomDataIndex()
		currentData := &d.vertexData[currentIdx]

		// Skip if already visited
		if currentData.visited {
			continue
		}

		// Mark as visited and add to result
		currentData.visited = true
		*result = append(*result, current.GetId())

		// Add all unvisited neighbors to stack
		// We reverse the order to maintain the same traversal order as recursive version
		edges := current.GetEdges()
		for i := len(edges) - 1; i >= 0; i-- {
			neighbor := edges[i].GetTargetVertex()
			neighborIdx := neighbor.GetCustomDataIndex()
			neighborData := &d.vertexData[neighborIdx]

			if !neighborData.visited {
				neighborData.parent = current
				stack = append(stack, neighbor)
			}
		}
	}
}

// dfsTraverseWithCallback performs DFS traversal with a callback function.
// It marks all reachable vertices as visited and calls the callback for each vertex and edge.
// Uses an iterative approach with an explicit stack to avoid recursion.
func (d *DFS[I, C, V, E]) dfsTraverseWithCallback(startVertex *Vertex[I, C], startEdge *Edge[I, C], callback func(vertex *Vertex[I, C], edge *Edge[I, C])) {
	// Use a stack to store vertices and their incoming edges
	type stackItem struct {
		vertex *Vertex[I, C]
		edge   *Edge[I, C]
	}

	stack := []stackItem{{vertex: startVertex, edge: startEdge}}

	for len(stack) > 0 {
		// Pop vertex and edge from stack
		item := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		current := item.vertex
		incomingEdge := item.edge

		currentIdx := current.GetCustomDataIndex()
		currentData := &d.vertexData[currentIdx]

		// Skip if already visited
		if currentData.visited {
			continue
		}

		// Mark as visited and call callback
		currentData.visited = true
		callback(current, incomingEdge)

		// Add all unvisited neighbors to stack
		// We reverse the order to maintain the same traversal order as recursive version
		edges := current.GetEdges()
		for i := len(edges) - 1; i >= 0; i-- {
			neighbor := edges[i].GetTargetVertex()
			neighborIdx := neighbor.GetCustomDataIndex()
			neighborData := &d.vertexData[neighborIdx]

			if !neighborData.visited {
				neighborData.parent = current
				stack = append(stack, stackItem{vertex: neighbor, edge: &edges[i]})
			}
		}
	}
}

// dfsSearch performs DFS to find a path from start to target.
// Returns true if target is found, false otherwise.
// Uses an iterative approach with an explicit stack to avoid recursion.
// Always tracks parent pointers for potential path reconstruction.
func (d *DFS[I, C, V, E]) dfsSearch(start *Vertex[I, C], target interface{}) bool {
	// Use a stack to store vertices to visit
	stack := []*Vertex[I, C]{start}

	for len(stack) > 0 {
		// Pop vertex from stack
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		currentIdx := current.GetCustomDataIndex()
		currentData := &d.vertexData[currentIdx]

		// Skip if already visited
		if currentData.visited {
			continue
		}

		// Mark as visited
		currentData.visited = true

		// Check if we reached the target
		var found bool
		switch t := target.(type) {
		case *Vertex[I, C]:
			found = current == t
		case I:
			found = current.GetId() == t
		default:
			panic("invalid target type")
		}

		if found {
			return true
		}

		// Add all unvisited neighbors to stack
		// We reverse the order to maintain the same traversal order as recursive version
		edges := current.GetEdges()
		for i := len(edges) - 1; i >= 0; i-- {
			neighbor := edges[i].GetTargetVertex()
			neighborIdx := neighbor.GetCustomDataIndex()
			neighborData := &d.vertexData[neighborIdx]

			if !neighborData.visited {
				neighborData.parent = current
				stack = append(stack, neighbor)
			}
		}
	}

	return false
}

// FindCycles finds all cycles in the graph.
// Returns a slice of cycles, where each cycle is represented as a slice of vertex IDs.
// For directed graphs, this detects directed cycles.
// For undirected graphs, this detects any cycle (including simple back edges).
// Time complexity: O(V + E) where V is the number of vertices and E is the number of edges.
// Space complexity: O(V) where V is the number of vertices.
// WARNING: This function is not thread-safe and should not be called concurrently.
func (d *DFS[I, C, V, E]) FindCycles() [][]I {
	// Initialize vertex data for all vertices
	for i := range d.vertexData {
		d.vertexData[i].visited = false
		d.vertexData[i].parent = nil
		d.vertexData[i].visiting = false
	}

	var cycles [][]I
	visitedInCycles := make(map[I]bool) // Track vertices already part of found cycles

	// Check each unvisited vertex to find cycles
	for i := range d.graph.vertices {
		vertex := &d.graph.vertices[i]
		vertexIdx := vertex.GetCustomDataIndex()
		vertexData := &d.vertexData[vertexIdx]

		if !vertexData.visited && !visitedInCycles[vertex.GetId()] {
			cycle := d.findCycleFromVertex(vertex)
			if cycle != nil {
				cycles = append(cycles, cycle)
				// Mark all vertices in this cycle as visited to avoid duplicates
				for _, vertexId := range cycle {
					visitedInCycles[vertexId] = true
				}
			}
		}
	}

	return cycles
}

// HasCycle detects if the graph contains any cycles.
// This is a convenience method that returns true if FindCycles() finds any cycles.
// Time complexity: O(V + E) where V is the number of vertices and E is the number of edges.
// Space complexity: O(V) where V is the number of vertices.
// WARNING: This function is not thread-safe and should not be called concurrently.
func (d *DFS[I, C, V, E]) HasCycle() bool {
	cycles := d.FindCycles()
	return len(cycles) > 0
}

// findCycleFromVertex performs DFS from the given vertex to find cycles.
// Uses the "visiting" state to detect back edges in the current path.
// Returns the first cycle found as a slice of vertex IDs, or nil if no cycle.
func (d *DFS[I, C, V, E]) findCycleFromVertex(startVertex *Vertex[I, C]) []I {
	// Use a stack to store vertices and their state
	type stackItem struct {
		vertex  *Vertex[I, C]
		started bool // true if we've started processing this vertex
		path    []I  // current path from start vertex
	}

	stack := []stackItem{{vertex: startVertex, started: false, path: []I{}}}

	for len(stack) > 0 {
		// Pop vertex from stack
		item := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		current := item.vertex
		currentIdx := current.GetCustomDataIndex()
		currentData := &d.vertexData[currentIdx]

		if item.started {
			// We're finishing processing this vertex - mark as visited
			currentData.visiting = false
			currentData.visited = true
		} else {
			// Check if we've already visited this vertex
			if currentData.visited {
				continue
			}

			// Check if we're currently visiting this vertex (back edge = cycle)
			if currentData.visiting {
				// Found a cycle! Extract the cycle from the path
				cycleStart := -1
				for i, vertexId := range item.path {
					if vertexId == current.GetId() {
						cycleStart = i
						break
					}
				}
				if cycleStart >= 0 {
					// Return the cycle (from the back edge to the current vertex)
					cycle := make([]I, len(item.path)-cycleStart)
					copy(cycle, item.path[cycleStart:])
					return cycle
				}
			}

			// Mark as currently visiting
			currentData.visiting = true

			// Create new path with current vertex
			newPath := make([]I, len(item.path)+1)
			copy(newPath, item.path)
			newPath[len(item.path)] = current.GetId()

			// Push the vertex back to mark it as finished later
			stack = append(stack, stackItem{vertex: current, started: true, path: newPath})

			// Add all neighbors to stack
			edges := current.GetEdges()
			for i := len(edges) - 1; i >= 0; i-- {
				neighbor := edges[i].GetTargetVertex()
				neighborIdx := neighbor.GetCustomDataIndex()
				neighborData := &d.vertexData[neighborIdx]

				// Only process unvisited neighbors
				if !neighborData.visited {
					stack = append(stack, stackItem{vertex: neighbor, started: false, path: newPath})
				}
			}
		}
	}

	return nil
}
