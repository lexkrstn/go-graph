package graph

import (
	"testing"
)

// Test type constraints with various types
func TestTypeConstraints(t *testing.T) {
	// Test SInt constraint
	_ = testSInt[int](42)
	_ = testSInt[int8](42)
	_ = testSInt[int16](42)
	_ = testSInt[int32](42)
	_ = testSInt[int64](42)

	// Test UInt constraint
	_ = testUInt[uint](42)
	_ = testUInt[uint8](42)
	_ = testUInt[uint16](42)
	_ = testUInt[uint32](42)
	_ = testUInt[uint64](42)

	// Test Float constraint
	_ = testFloat[float32](3.14)
	_ = testFloat[float64](3.14)

	// Test Id constraint
	_ = testId[int](42)
	_ = testId[uint](42)
	_ = testId[string]("test")
	_ = testId[rune]('a')

	// Test Cost constraint
	_ = testCost[int](42)
	_ = testCost[uint](42)
	_ = testCost[float32](3.14)
	_ = testCost[float64](3.14)
}

func testSInt[T SInt](val T) T {
	return val
}

func testUInt[T UInt](val T) T {
	return val
}

func testFloat[T Float](val T) T {
	return val
}

func testId[T Id](val T) T {
	return val
}

func testCost[T Cost](val T) T {
	return val
}
