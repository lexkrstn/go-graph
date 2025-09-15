package graph

import (
	"fmt"
	"testing"
)

func TestNewDFS(t *testing.T) {
	t.Run("Create DFS for simple graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 10.0, "edge1-2")
		builder.AddEdge(2, 3, 15.0, "edge2-3")

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		if dfs == nil {
			t.Error("Expected DFS instance, got nil")
			return
		}

		if dfs.graph != graph {
			t.Error("Expected DFS graph to match input graph")
		}

		if dfs.vertexData == nil {
			t.Error("Expected vertexData to be initialized")
		}
	})

	t.Run("Create DFS for empty graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		if dfs == nil {
			t.Error("Expected DFS instance for empty graph, got nil")
		}
	})
}

func TestDFSGetAllReachable(t *testing.T) {
	t.Run("Get all reachable from single vertex", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		result := dfs.GetAllReachable(1)

		if len(result) != 1 {
			t.Errorf("Expected 1 vertex, got %d", len(result))
		}

		if result[0] != 1 {
			t.Errorf("Expected vertex 1, got %d", result[0])
		}
	})

	t.Run("Get all reachable from vertex in connected graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddVertex(4, "D")
		builder.AddEdge(1, 2, 10.0, "edge1-2")
		builder.AddEdge(2, 3, 15.0, "edge2-3")
		builder.AddEdge(1, 4, 5.0, "edge1-4")

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		result := dfs.GetAllReachable(1)

		if len(result) != 4 {
			t.Errorf("Expected 4 vertices, got %d", len(result))
		}

		// Check that all vertices are present
		expected := map[int]bool{1: true, 2: true, 3: true, 4: true}
		for _, vertex := range result {
			if !expected[vertex] {
				t.Errorf("Unexpected vertex %d in result", vertex)
			}
		}
	})

	t.Run("Get all reachable from vertex in disconnected graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddVertex(4, "D")
		builder.AddEdge(1, 2, 10.0, "edge1-2")
		builder.AddEdge(3, 4, 15.0, "edge3-4")

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		result := dfs.GetAllReachable(1)

		if len(result) != 2 {
			t.Errorf("Expected 2 vertices, got %d", len(result))
		}

		// Check that only vertices 1 and 2 are present
		expected := map[int]bool{1: true, 2: true}
		for _, vertex := range result {
			if !expected[vertex] {
				t.Errorf("Unexpected vertex %d in result", vertex)
			}
		}
	})

	t.Run("Get all reachable from non-existent vertex", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		result := dfs.GetAllReachable(999)

		if result != nil {
			t.Error("Expected nil result for non-existent vertex")
		}
	})
}

func TestDFSTraverseFrom(t *testing.T) {
	t.Run("Traverse with callback on single vertex", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		var visitedVertices []int
		var visitedEdges []*Edge[int, float64]

		dfs.TraverseFrom(1, func(vertex *Vertex[int, float64], edge *Edge[int, float64]) {
			visitedVertices = append(visitedVertices, vertex.GetId())
			visitedEdges = append(visitedEdges, edge)
		})

		if len(visitedVertices) != 1 {
			t.Errorf("Expected 1 vertex, got %d", len(visitedVertices))
		}

		if visitedVertices[0] != 1 {
			t.Errorf("Expected vertex 1, got %d", visitedVertices[0])
		}

		if len(visitedEdges) != 1 {
			t.Errorf("Expected 1 edge (nil for start), got %d", len(visitedEdges))
		}

		if visitedEdges[0] != nil {
			t.Error("Expected nil edge for start vertex")
		}
	})

	t.Run("Traverse with callback on connected graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddVertex(4, "D")
		builder.AddEdge(1, 2, 10.0, "edge1-2")
		builder.AddEdge(2, 3, 15.0, "edge2-3")
		builder.AddEdge(1, 4, 5.0, "edge1-4")

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		var visitedVertices []int
		var visitedEdges []*Edge[int, float64]

		dfs.TraverseFrom(1, func(vertex *Vertex[int, float64], edge *Edge[int, float64]) {
			visitedVertices = append(visitedVertices, vertex.GetId())
			visitedEdges = append(visitedEdges, edge)
		})

		if len(visitedVertices) != 4 {
			t.Errorf("Expected 4 vertices, got %d", len(visitedVertices))
		}

		// Check that all vertices are present
		expected := map[int]bool{1: true, 2: true, 3: true, 4: true}
		for _, vertex := range visitedVertices {
			if !expected[vertex] {
				t.Errorf("Unexpected vertex %d in result", vertex)
			}
		}

		// First vertex should have nil edge (start vertex)
		if visitedEdges[0] != nil {
			t.Error("Expected nil edge for start vertex")
		}

		// Other vertices should have non-nil edges
		for i := 1; i < len(visitedEdges); i++ {
			if visitedEdges[i] == nil {
				t.Errorf("Expected non-nil edge for vertex %d", visitedVertices[i])
			}
		}
	})

	t.Run("Traverse with callback on disconnected graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddVertex(4, "D")
		builder.AddEdge(1, 2, 10.0, "edge1-2")
		builder.AddEdge(3, 4, 15.0, "edge3-4")

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		var visitedVertices []int
		var visitedEdges []*Edge[int, float64]

		dfs.TraverseFrom(1, func(vertex *Vertex[int, float64], edge *Edge[int, float64]) {
			visitedVertices = append(visitedVertices, vertex.GetId())
			visitedEdges = append(visitedEdges, edge)
		})

		if len(visitedVertices) != 2 {
			t.Errorf("Expected 2 vertices, got %d", len(visitedVertices))
		}

		// Check that only vertices 1 and 2 are present
		expected := map[int]bool{1: true, 2: true}
		for _, vertex := range visitedVertices {
			if !expected[vertex] {
				t.Errorf("Unexpected vertex %d in result", vertex)
			}
		}
	})

	t.Run("Traverse with callback on non-existent vertex", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		var visitedVertices []int
		var visitedEdges []*Edge[int, float64]

		dfs.TraverseFrom(999, func(vertex *Vertex[int, float64], edge *Edge[int, float64]) {
			visitedVertices = append(visitedVertices, vertex.GetId())
			visitedEdges = append(visitedEdges, edge)
		})

		if len(visitedVertices) != 0 {
			t.Errorf("Expected 0 vertices for non-existent start, got %d", len(visitedVertices))
		}

		if len(visitedEdges) != 0 {
			t.Errorf("Expected 0 edges for non-existent start, got %d", len(visitedEdges))
		}
	})

	t.Run("Traverse with callback that processes edge data", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 10.0, "edge1-2")
		builder.AddEdge(2, 3, 15.0, "edge2-3")

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		var edgeCosts []float64
		var edgeData []string

		dfs.TraverseFrom(1, func(vertex *Vertex[int, float64], edge *Edge[int, float64]) {
			if edge != nil {
				edgeCosts = append(edgeCosts, edge.GetCost())
				// Get edge data from graph
				edgeDataPtr, _ := graph.GetEdgeData(edge)
				if edgeDataPtr != nil {
					edgeData = append(edgeData, *edgeDataPtr)
				}
			}
		})

		// Should have 2 edges (1->2 and 2->3)
		if len(edgeCosts) != 2 {
			t.Errorf("Expected 2 edges, got %d", len(edgeCosts))
		}

		expectedCosts := []float64{10.0, 15.0}
		for i, cost := range edgeCosts {
			if cost != expectedCosts[i] {
				t.Errorf("Expected edge cost %f, got %f", expectedCosts[i], cost)
			}
		}

		expectedData := []string{"edge1-2", "edge2-3"}
		for i, data := range edgeData {
			if data != expectedData[i] {
				t.Errorf("Expected edge data %s, got %s", expectedData[i], data)
			}
		}
	})

	t.Run("Traverse with callback that counts visits", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 10.0, "edge1-2")
		builder.AddEdge(2, 3, 15.0, "edge2-3")

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		visitCount := 0
		edgeCount := 0

		dfs.TraverseFrom(1, func(vertex *Vertex[int, float64], edge *Edge[int, float64]) {
			visitCount++
			if edge != nil {
				edgeCount++
			}
		})

		if visitCount != 3 {
			t.Errorf("Expected 3 vertex visits, got %d", visitCount)
		}

		if edgeCount != 2 {
			t.Errorf("Expected 2 edge visits, got %d", edgeCount)
		}
	})
}

func TestDFSFindPath(t *testing.T) {
	t.Run("Find path between connected vertices", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 10.0, "edge1-2")
		builder.AddEdge(2, 3, 15.0, "edge2-3")

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		path := dfs.FindPath(1, 3)

		if path == nil {
			t.Error("Expected path to exist")
			return
		}

		if len(path) != 3 {
			t.Errorf("Expected path length 3, got %d", len(path))
		}

		if path[0] != 1 || path[1] != 2 || path[2] != 3 {
			t.Errorf("Expected path [1, 2, 3], got %v", path)
		}
	})

	t.Run("Find path to same vertex", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		path := dfs.FindPath(1, 1)

		if path == nil {
			t.Error("Expected path to exist")
			return
		}

		if len(path) != 1 {
			t.Errorf("Expected path length 1, got %d", len(path))
		}

		if path[0] != 1 {
			t.Errorf("Expected path [1], got %v", path)
		}
	})

	t.Run("Find path between disconnected vertices", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 10.0, "edge1-2")

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		path := dfs.FindPath(1, 3)

		if path != nil {
			t.Error("Expected no path between disconnected vertices")
		}
	})

	t.Run("Find path with non-existent start vertex", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		path := dfs.FindPath(999, 1)

		if path != nil {
			t.Error("Expected nil path for non-existent start vertex")
		}
	})

	t.Run("Find path with non-existent end vertex", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		path := dfs.FindPath(1, 999)

		if path != nil {
			t.Error("Expected nil path for non-existent end vertex")
		}
	})
}

func TestDFSIsReachable(t *testing.T) {
	t.Run("Check reachability between connected vertices", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 10.0, "edge1-2")
		builder.AddEdge(2, 3, 15.0, "edge2-3")

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		if !dfs.IsReachable(1, 3) {
			t.Error("Expected vertex 3 to be reachable from vertex 1")
		}

		if !dfs.IsReachable(1, 2) {
			t.Error("Expected vertex 2 to be reachable from vertex 1")
		}

		if !dfs.IsReachable(2, 3) {
			t.Error("Expected vertex 3 to be reachable from vertex 2")
		}
	})

	t.Run("Check reachability to same vertex", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		if !dfs.IsReachable(1, 1) {
			t.Error("Expected vertex to be reachable from itself")
		}
	})

	t.Run("Check reachability between disconnected vertices", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 10.0, "edge1-2")

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		if dfs.IsReachable(1, 3) {
			t.Error("Expected vertex 3 to not be reachable from vertex 1")
		}

		if dfs.IsReachable(3, 1) {
			t.Error("Expected vertex 1 to not be reachable from vertex 3")
		}
	})

	t.Run("Check reachability with non-existent vertices", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		if dfs.IsReachable(999, 1) {
			t.Error("Expected false for non-existent start vertex")
		}

		if dfs.IsReachable(1, 999) {
			t.Error("Expected false for non-existent end vertex")
		}
	})
}

func TestDFSHasCycle(t *testing.T) {
	t.Run("No cycle in directed graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 10.0, "edge1-2")
		builder.AddEdge(2, 3, 15.0, "edge2-3")

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		hasCycle := dfs.HasCycle()
		if hasCycle {
			t.Error("Expected no cycle, but cycle was detected")
		}
	})

	t.Run("Cycle in directed graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 10.0, "edge1-2")
		builder.AddEdge(2, 3, 15.0, "edge2-3")
		builder.AddEdge(3, 1, 5.0, "edge3-1") // Creates a cycle

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		hasCycle := dfs.HasCycle()
		if !hasCycle {
			t.Error("Expected cycle, but none was detected")
		}
	})

	t.Run("Self-loop cycle", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddEdge(1, 1, 5.0, "self-loop") // Self-loop creates a cycle

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		hasCycle := dfs.HasCycle()
		if !hasCycle {
			t.Error("Expected cycle from self-loop, but none was detected")
		}
	})

	t.Run("Multiple components with cycle", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		// Component 1: No cycle
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddEdge(1, 2, 10.0, "edge1-2")

		// Component 2: With cycle
		builder.AddVertex(3, "C")
		builder.AddVertex(4, "D")
		builder.AddVertex(5, "E")
		builder.AddEdge(3, 4, 15.0, "edge3-4")
		builder.AddEdge(4, 5, 20.0, "edge4-5")
		builder.AddEdge(5, 3, 25.0, "edge5-3") // Creates a cycle

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		hasCycle := dfs.HasCycle()
		if !hasCycle {
			t.Error("Expected cycle in second component, but none was detected")
		}
	})

	t.Run("Complex cycle detection", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		// Create a more complex graph with multiple paths
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddVertex(4, "D")
		builder.AddVertex(5, "E")

		// Create multiple paths but no cycle
		builder.AddEdge(1, 2, 10.0, "edge1-2")
		builder.AddEdge(1, 3, 15.0, "edge1-3")
		builder.AddEdge(2, 4, 20.0, "edge2-4")
		builder.AddEdge(3, 4, 25.0, "edge3-4")
		builder.AddEdge(4, 5, 30.0, "edge4-5")

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		hasCycle := dfs.HasCycle()
		if hasCycle {
			t.Error("Expected no cycle in complex graph, but cycle was detected")
		}

		// Add a cycle
		builder.AddEdge(5, 2, 35.0, "edge5-2") // Creates cycle: 2->4->5->2

		graph = builder.BuildDirected()
		dfs = NewDFS(graph)

		hasCycle = dfs.HasCycle()
		if !hasCycle {
			t.Error("Expected cycle after adding edge 5->2, but none was detected")
		}
	})

	t.Run("Empty graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		hasCycle := dfs.HasCycle()
		if hasCycle {
			t.Error("Expected no cycle in empty graph, but cycle was detected")
		}
	})

	t.Run("Single vertex", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		hasCycle := dfs.HasCycle()
		if hasCycle {
			t.Error("Expected no cycle with single vertex, but cycle was detected")
		}
	})

	t.Run("Two vertices no cycle", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddEdge(1, 2, 10.0, "edge1-2")
		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		hasCycle := dfs.HasCycle()
		if hasCycle {
			t.Error("Expected no cycle with two vertices, but cycle was detected")
		}
	})

	t.Run("Two vertices with cycle", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddEdge(1, 2, 10.0, "edge1-2")
		builder.AddEdge(2, 1, 15.0, "edge2-1") // Creates cycle
		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		hasCycle := dfs.HasCycle()
		if !hasCycle {
			t.Error("Expected cycle with two vertices, but none was detected")
		}
	})
}

func TestDFSFindCycles(t *testing.T) {
	t.Run("No cycles in directed graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 10.0, "edge1-2")
		builder.AddEdge(2, 3, 15.0, "edge2-3")

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		cycles := dfs.FindCycles()
		if len(cycles) != 0 {
			t.Errorf("Expected no cycles, but found %d cycles: %v", len(cycles), cycles)
		}
	})

	t.Run("Single cycle in directed graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 10.0, "edge1-2")
		builder.AddEdge(2, 3, 15.0, "edge2-3")
		builder.AddEdge(3, 1, 5.0, "edge3-1") // Creates a cycle

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		cycles := dfs.FindCycles()
		if len(cycles) != 1 {
			t.Errorf("Expected 1 cycle, but found %d cycles: %v", len(cycles), cycles)
		}

		// Check that the cycle contains the expected vertices
		cycle := cycles[0]
		expectedVertices := map[int]bool{1: true, 2: true, 3: true}
		if len(cycle) != 3 {
			t.Errorf("Expected cycle length 3, got %d", len(cycle))
		}
		for _, vertex := range cycle {
			if !expectedVertices[vertex] {
				t.Errorf("Unexpected vertex %d in cycle", vertex)
			}
		}
	})

	t.Run("Self-loop cycle", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddEdge(1, 1, 5.0, "self-loop") // Self-loop creates a cycle

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		cycles := dfs.FindCycles()
		if len(cycles) != 1 {
			t.Errorf("Expected 1 cycle, but found %d cycles: %v", len(cycles), cycles)
		}

		// Check that the cycle is a self-loop
		cycle := cycles[0]
		if len(cycle) != 1 || cycle[0] != 1 {
			t.Errorf("Expected self-loop cycle [1], got %v", cycle)
		}
	})

	t.Run("Multiple cycles", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		// Cycle 1: 1->2->1
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddEdge(1, 2, 10.0, "edge1-2")
		builder.AddEdge(2, 1, 15.0, "edge2-1")

		// Cycle 2: 3->4->5->3
		builder.AddVertex(3, "C")
		builder.AddVertex(4, "D")
		builder.AddVertex(5, "E")
		builder.AddEdge(3, 4, 20.0, "edge3-4")
		builder.AddEdge(4, 5, 25.0, "edge4-5")
		builder.AddEdge(5, 3, 30.0, "edge5-3")

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		cycles := dfs.FindCycles()
		if len(cycles) < 2 {
			t.Errorf("Expected at least 2 cycles, but found %d cycles: %v", len(cycles), cycles)
		}

		// Verify that we have cycles of the expected lengths
		cycleLengths := make(map[int]int)
		for _, cycle := range cycles {
			cycleLengths[len(cycle)]++
		}

		// Should have at least one 2-vertex cycle and one 3-vertex cycle
		if cycleLengths[2] == 0 || cycleLengths[3] == 0 {
			t.Errorf("Expected cycles of length 2 and 3, got cycle lengths: %v", cycleLengths)
		}
	})

	t.Run("Complex graph with one cycle", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		// Create a complex graph with multiple paths but only one cycle
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddVertex(4, "D")
		builder.AddVertex(5, "E")

		// Create multiple paths but no cycle initially
		builder.AddEdge(1, 2, 10.0, "edge1-2")
		builder.AddEdge(1, 3, 15.0, "edge1-3")
		builder.AddEdge(2, 4, 20.0, "edge2-4")
		builder.AddEdge(3, 4, 25.0, "edge3-4")
		builder.AddEdge(4, 5, 30.0, "edge4-5")

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		cycles := dfs.FindCycles()
		if len(cycles) != 0 {
			t.Errorf("Expected no cycles initially, but found %d cycles: %v", len(cycles), cycles)
		}

		// Add a cycle
		builder.AddEdge(5, 2, 35.0, "edge5-2") // Creates cycle: 2->4->5->2

		graph = builder.BuildDirected()
		dfs = NewDFS(graph)

		cycles = dfs.FindCycles()
		if len(cycles) < 1 {
			t.Errorf("Expected at least 1 cycle after adding edge 5->2, but found %d cycles: %v", len(cycles), cycles)
		}

		// Check that at least one cycle contains the expected vertices
		foundExpectedCycle := false
		expectedVertices := map[int]bool{2: true, 4: true, 5: true}
		for _, cycle := range cycles {
			if len(cycle) == 3 {
				allExpected := true
				for _, vertex := range cycle {
					if !expectedVertices[vertex] {
						allExpected = false
						break
					}
				}
				if allExpected {
					foundExpectedCycle = true
					break
				}
			}
		}
		if !foundExpectedCycle {
			t.Errorf("Expected to find a cycle with vertices [2, 4, 5], but found cycles: %v", cycles)
		}
	})

	t.Run("Empty graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		cycles := dfs.FindCycles()
		if len(cycles) != 0 {
			t.Errorf("Expected no cycles in empty graph, but found %d cycles: %v", len(cycles), cycles)
		}
	})

	t.Run("Single vertex", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		cycles := dfs.FindCycles()
		if len(cycles) != 0 {
			t.Errorf("Expected no cycles with single vertex, but found %d cycles: %v", len(cycles), cycles)
		}
	})

	t.Run("Two vertices no cycle", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddEdge(1, 2, 10.0, "edge1-2")
		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		cycles := dfs.FindCycles()
		if len(cycles) != 0 {
			t.Errorf("Expected no cycles with two vertices, but found %d cycles: %v", len(cycles), cycles)
		}
	})

	t.Run("Two vertices with cycle", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddEdge(1, 2, 10.0, "edge1-2")
		builder.AddEdge(2, 1, 15.0, "edge2-1") // Creates cycle
		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		cycles := dfs.FindCycles()
		if len(cycles) != 1 {
			t.Errorf("Expected 1 cycle with two vertices, but found %d cycles: %v", len(cycles), cycles)
		}

		// Check that the cycle contains both vertices
		cycle := cycles[0]
		expectedVertices := map[int]bool{1: true, 2: true}
		if len(cycle) != 2 {
			t.Errorf("Expected cycle length 2, got %d", len(cycle))
		}
		for _, vertex := range cycle {
			if !expectedVertices[vertex] {
				t.Errorf("Unexpected vertex %d in cycle", vertex)
			}
		}
	})

	t.Run("HasCycle consistency with FindCycles", func(t *testing.T) {
		// Test that HasCycle() is consistent with FindCycles()
		testCases := []struct {
			name     string
			vertices []int
			edges    [][3]int // [from, to, cost]
			hasCycle bool
		}{
			{
				name:     "No cycle",
				vertices: []int{1, 2, 3},
				edges:    [][3]int{{1, 2, 10}, {2, 3, 15}},
				hasCycle: false,
			},
			{
				name:     "Simple cycle",
				vertices: []int{1, 2, 3},
				edges:    [][3]int{{1, 2, 10}, {2, 3, 15}, {3, 1, 5}},
				hasCycle: true,
			},
			{
				name:     "Self-loop",
				vertices: []int{1},
				edges:    [][3]int{{1, 1, 5}},
				hasCycle: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				builder := &Builder[int, float64, string, string]{}
				for _, v := range tc.vertices {
					builder.AddVertex(v, fmt.Sprintf("V%d", v))
				}
				for _, edge := range tc.edges {
					builder.AddEdge(edge[0], edge[1], float64(edge[2]), fmt.Sprintf("edge%d-%d", edge[0], edge[1]))
				}

				graph := builder.BuildDirected()
				dfs := NewDFS(graph)

				cycles := dfs.FindCycles()
				hasCycle := dfs.HasCycle()

				expectedCycleCount := 0
				if tc.hasCycle {
					expectedCycleCount = 1
				}

				if len(cycles) != expectedCycleCount {
					t.Errorf("Expected %d cycles, got %d", expectedCycleCount, len(cycles))
				}

				if hasCycle != tc.hasCycle {
					t.Errorf("Expected HasCycle() to return %v, got %v", tc.hasCycle, hasCycle)
				}

				// Verify consistency
				if hasCycle != (len(cycles) > 0) {
					t.Errorf("HasCycle() (%v) inconsistent with FindCycles() length (%d)", hasCycle, len(cycles))
				}
			})
		}
	})
}

func TestDFSEmptyGraph(t *testing.T) {
	t.Run("DFS operations on empty graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		// All operations should handle empty graph gracefully
		if dfs.GetAllReachable(1) != nil {
			t.Error("Expected nil for traverse from non-existent vertex in empty graph")
		}

		if dfs.FindPath(1, 2) != nil {
			t.Error("Expected nil for find path in empty graph")
		}

		if dfs.IsReachable(1, 2) {
			t.Error("Expected false for reachability in empty graph")
		}

		if dfs.GetAllReachable(1) != nil {
			t.Error("Expected nil for get all reachable in empty graph")
		}
	})
}

func TestDFSIterativeDeepGraph(t *testing.T) {
	t.Run("Iterative DFS with deep graph", func(t *testing.T) {
		// Create a deep chain graph to test iterative implementation
		builder := &Builder[int, float64, string, string]{}

		// Create a chain of 1000 vertices (deep enough to test stack vs heap)
		for i := 1; i <= 1000; i++ {
			builder.AddVertex(i, "Vertex")
			if i > 1 {
				builder.AddEdge(i-1, i, 1.0, "Edge")
			}
		}

		graph := builder.BuildDirected()
		dfs := NewDFS(graph)

		// Test traversal from start
		result := dfs.GetAllReachable(1)
		if len(result) != 1000 {
			t.Errorf("Expected 1000 vertices, got %d", len(result))
		}

		// Test path finding in deep graph
		path := dfs.FindPath(1, 1000)
		if path == nil {
			t.Error("Expected path from 1 to 1000")
		} else if len(path) != 1000 {
			t.Errorf("Expected path length 1000, got %d", len(path))
		}

		// Test reachability
		if !dfs.IsReachable(1, 1000) {
			t.Error("Expected 1000 to be reachable from 1")
		}

		// Test all reachable
		reachable := dfs.GetAllReachable(1)
		if len(reachable) != 1000 {
			t.Errorf("Expected 1000 reachable vertices, got %d", len(reachable))
		}
	})
}
