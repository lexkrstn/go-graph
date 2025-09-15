package graph

import "errors"

// Graph represents a directed graph with vertices and edges.
// The graph encapsulates edges and vertices with support for custom data types.
// Generic types: I (Id), C (Cost), V (Vertex data), E (Edge data).
type Graph[I Id, C Cost, V any, E any] struct {
	vertices         []Vertex[I, C] // Array of all vertices in the graph
	idToIndex        map[I]int      // Mapping from vertex ID to array index for O(1) lookups
	customVertexData []V            // Array of custom data associated with each vertex
	customEdgeData   []E            // Array of custom data associated with each edge
	edgeCount        int            // Total number of directed edges in the graph
	biEdgeCount      int            // Number of bidirectional edges (unique vertex pairs)
}

// GetVertexCount returns the total number of vertices in the graph.
// Returns the length of the vertices slice.
func (g *Graph[I, C, V, E]) GetVertexCount() int {
	return len(g.vertices)
}

// GetEdgeCount returns the total number of directed edges in the graph.
// This includes all edges, including both directions of bidirectional edges.
func (g *Graph[I, C, V, E]) GetEdgeCount() int {
	return g.edgeCount
}

// GetBiEdgeCount returns the number of bidirectional edges in the graph.
// A bidirectional edge is counted as one edge between a pair of vertices,
// regardless of whether there are edges in both directions.
func (g *Graph[I, C, V, E]) GetBiEdgeCount() int {
	return g.biEdgeCount
}

// GetVertexById retrieves a vertex by its unique identifier.
// Returns a pointer to the vertex if found, or an error if the ID doesn't exist.
// Time complexity: O(1) due to the idToIndex map.
func (g *Graph[I, C, V, E]) GetVertexById(id I) (*Vertex[I, C], error) {
	idx, ok := g.idToIndex[id]
	if !ok {
		return nil, errors.New("vertex id not found")
	}
	return &g.vertices[idx], nil
}

// GetVertexByIndex retrieves a vertex by its array index.
// Returns a pointer to the vertex if the index is valid, or an error if out of range.
// Time complexity: O(1) array access.
func (g *Graph[I, C, V, E]) GetVertexByIndex(idx int) (*Vertex[I, C], error) {
	if idx < 0 || idx >= len(g.vertices) {
		return nil, errors.New("index out of range")
	}
	return &g.vertices[idx], nil
}

// GetVertexData retrieves the custom data associated with a vertex.
// Returns a pointer to the vertex's custom data if the vertex is valid, or an error if nil.
func (g *Graph[I, C, V, E]) GetVertexData(vertex *Vertex[I, C]) (*V, error) {
	if vertex == nil {
		return nil, errors.New("vertex ptr is nil")
	}
	return &g.customVertexData[vertex.customDataIndex], nil
}

// GetEdgeData retrieves the custom data associated with an edge.
// Returns a pointer to the edge's custom data if the edge is valid, or an error if nil.
func (g *Graph[I, C, V, E]) GetEdgeData(edge *Edge[I, C]) (*E, error) {
	if edge == nil {
		return nil, errors.New("edge ptr is nil")
	}
	return &g.customEdgeData[edge.customDataIndex], nil
}

// GetAllVertices returns all vertices in the graph as DTOs.
// Takes a factory function to create new vertex DTOs.
// Returns a slice of VertexDto objects containing all vertex data.
func (g *Graph[I, C, V, E]) GetAllVertices(newVertex func() VertexDto[I, V]) []VertexDto[I, V] {
	dtos := make([]VertexDto[I, V], len(g.vertices))
	for i := range g.vertices {
		dtos[i] = newVertex()
		dtos[i].SetId(g.vertices[i].id)
		dtos[i].SetData(g.customVertexData[g.vertices[i].customDataIndex])
	}
	return dtos
}

// GetAllEdges returns all directed edges in the graph as DTOs.
// Takes a factory function to create new edge DTOs.
// Returns a slice of EdgeDto objects containing all edge data.
// Note: This includes all edges, so bidirectional connections appear twice.
func (g *Graph[I, C, V, E]) GetAllEdges(newEdge func() EdgeDto[I, C, E]) []EdgeDto[I, C, E] {
	dtos := make([]EdgeDto[I, C, E], g.edgeCount)
	k := 0
	for i := range g.vertices {
		for j := range g.vertices[i].edges {
			dtos[k] = newEdge()
			dtos[k].SetOrigin(g.vertices[i].id)
			dtos[k].SetTarget(g.vertices[i].edges[j].targetVertex.id)
			dtos[k].SetCost(g.vertices[i].edges[j].cost)
			dtos[k].SetData(g.customEdgeData[g.vertices[i].edges[j].customDataIndex])
			k++
		}
	}
	return dtos
}

// GetAllBiEdges returns all bidirectional edges in the graph as DTOs.
// Takes a factory function to create new edge DTOs.
// Returns a slice of EdgeDto objects where each bidirectional connection appears only once.
// Uses a map to deduplicate edges between the same vertex pairs.
func (g *Graph[I, C, V, E]) GetAllBiEdges(newEdge func() EdgeDto[I, C, E]) []EdgeDto[I, C, E] {
	dtos := make([]EdgeDto[I, C, E], g.biEdgeCount)
	existing := make(map[biEdgeKey[I]]struct{}, g.biEdgeCount)
	k := 0
	var key biEdgeKey[I]
	for i := range g.vertices {
		for j := range g.vertices[i].edges {
			key.origin = g.vertices[i].id
			key.target = g.vertices[i].edges[j].targetVertex.id
			if key.origin > key.target {
				key.target, key.origin = key.origin, key.target
			}
			if _, exists := existing[key]; exists {
				continue
			}
			dtos[k] = newEdge()
			dtos[k].SetOrigin(key.origin)
			dtos[k].SetTarget(key.target)
			dtos[k].SetCost(g.vertices[i].edges[j].cost)
			dtos[k].SetData(g.customEdgeData[g.vertices[i].edges[j].customDataIndex])
			existing[key] = struct{}{}
			k++
		}
	}
	return dtos
}

// VisitEdges applies a visitor function to every edge in the graph.
// The visitor function receives both the source vertex and the edge.
// This allows for edge traversal with access to both vertex and edge data.
func (g *Graph[I, C, V, E]) VisitEdges(visitor func(*Vertex[I, C], *Edge[I, C])) {
	for i := range g.vertices {
		for j := range g.vertices[i].edges {
			visitor(&g.vertices[i], &g.vertices[i].edges[j])
		}
	}
}

// SomeEdges checks if any edge satisfies the given predicate.
// Returns true if at least one edge matches the predicate, false otherwise.
// Stops iteration as soon as a matching edge is found.
func (g *Graph[I, C, V, E]) SomeEdges(predicate func(*Vertex[I, C], *Edge[I, C]) bool) bool {
	for i := range g.vertices {
		for j := range g.vertices[i].edges {
			if predicate(&g.vertices[i], &g.vertices[i].edges[j]) {
				return true
			}
		}
	}
	return false
}

// EveryEdge checks if all edges satisfy the given predicate.
// Returns true if every edge matches the predicate, false otherwise.
// Stops iteration as soon as a non-matching edge is found.
func (g *Graph[I, C, V, E]) EveryEdge(predicate func(*Vertex[I, C], *Edge[I, C]) bool) bool {
	for i := range g.vertices {
		for j := range g.vertices[i].edges {
			if !predicate(&g.vertices[i], &g.vertices[i].edges[j]) {
				return false
			}
		}
	}
	return true
}

// VisitVertices applies a visitor function to every vertex in the graph.
// The visitor function receives a pointer to each vertex.
// This allows for vertex traversal with access to vertex data.
func (g *Graph[I, C, V, E]) VisitVertices(visitor func(*Vertex[I, C])) {
	for i := range g.vertices {
		visitor(&g.vertices[i])
	}
}

// SomeVertices checks if any vertex satisfies the given predicate.
// Returns true if at least one vertex matches the predicate, false otherwise.
// Stops iteration as soon as a matching vertex is found.
func (g *Graph[I, C, V, E]) SomeVertices(predicate func(*Vertex[I, C]) bool) bool {
	for i := range g.vertices {
		if predicate(&g.vertices[i]) {
			return true
		}
	}
	return false
}

// EveryVertex checks if all vertices satisfy the given predicate.
// Returns true if every vertex matches the predicate, false otherwise.
// Stops iteration as soon as a non-matching vertex is found.
func (g *Graph[I, C, V, E]) EveryVertex(predicate func(*Vertex[I, C]) bool) bool {
	for i := range g.vertices {
		if !predicate(&g.vertices[i]) {
			return false
		}
	}
	return true
}
