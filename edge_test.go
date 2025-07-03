package graph

import (
	"testing"
)

func TestEdge(t *testing.T) {
	t.Run("Edge with int ID and float cost", func(t *testing.T) {
		targetVertex := &Vertex[int, float64]{
			id:              2,
			customDataIndex: 1,
			edges:           []Edge[int, float64]{},
		}

		edge := &Edge[int, float64]{
			cost:            15.5,
			targetVertex:    targetVertex,
			customDataIndex: 0,
		}

		if edge.GetCost() != 15.5 {
			t.Errorf("Expected cost 15.5, got %v", edge.GetCost())
		}

		if edge.GetTargetVertex().GetId() != 2 {
			t.Errorf("Expected target vertex ID 2, got %v", edge.GetTargetVertex().GetId())
		}

		if edge.GetCustomDataIndex() != 0 {
			t.Errorf("Expected custom data index 0, got %v", edge.GetCustomDataIndex())
		}
	})

	t.Run("Edge cost modification", func(t *testing.T) {
		targetVertex := &Vertex[string, int]{
			id:              "B",
			customDataIndex: 1,
			edges:           []Edge[string, int]{},
		}

		edge := &Edge[string, int]{
			cost:            10,
			targetVertex:    targetVertex,
			customDataIndex: 0,
		}

		edge.SetCost(20)
		if edge.GetCost() != 20 {
			t.Errorf("Expected cost 20 after SetCost, got %v", edge.GetCost())
		}
	})

	t.Run("Edge clone", func(t *testing.T) {
		targetVertex := &Vertex[int, float64]{
			id:              3,
			customDataIndex: 2,
			edges:           []Edge[int, float64]{},
		}

		originalEdge := &Edge[int, float64]{
			cost:            25.0,
			targetVertex:    targetVertex,
			customDataIndex: 1,
		}

		clonedEdge := originalEdge.Clone()

		// Test that clone has same values
		if clonedEdge.GetCost() != originalEdge.GetCost() {
			t.Errorf("Expected cloned cost %v, got %v", originalEdge.GetCost(), clonedEdge.GetCost())
		}

		if clonedEdge.GetTargetVertex().GetId() != originalEdge.GetTargetVertex().GetId() {
			t.Errorf("Expected cloned target vertex ID %v, got %v", originalEdge.GetTargetVertex().GetId(), clonedEdge.GetTargetVertex().GetId())
		}

		if clonedEdge.GetCustomDataIndex() != originalEdge.GetCustomDataIndex() {
			t.Errorf("Expected cloned custom data index %v, got %v", originalEdge.GetCustomDataIndex(), clonedEdge.GetCustomDataIndex())
		}

		// Test that modifying original doesn't affect clone
		originalEdge.SetCost(30.0)
		if clonedEdge.GetCost() == 30.0 {
			t.Error("Modifying original edge cost affected cloned edge")
		}
	})

	t.Run("Edge with different types", func(t *testing.T) {
		targetVertex := &Vertex[rune, uint]{
			id:              'A',
			customDataIndex: 1,
			edges:           []Edge[rune, uint]{},
		}

		edge := &Edge[rune, uint]{
			cost:            100,
			targetVertex:    targetVertex,
			customDataIndex: 0,
		}

		if edge.GetCost() != 100 {
			t.Errorf("Expected cost 100, got %v", edge.GetCost())
		}

		if edge.GetTargetVertex().GetId() != 'A' {
			t.Errorf("Expected target vertex ID 'A', got %v", edge.GetTargetVertex().GetId())
		}
	})
}
