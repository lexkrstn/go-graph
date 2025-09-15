package graph

import (
	"math"
	"reflect"
)

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

// Implements std::numeric_limits<T>::max() in a Go way.
func assignMaxNumber(v interface{}) {
	val := reflect.ValueOf(v).Elem()
	switch val.Kind() {
	case reflect.Int:
		val.SetInt(math.MaxInt)
	case reflect.Int8:
		val.SetInt(math.MaxInt8)
	case reflect.Int16:
		val.SetInt(math.MaxInt16)
	case reflect.Int32:
		val.SetInt(math.MaxInt32)
	case reflect.Int64:
		val.SetInt(math.MaxInt64)
	case reflect.Uint:
		val.SetUint(math.MaxUint)
	case reflect.Uint8:
		val.SetUint(math.MaxUint8)
	case reflect.Uint16:
		val.SetUint(math.MaxUint16)
	case reflect.Uint32:
		val.SetUint(math.MaxUint32)
	case reflect.Uint64:
		val.SetUint(math.MaxUint64)
	case reflect.Float32:
		val.SetFloat(math.MaxFloat32)
	case reflect.Float64:
		val.SetFloat(math.MaxFloat64)
	default:
		panic("Unsupported type")
	}
}
