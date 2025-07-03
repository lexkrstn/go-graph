package graph

// Edge represents a directed connection between two vertices in the graph.
// The generic type I represents the vertex ID type, and C represents the edge cost type.
type Edge[I Id, C Cost] struct {
	cost            C             // Cost or weight associated with this edge
	targetVertex    *Vertex[I, C] // Pointer to the destination vertex
	customDataIndex int           // Index into the graph's custom edge data array
}

// GetCost returns the cost or weight associated with this edge.
// Returns the edge cost of type C.
func (e *Edge[I, C]) GetCost() C {
	return e.cost
}

// SetCost updates the cost or weight associated with this edge.
// Takes a new cost value of type C.
func (e *Edge[I, C]) SetCost(cost C) {
	e.cost = cost
}

// GetTargetVertex returns a pointer to the destination vertex of this edge.
// Returns a pointer to the Vertex[I, C] that this edge points to.
func (e *Edge[I, C]) GetTargetVertex() *Vertex[I, C] {
	return e.targetVertex
}

// GetCustomDataIndex returns the index into the graph's custom edge data array.
// This index is used to retrieve associated custom data for this edge.
func (e *Edge[I, C]) GetCustomDataIndex() int {
	return e.customDataIndex
}

// Clone creates a deep copy of this edge.
// Returns a new Edge[I, C] with the same cost, target vertex pointer, and custom data index.
// Note: The target vertex pointer is shared between the original and cloned edge.
func (e *Edge[I, C]) Clone() Edge[I, C] {
	return Edge[I, C]{
		cost:            e.cost,
		targetVertex:    e.targetVertex,
		customDataIndex: e.customDataIndex,
	}
}
