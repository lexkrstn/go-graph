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
func zeroHeuristic[I Id, C Cost, V any, E any](current *Vertex[I, C], target *Vertex[I, C]) C {
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
		heuristic := func(current *Vertex[int, float64], target *Vertex[int, float64]) float64 {
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
		heuristic := func(current *Vertex[int, float64], target *Vertex[int, float64]) float64 {
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
		heuristic := zeroHeuristic[int, float64, string, string]
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
		heuristic := zeroHeuristic[int, float64, string, string]
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
		heuristic := zeroHeuristic[int, float64, string, string]
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
		heuristic := zeroHeuristic[int, float64, string, string]
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
		heuristic := zeroHeuristic[int, float64, string, string]
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
		heuristic := zeroHeuristic[int, float64, string, string]
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
		heuristic := func(current *Vertex[int, float64], goal *Vertex[int, float64]) float64 {
			// Convert vertex ID to 2D coordinates (1-based)
			currentX := (current.id-1)%3 + 1
			currentY := (current.id-1)/3 + 1
			goalX := (goal.id-1)%3 + 1
			goalY := (goal.id-1)/3 + 1
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
		heuristic := func(current *Vertex[int, float64], goal *Vertex[int, float64]) float64 {
			return euclideanDistance(current.id, 0, goal.id, 0)
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
		zeroHeur := zeroHeuristic[int, float64, string, string]
		astarZero := NewAStar(graph, zeroHeur)

		// Good heuristic (should guide towards goal)
		goodHeur := func(current *Vertex[int, float64], goal *Vertex[int, float64]) float64 {
			return float64(goal.id - current.id) // Simple linear heuristic
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

func TestAStarWithAmplifier(t *testing.T) {
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
		heuristic := zeroHeuristic[int, float64, string, string]
		astar := NewAStar(graph, heuristic)

		// Without amplifier: should take 1->2->4 (cost 3)
		pathWithoutAmplifier := astar.FindShortestPath(1, 4)
		expectedWithout := []int{1, 2, 4}
		if !slicesEqualAStar(pathWithoutAmplifier, expectedWithout) {
			t.Errorf("Without amplifier: expected %v, got %v", expectedWithout, pathWithoutAmplifier)
		}

		// With amplifier that disables edge 1->2: should take 1->3->4 (cost 5)
		astar.Amplifier = func(origin *Vertex[int, float64], edge *Edge[int, float64]) (float64, bool) {
			if origin.id == 1 && edge.targetVertex.id == 2 {
				return 0.0, false // Disable edge 1->2
			}
			return edge.cost, true // Keep other edges unchanged
		}

		pathWithAmplifier := astar.FindShortestPath(1, 4)
		expectedWith := []int{1, 3, 4}
		if !slicesEqualAStar(pathWithAmplifier, expectedWith) {
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
		heuristic := zeroHeuristic[int, float64, string, string]
		astar := NewAStar(graph, heuristic)

		// Amplifier that triples the cost of edge 1->2
		astar.Amplifier = func(origin *Vertex[int, float64], edge *Edge[int, float64]) (float64, bool) {
			if origin.id == 1 && edge.targetVertex.id == 2 {
				return edge.cost * 3.0, true // Triple the cost
			}
			return edge.cost, true // Keep other edges unchanged
		}

		path := astar.FindShortestPath(1, 4)
		// Now 1->2->4 costs 7 (1*3 + 2), so 1->3->4 (cost 5) should be chosen
		expected := []int{1, 3, 4}
		if !slicesEqualAStar(path, expected) {
			t.Errorf("Expected %v, got %v", expected, path)
		}
	})

	t.Run("Amplifier with heuristic interaction", func(t *testing.T) {
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

		// Heuristic that guides towards vertex 3
		heuristic := func(current *Vertex[int, float64], goal *Vertex[int, float64]) float64 {
			if current.id == 1 {
				return 1.0 // Favor going to 3 from 1
			}
			return 0.0
		}
		astar := NewAStar(graph, heuristic)

		// Without amplifier: should take 1->2->4 (cost 3) despite heuristic
		pathWithoutAmplifier := astar.FindShortestPath(1, 4)
		expectedWithout := []int{1, 2, 4}
		if !slicesEqualAStar(pathWithoutAmplifier, expectedWithout) {
			t.Errorf("Without amplifier: expected %v, got %v", expectedWithout, pathWithoutAmplifier)
		}

		// With amplifier that disables edge 1->2: should take 1->3->4 (cost 5)
		astar.Amplifier = func(origin *Vertex[int, float64], edge *Edge[int, float64]) (float64, bool) {
			if origin.id == 1 && edge.targetVertex.id == 2 {
				return 0.0, false // Disable edge 1->2
			}
			return edge.cost, true
		}

		pathWithAmplifier := astar.FindShortestPath(1, 4)
		expectedWith := []int{1, 3, 4}
		if !slicesEqualAStar(pathWithAmplifier, expectedWith) {
			t.Errorf("With amplifier: expected %v, got %v", expectedWith, pathWithAmplifier)
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
		heuristic := zeroHeuristic[int, float64, string, string]
		astar := NewAStar(graph, heuristic)

		// Amplifier that disables all edges
		astar.Amplifier = func(origin *Vertex[int, float64], edge *Edge[int, float64]) (float64, bool) {
			return 0.0, false // Disable all edges
		}

		path := astar.FindShortestPath(1, 3)
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
		heuristic := zeroHeuristic[int, float64, string, string]
		astar := NewAStar(graph, heuristic)

		// Amplifier that makes edge 1->2 free
		astar.Amplifier = func(origin *Vertex[int, float64], edge *Edge[int, float64]) (float64, bool) {
			if origin.id == 1 && edge.targetVertex.id == 2 {
				return 0.0, true // Make edge 1->2 free
			}
			return edge.cost, true
		}

		path := astar.FindShortestPath(1, 3)
		// Should take 1->2->3 (cost 0+5=5) instead of 1->3 (cost 10)
		expected := []int{1, 2, 3}
		if !slicesEqualAStar(path, expected) {
			t.Errorf("Expected %v, got %v", expected, path)
		}
	})

	t.Run("Amplifier with Manhattan distance heuristic", func(t *testing.T) {
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
		heuristic := func(current *Vertex[int, float64], goal *Vertex[int, float64]) float64 {
			currentX := (current.id-1)%3 + 1
			currentY := (current.id-1)/3 + 1
			goalX := (goal.id-1)%3 + 1
			goalY := (goal.id-1)/3 + 1
			return manhattanDistance(currentX, currentY, goalX, goalY)
		}

		astar := NewAStar(graph, heuristic)

		// Test path from corner to corner (1 to 9) without amplifier
		pathWithoutAmplifier := astar.FindShortestPath(1, 9)
		if len(pathWithoutAmplifier) == 0 {
			t.Error("Expected a path from 1 to 9 without amplifier")
		}

		// With amplifier that disables horizontal edges in the middle row
		astar.Amplifier = func(origin *Vertex[int, float64], edge *Edge[int, float64]) (float64, bool) {
			// Disable horizontal edges from vertices 4, 5, 6
			if (origin.id == 4 || origin.id == 5 || origin.id == 6) &&
				(edge.targetVertex.id == origin.id+1 || edge.targetVertex.id == origin.id-1) {
				return 0.0, false
			}
			return edge.cost, true
		}

		pathWithAmplifier := astar.FindShortestPath(1, 9)
		// Should still find a path, but it will be different due to disabled edges
		if pathWithAmplifier == nil {
			t.Error("Expected a path from 1 to 9 with amplifier")
		}

		// Verify the path doesn't use disabled edges
		for i := 0; i < len(pathWithAmplifier)-1; i++ {
			current := pathWithAmplifier[i]
			next := pathWithAmplifier[i+1]

			// Check if this would be a disabled horizontal edge
			if (current == 4 || current == 5 || current == 6) &&
				(next == current+1 || next == current-1) {
				t.Errorf("Path uses disabled horizontal edge: %d->%d", current, next)
			}
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
		heuristic := zeroHeuristic[int, float64, string, string]
		astar := NewAStar(graph, heuristic)

		// Amplifier that makes edge 1->2 very expensive
		astar.Amplifier = func(origin *Vertex[int, float64], edge *Edge[int, float64]) (float64, bool) {
			if origin.id == 1 && edge.targetVertex.id == 2 {
				return 1000.0, true // Make edge 1->2 very expensive
			}
			return edge.cost, true
		}

		path := astar.FindShortestPath(1, 3)
		// Should take direct path 1->3 (cost 2) instead of 1->2->3 (cost 1000+1=1001)
		expected := []int{1, 3}
		if !slicesEqualAStar(path, expected) {
			t.Errorf("Expected %v, got %v", expected, path)
		}
	})

	t.Run("Amplifier with nil check", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddEdge(1, 2, 1.0, "1-2")

		graph := builder.BuildDirected()
		heuristic := zeroHeuristic[int, float64, string, string]
		astar := NewAStar(graph, heuristic)

		// Ensure nil amplifier doesn't cause issues
		astar.Amplifier = nil

		path := astar.FindShortestPath(1, 2)
		expected := []int{1, 2}
		if !slicesEqualAStar(path, expected) {
			t.Errorf("Expected %v, got %v", expected, path)
		}
	})

	t.Run("Amplifier with complex graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddVertex(4, "D")
		builder.AddVertex(5, "E")

		// Create a complex graph with multiple paths
		builder.AddEdge(1, 2, 1.0, "1-2")
		builder.AddEdge(1, 3, 2.0, "1-3")
		builder.AddEdge(2, 4, 1.0, "2-4")
		builder.AddEdge(3, 4, 1.0, "3-4")
		builder.AddEdge(2, 5, 2.0, "2-5")
		builder.AddEdge(4, 5, 1.0, "4-5")

		graph := builder.BuildDirected()
		heuristic := zeroHeuristic[int, float64, string, string]
		astar := NewAStar(graph, heuristic)

		// Amplifier that disables edges from vertex 2
		astar.Amplifier = func(origin *Vertex[int, float64], edge *Edge[int, float64]) (float64, bool) {
			if origin.id == 2 {
				return 0.0, false // Disable all edges from vertex 2
			}
			return edge.cost, true
		}

		path := astar.FindShortestPath(1, 5)
		// Should take 1->3->4->5 (cost 2+1+1=4) since 1->2->5 and 1->2->4->5 are blocked
		expected := []int{1, 3, 4, 5}
		if !slicesEqualAStar(path, expected) {
			t.Errorf("Expected %v, got %v", expected, path)
		}
	})

	t.Run("Amplifier with different cost types", func(t *testing.T) {
		builder := &Builder[string, int, string, string]{}
		builder.AddVertex("start", "Start")
		builder.AddVertex("middle", "Middle")
		builder.AddVertex("end", "End")
		builder.AddEdge("start", "middle", 10, "start-middle")
		builder.AddEdge("start", "end", 20, "start-end")
		builder.AddEdge("middle", "end", 5, "middle-end")

		graph := builder.BuildDirected()
		heuristic := func(current *Vertex[string, int], goal *Vertex[string, int]) int {
			return 0 // Zero heuristic
		}
		astar := NewAStar(graph, heuristic)

		// Amplifier that halves the cost of start->middle
		astar.Amplifier = func(origin *Vertex[string, int], edge *Edge[string, int]) (int, bool) {
			if origin.id == "start" && edge.targetVertex.id == "middle" {
				return edge.cost / 2, true // Halve the cost
			}
			return edge.cost, true
		}

		path := astar.FindShortestPath("start", "end")
		// Should take start->middle->end (cost 5+5=10) instead of start->end (cost 20)
		expected := []string{"start", "middle", "end"}
		if !slicesEqualString(path, expected) {
			t.Errorf("Expected %v, got %v", expected, path)
		}
	})
}

// Helper functions to compare slices
func slicesEqualAStar(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func slicesEqualString(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
