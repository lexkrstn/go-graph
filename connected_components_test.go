package graph

import (
	"testing"
)

func TestFindConnectedComponents(t *testing.T) {
	t.Run("Create ConnectedComponents for simple graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 10.0, "edge1-2")
		builder.AddEdge(2, 3, 15.0, "edge2-3")

		graph := builder.BuildDirected()
		cc := FindConnectedComponents(graph)

		if cc == nil {
			t.Error("Expected ConnectedComponents instance, got nil")
			return
		}

		if cc.graph != graph {
			t.Error("Expected ConnectedComponents graph to match input graph")
		}

		if cc.components == nil {
			t.Error("Expected components to be initialized")
		}
	})

	t.Run("Create ConnectedComponents for empty graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		graph := builder.BuildDirected()
		cc := FindConnectedComponents(graph)

		if cc == nil {
			t.Error("Expected ConnectedComponents instance for empty graph, got nil")
			return
		}

		if len(cc.components) != 0 {
			t.Errorf("Expected 0 components for empty graph, got %d", len(cc.components))
		}
	})
}

func TestGetComponents(t *testing.T) {
	t.Run("Single connected component", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 10.0, "edge1-2")
		builder.AddEdge(2, 3, 15.0, "edge2-3")

		graph := builder.BuildDirected()
		cc := FindConnectedComponents(graph)

		components := cc.GetComponents()

		if len(components) != 1 {
			t.Errorf("Expected 1 component, got %d", len(components))
		}

		if len(components[0]) != 3 {
			t.Errorf("Expected component to have 3 vertices, got %d", len(components[0]))
		}

		// Check that all vertices are in the same component
		expectedVertices := map[int]bool{1: true, 2: true, 3: true}
		for _, vertex := range components[0] {
			if !expectedVertices[vertex] {
				t.Errorf("Unexpected vertex %d in component", vertex)
			}
		}
	})

	t.Run("Multiple disconnected components", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddVertex(4, "D")
		builder.AddEdge(1, 2, 10.0, "edge1-2")
		builder.AddEdge(3, 4, 15.0, "edge3-4")
		// No edges between 1,2 and 3,4

		graph := builder.BuildDirected()
		cc := FindConnectedComponents(graph)

		components := cc.GetComponents()

		if len(components) != 2 {
			t.Errorf("Expected 2 components, got %d", len(components))
		}

		// Sort components by size for consistent testing
		if len(components[0]) > len(components[1]) {
			components[0], components[1] = components[1], components[0]
		}

		if len(components[0]) != 2 || len(components[1]) != 2 {
			t.Errorf("Expected both components to have 2 vertices, got %d and %d", len(components[0]), len(components[1]))
		}
	})

	t.Run("Single vertex components", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		// No edges at all

		graph := builder.BuildDirected()
		cc := FindConnectedComponents(graph)

		components := cc.GetComponents()

		if len(components) != 3 {
			t.Errorf("Expected 3 components, got %d", len(components))
		}

		for _, component := range components {
			if len(component) != 1 {
				t.Errorf("Expected each component to have 1 vertex, got %d", len(component))
			}
		}
	})

	t.Run("Empty graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		graph := builder.BuildDirected()
		cc := FindConnectedComponents(graph)

		components := cc.GetComponents()

		if len(components) != 0 {
			t.Errorf("Expected 0 components for empty graph, got %d", len(components))
		}
	})

	t.Run("Complex graph with multiple components", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		// Component 1: 1-2-3
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 1.0, "1-2")
		builder.AddEdge(2, 3, 1.0, "2-3")

		// Component 2: 4-5
		builder.AddVertex(4, "D")
		builder.AddVertex(5, "E")
		builder.AddEdge(4, 5, 1.0, "4-5")

		// Component 3: 6 (isolated)
		builder.AddVertex(6, "F")

		// Component 4: 7-8-9-10
		builder.AddVertex(7, "G")
		builder.AddVertex(8, "H")
		builder.AddVertex(9, "I")
		builder.AddVertex(10, "J")
		builder.AddEdge(7, 8, 1.0, "7-8")
		builder.AddEdge(8, 9, 1.0, "8-9")
		builder.AddEdge(9, 10, 1.0, "9-10")

		graph := builder.BuildDirected()
		cc := FindConnectedComponents(graph)

		components := cc.GetComponents()

		if len(components) != 4 {
			t.Errorf("Expected 4 components, got %d", len(components))
		}

		// Sort components by size for consistent testing
		for i := 0; i < len(components); i++ {
			for j := i + 1; j < len(components); j++ {
				if len(components[i]) > len(components[j]) {
					components[i], components[j] = components[j], components[i]
				}
			}
		}

		expectedSizes := []int{1, 2, 3, 4}
		for i, component := range components {
			if len(component) != expectedSizes[i] {
				t.Errorf("Expected component %d to have %d vertices, got %d", i, expectedSizes[i], len(component))
			}
		}
	})

	t.Run("Bidirectional edges", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 1.0, "1-2")
		builder.AddEdge(2, 1, 1.0, "2-1") // Bidirectional
		builder.AddEdge(2, 3, 1.0, "2-3")
		builder.AddEdge(3, 2, 1.0, "3-2") // Bidirectional

		graph := builder.BuildDirected()
		cc := FindConnectedComponents(graph)

		components := cc.GetComponents()

		if len(components) != 1 {
			t.Errorf("Expected 1 component with bidirectional edges, got %d", len(components))
		}

		if len(components[0]) != 3 {
			t.Errorf("Expected component to have 3 vertices, got %d", len(components[0]))
		}
	})
}

func TestGetComponentCount(t *testing.T) {
	t.Run("Single component", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddEdge(1, 2, 1.0, "1-2")

		graph := builder.BuildDirected()
		cc := FindConnectedComponents(graph)

		count := cc.GetComponentCount()
		if count != 1 {
			t.Errorf("Expected 1 component, got %d", count)
		}
	})

	t.Run("Multiple components", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 1.0, "1-2")
		// 3 is isolated

		graph := builder.BuildDirected()
		cc := FindConnectedComponents(graph)

		count := cc.GetComponentCount()
		if count != 2 {
			t.Errorf("Expected 2 components, got %d", count)
		}
	})

	t.Run("Empty graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		graph := builder.BuildDirected()
		cc := FindConnectedComponents(graph)

		count := cc.GetComponentCount()
		if count != 0 {
			t.Errorf("Expected 0 components for empty graph, got %d", count)
		}
	})
}

func TestIsConnected(t *testing.T) {
	t.Run("Connected graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 1.0, "1-2")
		builder.AddEdge(2, 3, 1.0, "2-3")

		graph := builder.BuildDirected()
		cc := FindConnectedComponents(graph)

		connected := cc.IsConnected()
		if !connected {
			t.Error("Expected graph to be connected")
		}
	})

	t.Run("Disconnected graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 1.0, "1-2")
		// 3 is isolated

		graph := builder.BuildDirected()
		cc := FindConnectedComponents(graph)

		connected := cc.IsConnected()
		if connected {
			t.Error("Expected graph to be disconnected")
		}
	})

	t.Run("Empty graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		graph := builder.BuildDirected()
		cc := FindConnectedComponents(graph)

		connected := cc.IsConnected()
		if connected {
			t.Error("Expected empty graph to be disconnected")
		}
	})

	t.Run("Single vertex", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")

		graph := builder.BuildDirected()
		cc := FindConnectedComponents(graph)

		connected := cc.IsConnected()
		if !connected {
			t.Error("Expected single vertex graph to be connected")
		}
	})
}

func TestGetComponentForVertex(t *testing.T) {
	t.Run("Vertex in connected component", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 1.0, "1-2")
		builder.AddEdge(2, 3, 1.0, "2-3")

		graph := builder.BuildDirected()
		cc := FindConnectedComponents(graph)

		component := cc.GetComponentForVertex(2)
		if component == nil {
			t.Error("Expected component for vertex 2, got nil")
		}

		if len(component) != 3 {
			t.Errorf("Expected component to have 3 vertices, got %d", len(component))
		}

		// Check that all expected vertices are in the component
		expectedVertices := map[int]bool{1: true, 2: true, 3: true}
		for _, vertex := range component {
			if !expectedVertices[vertex] {
				t.Errorf("Unexpected vertex %d in component", vertex)
			}
		}
	})

	t.Run("Isolated vertex", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddVertex(3, "C")
		builder.AddEdge(1, 2, 1.0, "1-2")
		// 3 is isolated

		graph := builder.BuildDirected()
		cc := FindConnectedComponents(graph)

		component := cc.GetComponentForVertex(3)
		if component == nil {
			t.Error("Expected component for vertex 3, got nil")
		}

		if len(component) != 1 {
			t.Errorf("Expected isolated vertex component to have 1 vertex, got %d", len(component))
		}

		if component[0] != 3 {
			t.Errorf("Expected component to contain vertex 3, got %d", component[0])
		}
	})

	t.Run("Non-existent vertex", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		builder.AddVertex(1, "A")
		builder.AddVertex(2, "B")
		builder.AddEdge(1, 2, 1.0, "1-2")

		graph := builder.BuildDirected()
		cc := FindConnectedComponents(graph)

		component := cc.GetComponentForVertex(999)
		if component != nil {
			t.Error("Expected nil component for non-existent vertex, got non-nil")
		}
	})

	t.Run("Empty graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}
		graph := builder.BuildDirected()
		cc := FindConnectedComponents(graph)

		component := cc.GetComponentForVertex(1)
		if component != nil {
			t.Error("Expected nil component for vertex in empty graph, got non-nil")
		}
	})
}

func TestConnectedComponentsWithDifferentTypes(t *testing.T) {
	t.Run("String IDs", func(t *testing.T) {
		builder := &Builder[string, float64, string, string]{}
		builder.AddVertex("A", "VertexA")
		builder.AddVertex("B", "VertexB")
		builder.AddVertex("C", "VertexC")
		builder.AddEdge("A", "B", 1.0, "A-B")
		builder.AddEdge("B", "C", 1.0, "B-C")

		graph := builder.BuildDirected()
		cc := FindConnectedComponents(graph)

		components := cc.GetComponents()
		if len(components) != 1 {
			t.Errorf("Expected 1 component, got %d", len(components))
		}

		expectedVertices := map[string]bool{"A": true, "B": true, "C": true}
		for _, vertex := range components[0] {
			if !expectedVertices[vertex] {
				t.Errorf("Unexpected vertex %s in component", vertex)
			}
		}
	})

	t.Run("Integer IDs with different cost types", func(t *testing.T) {
		builder := &Builder[uint32, int, string, string]{}
		builder.AddVertex(1, "First")
		builder.AddVertex(2, "Second")
		builder.AddVertex(3, "Third")
		builder.AddEdge(1, 2, 10, "1-2")
		builder.AddEdge(2, 3, 20, "2-3")

		graph := builder.BuildDirected()
		cc := FindConnectedComponents(graph)

		components := cc.GetComponents()
		if len(components) != 1 {
			t.Errorf("Expected 1 component, got %d", len(components))
		}

		expectedVertices := map[uint32]bool{1: true, 2: true, 3: true}
		for _, vertex := range components[0] {
			if !expectedVertices[vertex] {
				t.Errorf("Unexpected vertex %d in component", vertex)
			}
		}
	})
}

func TestConnectedComponentsPerformance(t *testing.T) {
	t.Run("Large connected graph", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}

		// Create a large connected graph (chain)
		numVertices := 1000
		for i := 1; i <= numVertices; i++ {
			builder.AddVertex(i, "Vertex")
		}

		for i := 1; i < numVertices; i++ {
			builder.AddEdge(i, i+1, 1.0, "edge")
		}

		graph := builder.BuildDirected()
		cc := FindConnectedComponents(graph)

		components := cc.GetComponents()
		if len(components) != 1 {
			t.Errorf("Expected 1 component for large connected graph, got %d", len(components))
		}

		if len(components[0]) != numVertices {
			t.Errorf("Expected component to have %d vertices, got %d", numVertices, len(components[0]))
		}
	})

	t.Run("Many small components", func(t *testing.T) {
		builder := &Builder[int, float64, string, string]{}

		// Create many small components (pairs)
		numComponents := 100
		for i := 0; i < numComponents; i++ {
			vertex1 := i*2 + 1
			vertex2 := i*2 + 2
			builder.AddVertex(vertex1, "Vertex")
			builder.AddVertex(vertex2, "Vertex")
			builder.AddEdge(vertex1, vertex2, 1.0, "edge")
		}

		graph := builder.BuildDirected()
		cc := FindConnectedComponents(graph)

		components := cc.GetComponents()
		if len(components) != numComponents {
			t.Errorf("Expected %d components, got %d", numComponents, len(components))
		}

		for _, component := range components {
			if len(component) != 2 {
				t.Errorf("Expected each component to have 2 vertices, got %d", len(component))
			}
		}
	})
}
