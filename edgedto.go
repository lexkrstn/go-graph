package graph

// EdgeDto defines the interface for edge data transfer objects.
// This interface provides methods to get and set edge origin, target, cost, and custom data.
// The generic type I represents the vertex ID type, C represents the edge cost type, and E represents the custom edge data type.
type EdgeDto[I Id, C Cost, E any] interface {
	GetOrigin() I   // Returns the origin vertex identifier
	SetOrigin(id I) // Sets the origin vertex identifier
	GetTarget() I   // Returns the target vertex identifier
	SetTarget(id I) // Sets the target vertex identifier
	GetCost() C     // Returns the edge cost
	SetCost(cost C) // Sets the edge cost
	GetData() E     // Returns the custom edge data
	SetData(data E) // Sets the custom edge data
}

// BasicEdgeDto is a simple implementation of the EdgeDto interface.
// It provides a basic structure for transferring edge data with JSON serialization support.
// The generic type I represents the vertex ID type, C represents the edge cost type, and E represents the custom edge data type.
type BasicEdgeDto[I Id, C Cost, E any] struct {
	Origin I `json:"origin"` // Origin vertex identifier with JSON tag
	Target I `json:"target"` // Target vertex identifier with JSON tag
	Cost   C `json:"cost"`   // Edge cost with JSON tag
	Data   E `json:"data"`   // Custom edge data with JSON tag
}

// GetOrigin returns the origin vertex identifier.
// Implements the EdgeDto interface.
func (dto *BasicEdgeDto[I, C, E]) GetOrigin() I {
	return dto.Origin
}

// SetOrigin sets the origin vertex identifier.
// Implements the EdgeDto interface.
func (dto *BasicEdgeDto[I, C, E]) SetOrigin(id I) {
	dto.Origin = id
}

// GetTarget returns the target vertex identifier.
// Implements the EdgeDto interface.
func (dto *BasicEdgeDto[I, C, E]) GetTarget() I {
	return dto.Target
}

// SetTarget sets the target vertex identifier.
// Implements the EdgeDto interface.
func (dto *BasicEdgeDto[I, C, E]) SetTarget(id I) {
	dto.Target = id
}

// GetCost returns the edge cost.
// Implements the EdgeDto interface.
func (dto *BasicEdgeDto[I, C, E]) GetCost() C {
	return dto.Cost
}

// SetCost sets the edge cost.
// Implements the EdgeDto interface.
func (dto *BasicEdgeDto[I, C, E]) SetCost(cost C) {
	dto.Cost = cost
}

// GetData returns the custom edge data.
// Implements the EdgeDto interface.
func (dto *BasicEdgeDto[I, C, E]) GetData() E {
	return dto.Data
}

// SetData sets the custom edge data.
// Implements the EdgeDto interface.
func (dto *BasicEdgeDto[I, C, E]) SetData(data E) {
	dto.Data = data
}
