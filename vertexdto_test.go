package graph

import (
	"testing"
)

func TestVertexDto(t *testing.T) {
	t.Run("BasicVertexDto with int ID and string data", func(t *testing.T) {
		dto := &BasicVertexDto[int, string]{
			Id:   1,
			Data: "test data",
		}

		if dto.GetId() != 1 {
			t.Errorf("Expected ID 1, got %v", dto.GetId())
		}

		if dto.GetData() != "test data" {
			t.Errorf("Expected data 'test data', got %v", dto.GetData())
		}
	})

	t.Run("BasicVertexDto with string ID and int data", func(t *testing.T) {
		dto := &BasicVertexDto[string, int]{
			Id:   "vertex1",
			Data: 42,
		}

		if dto.GetId() != "vertex1" {
			t.Errorf("Expected ID 'vertex1', got %v", dto.GetId())
		}

		if dto.GetData() != 42 {
			t.Errorf("Expected data 42, got %v", dto.GetData())
		}
	})

	t.Run("BasicVertexDto setter methods", func(t *testing.T) {
		dto := &BasicVertexDto[int, string]{}

		dto.SetId(100)
		if dto.GetId() != 100 {
			t.Errorf("Expected ID 100 after SetId, got %v", dto.GetId())
		}

		dto.SetData("new data")
		if dto.GetData() != "new data" {
			t.Errorf("Expected data 'new data' after SetData, got %v", dto.GetData())
		}
	})

	t.Run("BasicVertexDto with rune ID and float data", func(t *testing.T) {
		dto := &BasicVertexDto[rune, float64]{
			Id:   'A',
			Data: 3.14159,
		}

		if dto.GetId() != 'A' {
			t.Errorf("Expected ID 'A', got %v", dto.GetId())
		}

		if dto.GetData() != 3.14159 {
			t.Errorf("Expected data 3.14159, got %v", dto.GetData())
		}
	})

	t.Run("BasicVertexDto with struct data", func(t *testing.T) {
		type TestData struct {
			Name  string
			Value int
		}

		testData := TestData{Name: "test", Value: 123}
		dto := &BasicVertexDto[int, TestData]{
			Id:   999,
			Data: testData,
		}

		if dto.GetId() != 999 {
			t.Errorf("Expected ID 999, got %v", dto.GetId())
		}

		retrievedData := dto.GetData()
		if retrievedData.Name != "test" || retrievedData.Value != 123 {
			t.Errorf("Expected data {Name: 'test', Value: 123}, got %v", retrievedData)
		}
	})
}
