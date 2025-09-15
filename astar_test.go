package graph

import (
	"math"
	"testing"
)

// Simple Manhattan distance heuristic for testing
func manhattanDistance(x1, y1, x2, y2 int) float64 {
	return math.Abs(float64(x1-x2)) + math.Abs(float64(y1-y2))
}

// Zero heuristic (makes A* behave like Dijkstra)
func zeroHeuristic[I Id, C Cost](current I, goal I) C {
	var zero C
	return zero
}

// Euclidean distance heuristic for 2D coordinates
func euclideanDistance(x1, y1, x2, y2 int) float64 {
	dx := float64(x1 - x2)
	dy := float64(y1 - y2)
	return math.Sqrt(dx*dx + dy*dy)
}

func TestNewAStar(t *testing.T) {
	t.Run("Create A* for simple graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 10.0, "edge1-2")
		builder.AddEdge(2, 3, 15.0, "edge2-3")

		graph := builder.BuildDirected()
		heuristic := func(current int, goal int) float64 {
			return 0.0 // Zero heuristic
		}
		astar := NewAStar(graph, heuristic)

		if astar == nil {
			t.Error("Expected A* instance, got nil")
			return
		}

		if astar.graph != graph {
			t.Error("Expected A* graph to match input graph")
		}

		if astar.heap == nil {
			t.Error("Expected heap to be initialized")
		}

		if astar.heuristic == nil {
			t.Error("Expected heuristic function to be set")
		}
	})

	t.Run("Create A* for empty graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		graph := builder.BuildDirected()
		heuristic := func(current int, goal int) float64 {
			return 0.0
		}
		astar := NewAStar(graph, heuristic)

		if astar == nil {
			t.Error("Expected A* instance for empty graph, got nil")
		}
	})
}

func TestAStarFindShortestPath(t *testing.T) {
	t.Run("Simple path between two vertices", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddEdge(1, 2, 10.0, "edge1-2")

		graph := builder.BuildDirected()
		heuristic := zeroHeuristic[int, float64]
		astar := NewAStar(graph, heuristic)

		path := astar.FindShortestPath(1, 2)
		expected := []int{1, 2}

		if len(path) != len(expected) {
			t.Errorf("Expected path length %d, got %d", len(expected), len(path))
			return
		}

		for i, v := range expected {
			if path[i] != v {
				t.Errorf("Expected path[%d] = %d, got %d", i, v, path[i])
			}
		}
	})

	t.Run("Path to same vertex", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddEdge(1, 2, 10.0, "edge1-2")

		graph := builder.BuildDirected()
		heuristic := zeroHeuristic[int, float64]
		astar := NewAStar(graph, heuristic)

		path := astar.FindShortestPath(1, 1)
		expected := []int{1}

		if len(path) != len(expected) {
			t.Errorf("Expected path length %d, got %d", len(expected), len(path))
			return
		}

		if path[0] != expected[0] {
			t.Errorf("Expected path[0] = %d, got %d", expected[0], path[0])
		}
	})

	t.Run("No path exists", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 10.0, "edge1-2")
		// No edge from 2 to 3

		graph := builder.BuildDirected()
		heuristic := zeroHeuristic[int, float64]
		astar := NewAStar(graph, heuristic)

		path := astar.FindShortestPath(1, 3)
		if path != nil {
			t.Errorf("Expected nil path, got %v", path)
		}
	})

	t.Run("Non-existent start vertex", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddEdge(1, 2, 10.0, "edge1-2")

		graph := builder.BuildDirected()
		heuristic := zeroHeuristic[int, float64]
		astar := NewAStar(graph, heuristic)

		path := astar.FindShortestPath(99, 2)
		if path != nil {
			t.Errorf("Expected nil path for non-existent start, got %v", path)
		}
	})

	t.Run("Non-existent end vertex", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddEdge(1, 2, 10.0, "edge1-2")

		graph := builder.BuildDirected()
		heuristic := zeroHeuristic[int, float64]
		astar := NewAStar(graph, heuristic)

		path := astar.FindShortestPath(1, 99)
		if path != nil {
			t.Errorf("Expected nil path for non-existent end, got %v", path)
		}
	})

	t.Run("Complex graph with multiple paths", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddVertex(4, "D")
		builder.AddEdge(1, 2, 1.0, "edge1-2")
		builder.AddEdge(1, 3, 4.0, "edge1-3")
		builder.AddEdge(2, 4, 2.0, "edge2-4")
		builder.AddEdge(3, 4, 1.0, "edge3-4")

		graph := builder.BuildDirected()
		heuristic := zeroHeuristic[int, float64]
		astar := NewAStar(graph, heuristic)

		path := astar.FindShortestPath(1, 4)
		expected := []int{1, 2, 4} // Should find the shorter path: 1->2->4 (cost 3) vs 1->3->4 (cost 5)

		if len(path) != len(expected) {
			t.Errorf("Expected path length %d, got %d", len(expected), len(path))
			return
		}

		for i, v := range expected {
			if path[i] != v {
				t.Errorf("Expected path[%d] = %d, got %d", i, v, path[i])
			}
		}
	})

	t.Run("A* with Manhattan distance heuristic", func(t *testing.T) {
		// Create a grid-like graph where vertices represent 2D coordinates
		builder := &Builder[int, float64, string, string]{}

		// Add vertices representing a 3x3 grid
		// 1 2 3
		// 4 5 6
		// 7 8 9
		for i := 1; i <= 9; i++ {
			builder.AddVertex(i, "vertex")
		}

		// Add horizontal edges
		builder.AddEdge(1, 2, 1.0, "h1-2")
		builder.AddEdge(2, 3, 1.0, "h2-3")
		builder.AddEdge(4, 5, 1.0, "h4-5")
		builder.AddEdge(5, 6, 1.0, "h5-6")
		builder.AddEdge(7, 8, 1.0, "h7-8")
		builder.AddEdge(8, 9, 1.0, "h8-9")

		// Add vertical edges
		builder.AddEdge(1, 4, 1.0, "v1-4")
		builder.AddEdge(4, 7, 1.0, "v4-7")
		builder.AddEdge(2, 5, 1.0, "v2-5")
		builder.AddEdge(5, 8, 1.0, "v5-8")
		builder.AddEdge(3, 6, 1.0, "v3-6")
		builder.AddEdge(6, 9, 1.0, "v6-9")

		graph := builder.BuildDirected()

		// Manhattan distance heuristic
		heuristic := func(current int, goal int) float64 {
			// Convert vertex ID to 2D coordinates (1-based)
			currentX := (current-1)%3 + 1
			currentY := (current-1)/3 + 1
			goalX := (goal-1)%3 + 1
			goalY := (goal-1)/3 + 1
			return manhattanDistance(currentX, currentY, goalX, goalY)
		}

		astar := NewAStar(graph, heuristic)

		// Test path from corner to corner (1 to 9)
		path := astar.FindShortestPath(1, 9)
		expected := []int{1, 2, 3, 6, 9} // One possible optimal path

		if len(path) != len(expected) {
			t.Errorf("Expected path length %d, got %d", len(expected), len(path))
			return
		}

		// Verify it's a valid path
		for i := 0; i < len(path)-1; i++ {
			current := path[i]
			next := path[i+1]

			// Check if there's an edge from current to next
			vertex, _ := graph.GetVertexById(current)
			found := false
			for _, edge := range vertex.GetEdges() {
				if edge.GetTargetVertex().GetId() == next {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("No edge found from %d to %d in path", current, next)
			}
		}
	})

	t.Run("A* with Euclidean distance heuristic", func(t *testing.T) {
		// Create a simple line graph
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddVertex(4, "D")
		builder.AddEdge(1, 2, 1.0, "edge1-2")
		builder.AddEdge(2, 3, 1.0, "edge2-3")
		builder.AddEdge(3, 4, 1.0, "edge3-4")

		graph := builder.BuildDirected()

		// Euclidean distance heuristic
		heuristic := func(current int, goal int) float64 {
			return euclideanDistance(current, 0, goal, 0)
		}

		astar := NewAStar(graph, heuristic)

		path := astar.FindShortestPath(1, 4)
		expected := []int{1, 2, 3, 4}

		if len(path) != len(expected) {
			t.Errorf("Expected path length %d, got %d", len(expected), len(path))
			return
		}

		for i, v := range expected {
			if path[i] != v {
				t.Errorf("Expected path[%d] = %d, got %d", i, v, path[i])
			}
		}
	})

	t.Run("A* vs Dijkstra comparison", func(t *testing.T) {
		// Create a graph where A* should be more efficient than Dijkstra
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "start")
		builder.AddVertex(2, "middle1")
		builder.AddVertex(3, "middle2")
		builder.AddVertex(4, "goal")

		// Direct path (expensive)
		builder.AddEdge(1, 4, 100.0, "direct")

		// Indirect path (cheaper)
		builder.AddEdge(1, 2, 1.0, "start-middle1")
		builder.AddEdge(2, 3, 1.0, "middle1-middle2")
		builder.AddEdge(3, 4, 1.0, "middle2-goal")

		graph := builder.BuildDirected()

		// Zero heuristic (should behave like Dijkstra)
		zeroHeur := zeroHeuristic[int, float64]
		astarZero := NewAStar(graph, zeroHeur)

		// Good heuristic (should guide towards goal)
		goodHeur := func(current int, goal int) float64 {
			return float64(goal - current) // Simple linear heuristic
		}
		astarGood := NewAStar(graph, goodHeur)

		// Dijkstra for comparison
		dijkstra := NewDijkstra(graph)

		pathAStarZero := astarZero.FindShortestPath(1, 4)
		pathAStarGood := astarGood.FindShortestPath(1, 4)
		pathDijkstra := dijkstra.FindShortestPath(1, 4)

		// All should find the same optimal path
		expected := []int{1, 2, 3, 4}

		// Check A* with zero heuristic
		if len(pathAStarZero) != len(expected) {
			t.Errorf("A* zero heuristic: Expected path length %d, got %d", len(expected), len(pathAStarZero))
		}

		// Check A* with good heuristic
		if len(pathAStarGood) != len(expected) {
			t.Errorf("A* good heuristic: Expected path length %d, got %d", len(expected), len(pathAStarGood))
		}

		// Check Dijkstra
		if len(pathDijkstra) != len(expected) {
			t.Errorf("Dijkstra: Expected path length %d, got %d", len(expected), len(pathDijkstra))
		}

		// All paths should be identical
		for i, v := range expected {
			if pathAStarZero[i] != v {
				t.Errorf("A* zero heuristic: Expected path[%d] = %d, got %d", i, v, pathAStarZero[i])
			}
			if pathAStarGood[i] != v {
				t.Errorf("A* good heuristic: Expected path[%d] = %d, got %d", i, v, pathAStarGood[i])
			}
			if pathDijkstra[i] != v {
				t.Errorf("Dijkstra: Expected path[%d] = %d, got %d", i, v, pathDijkstra[i])
			}
		}
	})
}
