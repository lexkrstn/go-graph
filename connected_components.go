package graph

// The data that is attached to the vertices by the ConnectedComponents algorithm.
type connectedComponentsVertexData[I Id] struct {
	visited     bool
	componentId int
}

// The ConnectedComponents algorithm Use-Case (aka Command) object.
// It contains the precomputed connected components data and provides
// methods to query the results without recomputing.
type ConnectedComponents[I Id, C Cost, V any, E any] struct {
	graph      *Graph[I, C, V, E]
	components [][]I
}

// FindConnectedComponents finds all connected components in the graph.
// Returns a ConnectedComponents instance with precomputed results.
// Time complexity: O(V + E) where V is the number of vertices and E is the number of edges.
// Space complexity: O(V) where V is the number of vertices.
func FindConnectedComponents[I Id, C Cost, V any, E any](graph *Graph[I, C, V, E]) *ConnectedComponents[I, C, V, E] {
	vertexData := make([]connectedComponentsVertexData[I], len(graph.vertices))
	cc := &ConnectedComponents[I, C, V, E]{
		graph: graph,
	}

	// Initialize vertex data for all vertices
	for i := range vertexData {
		vertexData[i].visited = false
		vertexData[i].componentId = -1
	}

	var components [][]I
	componentId := 0

	// Visit all vertices to find connected components
	for i := range cc.graph.vertices {
		vertex := &cc.graph.vertices[i]
		vertexIdx := vertex.GetCustomDataIndex()

		// If vertex hasn't been visited, start a new component
		if !vertexData[vertexIdx].visited {
			component := dfs(cc, vertex, vertexData, componentId)
			if len(component) > 0 {
				components = append(components, component)
				componentId++
			}
		}
	}

	cc.components = components
	return cc
}

// GetComponents returns the precomputed connected components.
// Returns a slice of slices, where each inner slice contains the vertex IDs
// that belong to the same connected component.
// Time complexity: O(1) - returns precomputed data.
func (cc *ConnectedComponents[I, C, V, E]) GetComponents() [][]I {
	return cc.components
}

// dfs performs depth-first search starting from the given vertex.
// It marks all reachable vertices as visited and assigns them the same component ID.
// For directed graphs, this considers both incoming and outgoing edges to find
// all vertices in the same strongly connected component.
// Returns a slice of vertex IDs in the connected component.
func dfs[I Id, C Cost, V any, E any](
	cc *ConnectedComponents[I, C, V, E],
	vertex *Vertex[I, C],
	data []connectedComponentsVertexData[I],
	componentId int,
) []I {
	vertexIdx := vertex.GetCustomDataIndex()
	vertexData := &data[vertexIdx]

	// Mark as visited and assign component ID
	vertexData.visited = true
	vertexData.componentId = componentId

	component := []I{vertex.GetId()}

	// Visit all neighbors recursively (outgoing edges)
	for _, edge := range vertex.GetEdges() {
		neighbor := edge.GetTargetVertex()
		neighborIdx := neighbor.GetCustomDataIndex()
		neighborData := &data[neighborIdx]

		if !neighborData.visited {
			neighborComponent := dfs(cc, neighbor, data, componentId)
			component = append(component, neighborComponent...)
		}
	}

	// For directed graphs, also check incoming edges by searching through all vertices
	// This ensures we find all vertices that can reach the current vertex
	for i := range cc.graph.vertices {
		otherVertex := &cc.graph.vertices[i]
		otherIdx := otherVertex.GetCustomDataIndex()
		otherData := &data[otherIdx]

		if !otherData.visited {
			// Check if this vertex has an edge to our current vertex
			hasEdgeToCurrent := false
			for _, edge := range otherVertex.GetEdges() {
				if edge.GetTargetVertex() == vertex {
					hasEdgeToCurrent = true
					break
				}
			}

			if hasEdgeToCurrent {
				neighborComponent := dfs(cc, otherVertex, data, componentId)
				component = append(component, neighborComponent...)
			}
		}
	}

	return component
}

// GetComponentCount returns the number of connected components in the graph.
// Time complexity: O(1) - returns precomputed data.
func (cc *ConnectedComponents[I, C, V, E]) GetComponentCount() int {
	return len(cc.components)
}

// IsConnected checks if the graph is connected (has only one connected component).
// Returns true if the graph is connected, false otherwise.
// Time complexity: O(1) - returns precomputed data.
func (cc *ConnectedComponents[I, C, V, E]) IsConnected() bool {
	return len(cc.components) == 1
}

// GetComponentForVertex returns the connected component that contains the given vertex.
// Returns a slice of vertex IDs in the same component as the given vertex.
// Returns nil if the vertex is not found in the graph.
// Time complexity: O(V) where V is the number of vertices in the component.
func (cc *ConnectedComponents[I, C, V, E]) GetComponentForVertex(vertexId I) []I {
	// Check if vertex exists
	_, err := cc.graph.GetVertexById(vertexId)
	if err != nil {
		return nil // Vertex not found
	}

	// Search through precomputed components to find the one containing the vertex
	for _, component := range cc.components {
		for _, id := range component {
			if id == vertexId {
				return component
			}
		}
	}

	return nil
}
