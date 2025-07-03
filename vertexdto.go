package graph

// VertexDto defines the interface for vertex data transfer objects.
// This interface provides methods to get and set vertex ID and custom data.
// The generic type I represents the vertex ID type, and V represents the custom vertex data type.
type VertexDto[I Id, V any] interface {
	GetId() I       // Returns the vertex identifier
	SetId(id I)     // Sets the vertex identifier
	GetData() V     // Returns the custom vertex data
	SetData(data V) // Sets the custom vertex data
}

// BasicVertexDto is a simple implementation of the VertexDto interface.
// It provides a basic structure for transferring vertex data with JSON serialization support.
// The generic type I represents the vertex ID type, and V represents the custom vertex data type.
type BasicVertexDto[I Id, V any] struct {
	Id   I `json:"id"`   // Vertex identifier with JSON tag
	Data V `json:"data"` // Custom vertex data with JSON tag
}

// GetId returns the vertex identifier.
// Implements the VertexDto interface.
func (dto *BasicVertexDto[I, V]) GetId() I {
	return dto.Id
}

// SetId sets the vertex identifier.
// Implements the VertexDto interface.
func (dto *BasicVertexDto[I, V]) SetId(id I) {
	dto.Id = id
}

// GetData returns the custom vertex data.
// Implements the VertexDto interface.
func (dto *BasicVertexDto[I, V]) GetData() V {
	return dto.Data
}

// SetData sets the custom vertex data.
// Implements the VertexDto interface.
func (dto *BasicVertexDto[I, V]) SetData(data V) {
	dto.Data = data
}
