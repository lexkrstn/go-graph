package graph

import (
	"testing"
)

func TestNewBellmanFord(t *testing.T) {
	t.Run("Create Bellman-Ford for simple graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 10.0, "edge1-2")
		builder.AddEdge(2, 3, 15.0, "edge2-3")

		graph := builder.BuildDirected()
		bf := NewBellmanFord(graph)

		if bf == nil {
			t.Error("Expected Bellman-Ford instance, got nil")
			return
		}

		if bf.graph != graph {
			t.Error("Expected Bellman-Ford graph to match input graph")
		}

		if bf.vertexData == nil {
			t.Error("Expected vertexData to be initialized")
		}
	})

	t.Run("Create Bellman-Ford for empty graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		graph := builder.BuildDirected()
		bf := NewBellmanFord(graph)

		if bf == nil {
			t.Error("Expected Bellman-Ford instance for empty graph, got nil")
		}
	})
}

func TestBellmanFordFindShortestPath(t *testing.T) {
	t.Run("Simple path between two vertices", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddEdge(1, 2, 10.0, "edge1-2")

		graph := builder.BuildDirected()
		bf := NewBellmanFord(graph)

		path := bf.FindShortestPath(1, 2)
		expectedPath := []int{1, 2}

		if len(path) != len(expectedPath) {
			t.Errorf("Expected path length %d, got %d", len(expectedPath), len(path))
		}

		for i, vertex := range path {
			if vertex != expectedPath[i] {
				t.Errorf("Expected vertex %d at position %d, got %d", expectedPath[i], i, vertex)
			}
		}
	})

	t.Run("Three vertex linear path", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 5.0, "edge1-2")
		builder.AddEdge(2, 3, 10.0, "edge2-3")

		graph := builder.BuildDirected()
		bf := NewBellmanFord(graph)

		path := bf.FindShortestPath(1, 3)
		expectedPath := []int{1, 2, 3}

		if len(path) != len(expectedPath) {
			t.Errorf("Expected path length %d, got %d", len(expectedPath), len(path))
		}

		for i, vertex := range path {
			if vertex != expectedPath[i] {
				t.Errorf("Expected vertex %d at position %d, got %d", expectedPath[i], i, vertex)
			}
		}
	})

	t.Run("Same start and end vertex", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddEdge(1, 2, 10.0, "edge1-2")

		graph := builder.BuildDirected()
		bf := NewBellmanFord(graph)

		path := bf.FindShortestPath(1, 1)
		expectedPath := []int{1}

		if len(path) != len(expectedPath) {
			t.Errorf("Expected path length %d, got %d", len(expectedPath), len(path))
		}

		if path[0] != 1 {
			t.Errorf("Expected single vertex path [1], got %v", path)
		}
	})

	t.Run("Non-existent start vertex", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddEdge(1, 2, 10.0, "edge1-2")

		graph := builder.BuildDirected()
		bf := NewBellmanFord(graph)

		path := bf.FindShortestPath(999, 2)

		if path != nil {
			t.Errorf("Expected nil path for non-existent start vertex, got %v", path)
		}
	})

	t.Run("Non-existent end vertex", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddEdge(1, 2, 10.0, "edge1-2")

		graph := builder.BuildDirected()
		bf := NewBellmanFord(graph)

		path := bf.FindShortestPath(1, 999)

		if path != nil {
			t.Errorf("Expected nil path for non-existent end vertex, got %v", path)
		}
	})

	t.Run("No path between vertices", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 10.0, "edge1-2")
		// No edge from 2 to 3, so no path from 1 to 3

		graph := builder.BuildDirected()
		bf := NewBellmanFord(graph)

		path := bf.FindShortestPath(1, 3)

		if path != nil {
			t.Errorf("Expected nil path for unreachable vertex, got %v", path)
		}
	})

	t.Run("Complex graph with multiple paths", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddVertex(4, "D")

		// Create two paths: 1->2->4 (cost 30) and 1->3->4 (cost 20)
		builder.AddEdge(1, 2, 10.0, "1-2")
		builder.AddEdge(1, 3, 5.0, "1-3")
		builder.AddEdge(2, 4, 20.0, "2-4")
		builder.AddEdge(3, 4, 15.0, "3-4")

		graph := builder.BuildDirected()
		bf := NewBellmanFord(graph)

		path := bf.FindShortestPath(1, 4)

		// Shortest path should be 1->3->4 (cost 20)
		expectedPath := []int{1, 3, 4}

		if len(path) != len(expectedPath) {
			t.Errorf("Expected path length %d, got %d. Path: %v", len(expectedPath), len(path), path)
		}

		for i, vertex := range path {
			if i < len(expectedPath) && vertex != expectedPath[i] {
				t.Errorf("Expected vertex %d at position %d, got %d. Full path: %v", expectedPath[i], i, vertex, path)
			}
		}
	})

	t.Run("Graph with negative edge weights", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddVertex(4, "D")

		// Create paths with negative weights: 1->2->4 (cost 5) and 1->3->4 (cost 8)
		builder.AddEdge(1, 2, 2.0, "1-2")
		builder.AddEdge(1, 3, 5.0, "1-3")
		builder.AddEdge(2, 4, 3.0, "2-4")
		builder.AddEdge(3, 4, 3.0, "3-4")
		builder.AddEdge(2, 3, -1.0, "2-3") // Negative edge

		graph := builder.BuildDirected()
		bf := NewBellmanFord(graph)

		path := bf.FindShortestPath(1, 4)

		// Shortest path should be 1->2->3->4 (cost 2+(-1)+3=4)
		expectedPath := []int{1, 2, 3, 4}

		if len(path) != len(expectedPath) {
			t.Errorf("Expected path length %d, got %d. Path: %v", len(expectedPath), len(path), path)
		}

		for i, vertex := range path {
			if i < len(expectedPath) && vertex != expectedPath[i] {
				t.Errorf("Expected vertex %d at position %d, got %d. Full path: %v", expectedPath[i], i, vertex, path)
			}
		}
	})

	t.Run("Graph with negative cycle", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")

		// Create a negative cycle: 1->2->3->1 with negative total weight
		builder.AddEdge(1, 2, 1.0, "1-2")
		builder.AddEdge(2, 3, 1.0, "2-3")
		builder.AddEdge(3, 1, -3.0, "3-1") // Negative cycle

		graph := builder.BuildDirected()
		bf := NewBellmanFord(graph)

		path := bf.FindShortestPath(1, 3)

		// Should return nil due to negative cycle
		if path != nil {
			t.Errorf("Expected nil path due to negative cycle, got %v", path)
		}
	})

	t.Run("Single vertex graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(42, "Singleton")

		graph := builder.BuildDirected()
		bf := NewBellmanFord(graph)

		path := bf.FindShortestPath(42, 42)
		expectedPath := []int{42}

		if len(path) != len(expectedPath) {
			t.Errorf("Expected path length %d, got %d", len(expectedPath), len(path))
		}

		if path[0] != 42 {
			t.Errorf("Expected path [42], got %v", path)
		}
	})

	t.Run("Zero cost edges", func(t *testing.T) {
		builder := &Builder[int, int, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 0, "free1-2")
		builder.AddEdge(2, 3, 0, "free2-3")

		graph := builder.BuildDirected()
		bf := NewBellmanFord(graph)

		path := bf.FindShortestPath(1, 3)
		expectedPath := []int{1, 2, 3}

		if len(path) != len(expectedPath) {
			t.Errorf("Expected path length %d, got %d", len(expectedPath), len(path))
		}

		for i, vertex := range path {
			if vertex != expectedPath[i] {
				t.Errorf("Expected vertex %d at position %d, got %d", expectedPath[i], i, vertex)
			}
		}
	})
}

func TestBellmanFordWithDifferentTypes(t *testing.T) {
	t.Run("String IDs with float costs", func(t *testing.T) {
		builder := &Builder[string, float32, int, bool]{}
		builder.AddVertex("start", 1)
		builder.AddVertex("middle", 2)
		builder.AddVertex("end", 3)
		builder.AddEdge("start", "middle", 1.5, true)
		builder.AddEdge("middle", "end", 2.5, false)

		graph := builder.BuildDirected()
		bf := NewBellmanFord(graph)

		path := bf.FindShortestPath("start", "end")
		expectedPath := []string{"start", "middle", "end"}

		if len(path) != len(expectedPath) {
			t.Errorf("Expected path length %d, got %d", len(expectedPath), len(path))
		}

		for i, vertex := range path {
			if vertex != expectedPath[i] {
				t.Errorf("Expected vertex %s at position %d, got %s", expectedPath[i], i, vertex)
			}
		}
	})

	t.Run("Integer IDs with integer costs", func(t *testing.T) {
		builder := &Builder[uint32, uint16, string, string]{}
		builder.AddVertex(1, "First")
		builder.AddVertex(2, "Second")
		builder.AddVertex(3, "Third")
		builder.AddEdge(1, 2, 100, "edge1")
		builder.AddEdge(2, 3, 200, "edge2")
		builder.AddEdge(1, 3, 250, "direct")

		graph := builder.BuildDirected()
		bf := NewBellmanFord(graph)

		path := bf.FindShortestPath(1, 3)
		// The shortest path should be the direct route 1->3 (cost 250)
		// rather than 1->2->3 (cost 300)
		expectedPath := []uint32{1, 3}

		if len(path) != len(expectedPath) {
			t.Errorf("Expected path length %d, got %d", len(expectedPath), len(path))
		}

		for i, vertex := range path {
			if vertex != expectedPath[i] {
				t.Errorf("Expected vertex %d at position %d, got %d", expectedPath[i], i, vertex)
			}
		}
	})
}

func TestBellmanFordWithAmplifier(t *testing.T) {
	t.Run("Amplifier disables specific edges", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddVertex(4, "D")

		// Create two paths: 1->2->4 (cost 3) and 1->3->4 (cost 5)
		builder.AddEdge(1, 2, 1.0, "1-2")
		builder.AddEdge(1, 3, 2.0, "1-3")
		builder.AddEdge(2, 4, 2.0, "2-4")
		builder.AddEdge(3, 4, 3.0, "3-4")

		graph := builder.BuildDirected()
		bf := NewBellmanFord(graph)

		// Without amplifier: should take 1->2->4 (cost 3)
		pathWithoutAmplifier := bf.FindShortestPath(1, 4)
		expectedWithout := []int{1, 2, 4}
		if !slicesEqual(pathWithoutAmplifier, expectedWithout) {
			t.Errorf("Without amplifier: expected %v, got %v", expectedWithout, pathWithoutAmplifier)
		}

		// With amplifier that disables edge 1->2: should take 1->3->4 (cost 5)
		bf.Amplifier = func(origin *Vertex[int, float64], edge *Edge[int, float64]) (float64, bool) {
			if origin.id == 1 && edge.targetVertex.id == 2 {
				return 0.0, false // Disable edge 1->2
			}
			return edge.cost, true // Keep other edges unchanged
		}

		pathWithAmplifier := bf.FindShortestPath(1, 4)
		expectedWith := []int{1, 3, 4}
		if !slicesEqual(pathWithAmplifier, expectedWith) {
			t.Errorf("With amplifier: expected %v, got %v", expectedWith, pathWithAmplifier)
		}
	})

	t.Run("Amplifier modifies edge costs", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddVertex(4, "D")

		// Create two paths: 1->2->4 (cost 3) and 1->3->4 (cost 5)
		builder.AddEdge(1, 2, 1.0, "1-2")
		builder.AddEdge(1, 3, 2.0, "1-3")
		builder.AddEdge(2, 4, 2.0, "2-4")
		builder.AddEdge(3, 4, 3.0, "3-4")

		graph := builder.BuildDirected()
		bf := NewBellmanFord(graph)

		// Amplifier that triples the cost of edge 1->2
		bf.Amplifier = func(origin *Vertex[int, float64], edge *Edge[int, float64]) (float64, bool) {
			if origin.id == 1 && edge.targetVertex.id == 2 {
				return edge.cost * 3.0, true // Triple the cost
			}
			return edge.cost, true // Keep other edges unchanged
		}

		path := bf.FindShortestPath(1, 4)
		// Now 1->2->4 costs 5 (1*3 + 2) and 1->3->4 costs 5 (2 + 3)
		// Both paths have equal cost, so either is acceptable
		expected1 := []int{1, 2, 4}
		expected2 := []int{1, 3, 4}
		if !slicesEqual(path, expected1) && !slicesEqual(path, expected2) {
			t.Errorf("Expected either %v or %v, got %v", expected1, expected2, path)
		}
	})

	t.Run("Amplifier disables all edges - no path", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 1.0, "1-2")
		builder.AddEdge(2, 3, 1.0, "2-3")

		graph := builder.BuildDirected()
		bf := NewBellmanFord(graph)

		// Amplifier that disables all edges
		bf.Amplifier = func(origin *Vertex[int, float64], edge *Edge[int, float64]) (float64, bool) {
			return 0.0, false // Disable all edges
		}

		path := bf.FindShortestPath(1, 3)
		if path != nil {
			t.Errorf("Expected nil path when all edges are disabled, got %v", path)
		}
	})

	t.Run("Amplifier with zero cost edges", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 5.0, "1-2")
		builder.AddEdge(1, 3, 10.0, "1-3")
		builder.AddEdge(2, 3, 5.0, "2-3")

		graph := builder.BuildDirected()
		bf := NewBellmanFord(graph)

		// Amplifier that makes edge 1->2 free
		bf.Amplifier = func(origin *Vertex[int, float64], edge *Edge[int, float64]) (float64, bool) {
			if origin.id == 1 && edge.targetVertex.id == 2 {
				return 0.0, true // Make edge 1->2 free
			}
			return edge.cost, true
		}

		path := bf.FindShortestPath(1, 3)
		// Should take 1->2->3 (cost 0+5=5) instead of 1->3 (cost 10)
		expected := []int{1, 2, 3}
		if !slicesEqual(path, expected) {
			t.Errorf("Expected %v, got %v", expected, path)
		}
	})

	t.Run("Amplifier with very high costs", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 1.0, "1-2")
		builder.AddEdge(1, 3, 2.0, "1-3")
		builder.AddEdge(2, 3, 1.0, "2-3")

		graph := builder.BuildDirected()
		bf := NewBellmanFord(graph)

		// Amplifier that makes edge 1->2 very expensive
		bf.Amplifier = func(origin *Vertex[int, float64], edge *Edge[int, float64]) (float64, bool) {
			if origin.id == 1 && edge.targetVertex.id == 2 {
				return 1000.0, true // Make edge 1->2 very expensive
			}
			return edge.cost, true
		}

		path := bf.FindShortestPath(1, 3)
		// Should take direct path 1->3 (cost 2) instead of 1->2->3 (cost 1000+1=1001)
		expected := []int{1, 3}
		if !slicesEqual(path, expected) {
			t.Errorf("Expected %v, got %v", expected, path)
		}
	})

	t.Run("Amplifier with nil check", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddEdge(1, 2, 1.0, "1-2")

		graph := builder.BuildDirected()
		bf := NewBellmanFord(graph)

		// Ensure nil amplifier doesn't cause issues
		bf.Amplifier = nil

		path := bf.FindShortestPath(1, 2)
		expected := []int{1, 2}
		if !slicesEqual(path, expected) {
			t.Errorf("Expected %v, got %v", expected, path)
		}
	})
}

func TestBellmanFordHasNegativeCycle(t *testing.T) {
	t.Run("No negative cycle", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 1.0, "1-2")
		builder.AddEdge(2, 3, 1.0, "2-3")
		builder.AddEdge(3, 1, 1.0, "3-1")

		graph := builder.BuildDirected()
		bf := NewBellmanFord(graph)

		hasCycle := bf.HasNegativeCycle(1)
		if hasCycle {
			t.Error("Expected no negative cycle, but one was detected")
		}
	})

	t.Run("Negative cycle present", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 1.0, "1-2")
		builder.AddEdge(2, 3, 1.0, "2-3")
		builder.AddEdge(3, 1, -4.0, "3-1") // Negative cycle

		graph := builder.BuildDirected()
		bf := NewBellmanFord(graph)

		hasCycle := bf.HasNegativeCycle(1)
		if !hasCycle {
			t.Error("Expected negative cycle, but none was detected")
		}
	})

	t.Run("Negative cycle not reachable from start", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddVertex(4, "D")
		builder.AddEdge(1, 2, 1.0, "1-2")
		builder.AddEdge(3, 4, 1.0, "3-4")
		builder.AddEdge(4, 3, -2.0, "4-3") // Negative cycle in disconnected component

		graph := builder.BuildDirected()
		bf := NewBellmanFord(graph)

		hasCycle := bf.HasNegativeCycle(1)
		if hasCycle {
			t.Error("Expected no negative cycle reachable from start, but one was detected")
		}
	})

	t.Run("Non-existent start vertex", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddEdge(1, 2, 1.0, "1-2")

		graph := builder.BuildDirected()
		bf := NewBellmanFord(graph)

		hasCycle := bf.HasNegativeCycle(999)
		if hasCycle {
			t.Error("Expected no negative cycle for non-existent start vertex, but one was detected")
		}
	})
}
