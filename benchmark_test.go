package graph

import (
	"testing"
)

func BenchmarkBuilderAddVertex(b *testing.B) {
	builder := &Builder[int, float64, string, bool]{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder.AddVertex(i, "vertex")
	}
}

func BenchmarkBuilderAddEdge(b *testing.B) {
	builder := &Builder[int, float64, string, bool]{}

	// Pre-add vertices
	for i := 0; i < 1000; i++ {
		builder.AddVertex(i, "vertex")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder.AddEdge(i%1000, (i+1)%1000, float64(i), true)
	}
}

func BenchmarkBuilderAddBiEdge(b *testing.B) {
	builder := &Builder[int, float64, string, bool]{}

	// Pre-add vertices
	for i := 0; i < 1000; i++ {
		builder.AddVertex(i, "vertex")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder.AddBiEdge(i%1000, (i+1)%1000, float64(i), true)
	}
}

func BenchmarkBuildDirected(b *testing.B) {
	builder := &Builder[int, float64, string, bool]{}

	// Pre-add vertices and edges
	for i := 0; i < 1000; i++ {
		builder.AddVertex(i, "vertex")
	}

	for i := 0; i < 5000; i++ {
		builder.AddEdge(i%1000, (i+1)%1000, float64(i), true)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = builder.BuildDirected()
	}
}

func BenchmarkGetVertexById(b *testing.B) {
	builder := &Builder[int, float64, string, bool]{}

	// Build a graph with 1000 vertices
	for i := 0; i < 1000; i++ {
		builder.AddVertex(i, "vertex")
	}

	graph := builder.BuildDirected()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = graph.GetVertexById(i % 1000)
	}
}

func BenchmarkGetVertexByIndex(b *testing.B) {
	builder := &Builder[int, float64, string, bool]{}

	// Build a graph with 1000 vertices
	for i := 0; i < 1000; i++ {
		builder.AddVertex(i, "vertex")
	}

	graph := builder.BuildDirected()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = graph.GetVertexByIndex(i % 1000)
	}
}

func BenchmarkGetVertexData(b *testing.B) {
	builder := &Builder[int, float64, string, bool]{}

	// Build a graph with 1000 vertices
	for i := 0; i < 1000; i++ {
		builder.AddVertex(i, "vertex")
	}

	graph := builder.BuildDirected()
	vertex, _ := graph.GetVertexById(0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = graph.GetVertexData(vertex)
	}
}

func BenchmarkGetEdgeData(b *testing.B) {
	builder := &Builder[int, float64, string, bool]{}

	// Build a graph with 1000 vertices and 5000 edges
	for i := 0; i < 1000; i++ {
		builder.AddVertex(i, "vertex")
	}

	for i := 0; i < 5000; i++ {
		builder.AddEdge(i%1000, (i+1)%1000, float64(i), true)
	}

	graph := builder.BuildDirected()
	vertex, _ := graph.GetVertexById(0)
	edge := &vertex.edges[0]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = graph.GetEdgeData(edge)
	}
}

func BenchmarkVisitVertices(b *testing.B) {
	builder := &Builder[int, float64, string, bool]{}

	// Build a graph with 1000 vertices
	for i := 0; i < 1000; i++ {
		builder.AddVertex(i, "vertex")
	}

	graph := builder.BuildDirected()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.VisitVertices(func(vertex *Vertex[int, float64]) {
			// Do nothing, just visit
		})
	}
}

func BenchmarkVisitEdges(b *testing.B) {
	builder := &Builder[int, float64, string, bool]{}

	// Build a graph with 1000 vertices and 5000 edges
	for i := 0; i < 1000; i++ {
		builder.AddVertex(i, "vertex")
	}

	for i := 0; i < 5000; i++ {
		builder.AddEdge(i%1000, (i+1)%1000, float64(i), true)
	}

	graph := builder.BuildDirected()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.VisitEdges(func(vertex *Vertex[int, float64], edge *Edge[int, float64]) {
			// Do nothing, just visit
		})
	}
}

func BenchmarkSomeVertices(b *testing.B) {
	builder := &Builder[int, float64, string, bool]{}

	// Build a graph with 1000 vertices
	for i := 0; i < 1000; i++ {
		builder.AddVertex(i, "vertex")
	}

	graph := builder.BuildDirected()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = graph.SomeVertices(func(vertex *Vertex[int, float64]) bool {
			return vertex.GetId() == 999 // Will be found
		})
	}
}

func BenchmarkEveryVertex(b *testing.B) {
	builder := &Builder[int, float64, string, bool]{}

	// Build a graph with 1000 vertices
	for i := 0; i < 1000; i++ {
		builder.AddVertex(i, "vertex")
	}

	graph := builder.BuildDirected()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = graph.EveryVertex(func(vertex *Vertex[int, float64]) bool {
			return vertex.GetId() >= 0 // Always true
		})
	}
}

func BenchmarkSomeEdges(b *testing.B) {
	builder := &Builder[int, float64, string, bool]{}

	// Build a graph with 1000 vertices and 5000 edges
	for i := 0; i < 1000; i++ {
		builder.AddVertex(i, "vertex")
	}

	for i := 0; i < 5000; i++ {
		builder.AddEdge(i%1000, (i+1)%1000, float64(i), true)
	}

	graph := builder.BuildDirected()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = graph.SomeEdges(func(vertex *Vertex[int, float64], edge *Edge[int, float64]) bool {
			return edge.GetCost() == 999.0 // Will be found
		})
	}
}

func BenchmarkEveryEdge(b *testing.B) {
	builder := &Builder[int, float64, string, bool]{}

	// Build a graph with 1000 vertices and 5000 edges
	for i := 0; i < 1000; i++ {
		builder.AddVertex(i, "vertex")
	}

	for i := 0; i < 5000; i++ {
		builder.AddEdge(i%1000, (i+1)%1000, float64(i), true)
	}

	graph := builder.BuildDirected()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = graph.EveryEdge(func(vertex *Vertex[int, float64], edge *Edge[int, float64]) bool {
			return edge.GetCost() >= 0.0 // Always true
		})
	}
}

func BenchmarkGetAllVertices(b *testing.B) {
	builder := &Builder[int, float64, string, bool]{}

	// Build a graph with 1000 vertices
	for i := 0; i < 1000; i++ {
		builder.AddVertex(i, "vertex")
	}

	graph := builder.BuildDirected()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = graph.GetAllVertices(func() VertexDto[int, string] {
			return &BasicVertexDto[int, string]{}
		})
	}
}

func BenchmarkGetAllEdges(b *testing.B) {
	builder := &Builder[int, float64, string, bool]{}

	// Build a graph with 1000 vertices and 5000 edges
	for i := 0; i < 1000; i++ {
		builder.AddVertex(i, "vertex")
	}

	for i := 0; i < 5000; i++ {
		builder.AddEdge(i%1000, (i+1)%1000, float64(i), true)
	}

	graph := builder.BuildDirected()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = graph.GetAllEdges(func() EdgeDto[int, float64, bool] {
			return &BasicEdgeDto[int, float64, bool]{}
		})
	}
}

func BenchmarkGetAllBiEdges(b *testing.B) {
	builder := &Builder[int, float64, string, bool]{}

	// Build a graph with 1000 vertices and 5000 bidirectional edges
	for i := 0; i < 1000; i++ {
		builder.AddVertex(i, "vertex")
	}

	for i := 0; i < 2500; i++ { // 2500 bidirectional edges = 5000 directed edges
		builder.AddBiEdge(i%1000, (i+1)%1000, float64(i), true)
	}

	graph := builder.BuildDirected()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = graph.GetAllBiEdges(func() EdgeDto[int, float64, bool] {
			return &BasicEdgeDto[int, float64, bool]{}
		})
	}
}

func BenchmarkCountBiEdges(b *testing.B) {
	builder := &Builder[int, float64, string, bool]{}

	// Add 2500 bidirectional edges
	for i := 0; i < 2500; i++ {
		builder.AddBiEdge(i%1000, (i+1)%1000, float64(i), true)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = builder.CountBiEdges()
	}
}
