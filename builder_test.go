package graph

import (
	"testing"
)

func TestBuilder(t *testing.T) {
	t.Run("Empty builder", func(t *testing.T) {
		builder := &Builder[int, float64, string, bool]{}

		if builder.edgeCount != 0 {
			t.Errorf("Expected edge count 0, got %d", builder.edgeCount)
		}

		if builder.vertexCount != 0 {
			t.Errorf("Expected vertex count 0, got %d", builder.vertexCount)
		}
	})

	t.Run("Add single edge", func(t *testing.T) {
		builder := &Builder[int, float64, string, bool]{}

		builder.AddEdge(1, 2, 10.5, true)

		if builder.edgeCount != 1 {
			t.Errorf("Expected edge count 1, got %d", builder.edgeCount)
		}
	})

	t.Run("Add multiple edges", func(t *testing.T) {
		builder := &Builder[int, float64, string, bool]{}

		builder.AddEdge(1, 2, 10.5, true)
		builder.AddEdge(2, 3, 15.0, false)
		builder.AddEdge(1, 3, 20.0, true)

		if builder.edgeCount != 3 {
			t.Errorf("Expected edge count 3, got %d", builder.edgeCount)
		}
	})

	t.Run("Add bidirectional edge", func(t *testing.T) {
		builder := &Builder[int, float64, string, bool]{}

		builder.AddBiEdge(1, 2, 10.5, true)

		if builder.edgeCount != 2 {
			t.Errorf("Expected edge count 2, got %d", builder.edgeCount)
		}
	})

	t.Run("Add vertex", func(t *testing.T) {
		builder := &Builder[int, float64, string, bool]{}

		builder.AddVertex(1, "vertex1")

		if builder.vertexCount != 1 {
			t.Errorf("Expected vertex count 1, got %d", builder.vertexCount)
		}
	})

	t.Run("Add multiple vertices", func(t *testing.T) {
		builder := &Builder[int, float64, string, bool]{}

		builder.AddVertex(1, "vertex1")
		builder.AddVertex(2, "vertex2")
		builder.AddVertex(3, "vertex3")

		if builder.vertexCount != 3 {
			t.Errorf("Expected vertex count 3, got %d", builder.vertexCount)
		}
	})

	t.Run("Add edge DTO", func(t *testing.T) {
		builder := &Builder[int, float64, string, bool]{}

		dto := &BasicEdgeDto[int, float64, bool]{
			Origin: 1,
			Target: 2,
			Cost:   10.5,
			Data:   true,
		}

		builder.AddEdgeDto(dto)

		if builder.edgeCount != 1 {
			t.Errorf("Expected edge count 1, got %d", builder.edgeCount)
		}
	})

	t.Run("Add vertex DTO", func(t *testing.T) {
		builder := &Builder[int, float64, string, bool]{}

		dto := &BasicVertexDto[int, string]{
			Id:   1,
			Data: "vertex1",
		}

		builder.AddVertexDto(dto)

		if builder.vertexCount != 1 {
			t.Errorf("Expected vertex count 1, got %d", builder.vertexCount)
		}
	})

	t.Run("Count bidirectional edges", func(t *testing.T) {
		builder := &Builder[int, float64, string, bool]{}

		// Add bidirectional edges
		builder.AddBiEdge(1, 2, 10.5, true)
		builder.AddBiEdge(2, 3, 15.0, false)

		// Add single edge
		builder.AddEdge(1, 3, 20.0, true)

		biEdgeCount := builder.CountBiEdges()
		expected := 3 // (1,2), (2,3), and (1,3) are all counted as bidirectional edges

		if biEdgeCount != expected {
			t.Errorf("Expected bidirectional edge count %d, got %d", expected, biEdgeCount)
		}
	})

	t.Run("Build directed graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, bool]{}

		// Add vertices
		builder.AddVertex(1, "vertex1")
		builder.AddVertex(2, "vertex2")
		builder.AddVertex(3, "vertex3")

		// Add edges
		builder.AddEdge(1, 2, 10.5, true)
		builder.AddEdge(2, 3, 15.0, false)
		builder.AddEdge(1, 3, 20.0, true)

		graph := builder.BuildDirected()

		if graph.GetVertexCount() != 3 {
			t.Errorf("Expected vertex count 3, got %d", graph.GetVertexCount())
		}

		if graph.GetEdgeCount() != 3 {
			t.Errorf("Expected edge count 3, got %d", graph.GetEdgeCount())
		}

		// Test vertex retrieval
		vertex1, err := graph.GetVertexById(1)
		if err != nil {
			t.Errorf("Failed to get vertex 1: %v", err)
		}
		if vertex1.GetId() != 1 {
			t.Errorf("Expected vertex ID 1, got %v", vertex1.GetId())
		}

		// Test vertex data retrieval
		vertexData, err := graph.GetVertexData(vertex1)
		if err != nil {
			t.Errorf("Failed to get vertex data: %v", err)
		}
		if *vertexData != "vertex1" {
			t.Errorf("Expected vertex data 'vertex1', got %v", *vertexData)
		}
	})

	t.Run("Build graph with implicit vertices", func(t *testing.T) {
		builder := &Builder[int, float64, string, bool]{}

		// Add edges only (vertices will be created implicitly)
		builder.AddEdge(1, 2, 10.5, true)
		builder.AddEdge(2, 3, 15.0, false)
		builder.AddEdge(1, 3, 20.0, true)

		graph := builder.BuildDirected()

		if graph.GetVertexCount() != 3 {
			t.Errorf("Expected vertex count 3, got %d", graph.GetVertexCount())
		}

		if graph.GetEdgeCount() != 3 {
			t.Errorf("Expected edge count 3, got %d", graph.GetEdgeCount())
		}

		// Test that all vertices exist
		for _, id := range []int{1, 2, 3} {
			vertex, err := graph.GetVertexById(id)
			if err != nil {
				t.Errorf("Failed to get vertex %d: %v", id, err)
			}
			if vertex.GetId() != id {
				t.Errorf("Expected vertex ID %d, got %v", id, vertex.GetId())
			}
		}
	})

	t.Run("Build graph with string IDs", func(t *testing.T) {
		builder := &Builder[string, int, bool, string]{}

		builder.AddVertex("A", true)
		builder.AddVertex("B", false)
		builder.AddEdge("A", "B", 10, "edge1")

		graph := builder.BuildDirected()

		if graph.GetVertexCount() != 2 {
			t.Errorf("Expected vertex count 2, got %d", graph.GetVertexCount())
		}

		if graph.GetEdgeCount() != 1 {
			t.Errorf("Expected edge count 1, got %d", graph.GetEdgeCount())
		}

		vertexA, err := graph.GetVertexById("A")
		if err != nil {
			t.Errorf("Failed to get vertex A: %v", err)
		}
		if vertexA.GetId() != "A" {
			t.Errorf("Expected vertex ID 'A', got %v", vertexA.GetId())
		}
	})

	t.Run("Large graph construction", func(t *testing.T) {
		builder := &Builder[int, float64, string, bool]{}

		// Add many vertices and edges
		for i := 0; i < 100; i++ {
			builder.AddVertex(i, "vertex")
		}

		for i := 0; i < 99; i++ {
			builder.AddEdge(i, i+1, float64(i), true)
		}

		graph := builder.BuildDirected()

		if graph.GetVertexCount() != 100 {
			t.Errorf("Expected vertex count 100, got %d", graph.GetVertexCount())
		}

		if graph.GetEdgeCount() != 99 {
			t.Errorf("Expected edge count 99, got %d", graph.GetEdgeCount())
		}
	})
}
