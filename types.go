package graph

// SInt represents signed integer types that can be used as vertex IDs or edge costs.
// This includes all standard signed integer types from int8 to int64.
type SInt interface {
	int | int8 | int16 | int32 | int64
}

// UInt represents unsigned integer types that can be used as vertex IDs or edge costs.
// This includes all standard unsigned integer types from uint8 to uint64.
type UInt interface {
	uint | uint8 | uint16 | uint32 | uint64
}

// Float represents floating-point types that can be used as edge costs.
// This includes float32 and float64 types.
type Float interface {
	float32 | float64
}

// Id represents the constraint for vertex identifier types.
// Vertex IDs can be signed integers, unsigned integers, strings, or runes.
// This provides flexibility in choosing appropriate identifier types for different use cases.
type Id interface {
	SInt | UInt | string | rune
}

// Cost represents the constraint for edge cost types.
// Edge costs can be signed integers, unsigned integers, or floating-point numbers.
// This allows for both integer and decimal cost representations depending on the application needs.
type Cost interface {
	SInt | UInt | Float
}
