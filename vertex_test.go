package graph

import (
	"testing"
)

func TestVertex(t *testing.T) {
	t.Run("Vertex with int ID and float cost", func(t *testing.T) {
		vertex := &Vertex[int, float64]{
			id:              1,
			customDataIndex: 0,
			edges:           []Edge[int, float64]{},
		}

		if vertex.GetId() != 1 {
			t.Errorf("Expected ID 1, got %v", vertex.GetId())
		}

		if vertex.GetCustomDataIndex() != 0 {
			t.Errorf("Expected custom data index 0, got %v", vertex.GetCustomDataIndex())
		}

		if len(vertex.GetEdges()) != 0 {
			t.Errorf("Expected empty edges slice, got %v", vertex.GetEdges())
		}
	})

	t.Run("Vertex with string ID and int cost", func(t *testing.T) {
		vertex := &Vertex[string, int]{
			id:              "A",
			customDataIndex: 5,
			edges:           []Edge[string, int]{},
		}

		if vertex.GetId() != "A" {
			t.Errorf("Expected ID 'A', got %v", vertex.GetId())
		}

		if vertex.GetCustomDataIndex() != 5 {
			t.Errorf("Expected custom data index 5, got %v", vertex.GetCustomDataIndex())
		}
	})

	t.Run("Vertex with edges", func(t *testing.T) {
		targetVertex := &Vertex[int, float64]{
			id:              2,
			customDataIndex: 1,
			edges:           []Edge[int, float64]{},
		}

		edge := Edge[int, float64]{
			cost:            10.5,
			targetVertex:    targetVertex,
			customDataIndex: 0,
		}

		vertex := &Vertex[int, float64]{
			id:              1,
			customDataIndex: 0,
			edges:           []Edge[int, float64]{edge},
		}

		edges := vertex.GetEdges()
		if len(edges) != 1 {
			t.Errorf("Expected 1 edge, got %d", len(edges))
		}

		if edges[0].GetCost() != 10.5 {
			t.Errorf("Expected edge cost 10.5, got %v", edges[0].GetCost())
		}

		if edges[0].GetTargetVertex().GetId() != 2 {
			t.Errorf("Expected target vertex ID 2, got %v", edges[0].GetTargetVertex().GetId())
		}
	})
}
