package graph

import (
	"testing"
)

func TestEdgeDto(t *testing.T) {
	t.Run("BasicEdgeDto with int IDs and float cost", func(t *testing.T) {
		dto := &BasicEdgeDto[int, float64, string]{
			Origin: 1,
			Target: 2,
			Cost:   10.5,
			Data:   "edge data",
		}

		if dto.GetOrigin() != 1 {
			t.Errorf("Expected origin 1, got %v", dto.GetOrigin())
		}

		if dto.GetTarget() != 2 {
			t.Errorf("Expected target 2, got %v", dto.GetTarget())
		}

		if dto.GetCost() != 10.5 {
			t.Errorf("Expected cost 10.5, got %v", dto.GetCost())
		}

		if dto.GetData() != "edge data" {
			t.Errorf("Expected data 'edge data', got %v", dto.GetData())
		}
	})

	t.Run("BasicEdgeDto with string IDs and int cost", func(t *testing.T) {
		dto := &BasicEdgeDto[string, int, bool]{
			Origin: "A",
			Target: "B",
			Cost:   15,
			Data:   true,
		}

		if dto.GetOrigin() != "A" {
			t.Errorf("Expected origin 'A', got %v", dto.GetOrigin())
		}

		if dto.GetTarget() != "B" {
			t.Errorf("Expected target 'B', got %v", dto.GetTarget())
		}

		if dto.GetCost() != 15 {
			t.Errorf("Expected cost 15, got %v", dto.GetCost())
		}

		if dto.GetData() != true {
			t.Errorf("Expected data true, got %v", dto.GetData())
		}
	})

	t.Run("BasicEdgeDto setter methods", func(t *testing.T) {
		dto := &BasicEdgeDto[int, float64, string]{}

		dto.SetOrigin(10)
		if dto.GetOrigin() != 10 {
			t.Errorf("Expected origin 10 after SetOrigin, got %v", dto.GetOrigin())
		}

		dto.SetTarget(20)
		if dto.GetTarget() != 20 {
			t.Errorf("Expected target 20 after SetTarget, got %v", dto.GetTarget())
		}

		dto.SetCost(25.5)
		if dto.GetCost() != 25.5 {
			t.Errorf("Expected cost 25.5 after SetCost, got %v", dto.GetCost())
		}

		dto.SetData("new edge data")
		if dto.GetData() != "new edge data" {
			t.Errorf("Expected data 'new edge data' after SetData, got %v", dto.GetData())
		}
	})

	t.Run("BasicEdgeDto with rune IDs and uint cost", func(t *testing.T) {
		dto := &BasicEdgeDto[rune, uint, int]{
			Origin: 'X',
			Target: 'Y',
			Cost:   100,
			Data:   42,
		}

		if dto.GetOrigin() != 'X' {
			t.Errorf("Expected origin 'X', got %v", dto.GetOrigin())
		}

		if dto.GetTarget() != 'Y' {
			t.Errorf("Expected target 'Y', got %v", dto.GetTarget())
		}

		if dto.GetCost() != 100 {
			t.Errorf("Expected cost 100, got %v", dto.GetCost())
		}

		if dto.GetData() != 42 {
			t.Errorf("Expected data 42, got %v", dto.GetData())
		}
	})

	t.Run("BasicEdgeDto with struct data", func(t *testing.T) {
		type EdgeData struct {
			Weight int
			Label  string
		}

		edgeData := EdgeData{Weight: 5, Label: "important"}
		dto := &BasicEdgeDto[int, float64, EdgeData]{
			Origin: 1,
			Target: 3,
			Cost:   7.5,
			Data:   edgeData,
		}

		if dto.GetOrigin() != 1 {
			t.Errorf("Expected origin 1, got %v", dto.GetOrigin())
		}

		if dto.GetTarget() != 3 {
			t.Errorf("Expected target 3, got %v", dto.GetTarget())
		}

		if dto.GetCost() != 7.5 {
			t.Errorf("Expected cost 7.5, got %v", dto.GetCost())
		}

		retrievedData := dto.GetData()
		if retrievedData.Weight != 5 || retrievedData.Label != "important" {
			t.Errorf("Expected data {Weight: 5, Label: 'important'}, got %v", retrievedData)
		}
	})
}
