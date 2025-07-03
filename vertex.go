package graph

// Vertex represents a node in the graph with an identifier, custom data index, and outgoing edges.
// The generic type I represents the vertex ID type, and C represents the edge cost type.
type Vertex[I Id, C Cost] struct {
	id              I            // Unique identifier for this vertex
	customDataIndex int          // Index into the graph's custom vertex data array
	edges           []Edge[I, C] // List of outgoing edges from this vertex
}

// GetId returns the unique identifier of this vertex.
// Returns the vertex ID of type I.
func (v *Vertex[I, C]) GetId() I {
	return v.id
}

// GetCustomDataIndex returns the index into the graph's custom vertex data array.
// This index is used to retrieve associated custom data for this vertex.
func (v *Vertex[I, C]) GetCustomDataIndex() int {
	return v.customDataIndex
}

// GetEdges returns a slice of all outgoing edges from this vertex.
// Returns a copy of the edges slice containing Edge[I, C] objects.
func (v *Vertex[I, C]) GetEdges() []Edge[I, C] {
	return v.edges
}
