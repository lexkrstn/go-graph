package graph

import (
	"testing"
)

func TestGraph(t *testing.T) {
	t.Run("Empty graph", func(t *testing.T) {
		graph := &Graph[int, float64, string, bool]{
			vertices:         []Vertex[int, float64]{},
			idToIndex:        make(map[int]int),
			customVertexData: []string{},
			customEdgeData:   []bool{},
			edgeCount:        0,
			biEdgeCount:      0,
		}

		if graph.GetVertexCount() != 0 {
			t.Errorf("Expected vertex count 0, got %d", graph.GetVertexCount())
		}

		if graph.GetEdgeCount() != 0 {
			t.Errorf("Expected edge count 0, got %d", graph.GetEdgeCount())
		}

		if graph.GetBiEdgeCount() != 0 {
			t.Errorf("Expected bidirectional edge count 0, got %d", graph.GetBiEdgeCount())
		}
	})

	t.Run("Get vertex by ID", func(t *testing.T) {
		graph := &Graph[int, float64, string, bool]{
			vertices: []Vertex[int, float64]{
				{id: 1, customDataIndex: 0, edges: []Edge[int, float64]{}},
				{id: 2, customDataIndex: 1, edges: []Edge[int, float64]{}},
			},
			idToIndex: map[int]int{
				1: 0,
				2: 1,
			},
			customVertexData: []string{"vertex1", "vertex2"},
			customEdgeData:   []bool{},
			edgeCount:        0,
			biEdgeCount:      0,
		}

		vertex, err := graph.GetVertexById(1)
		if err != nil {
			t.Errorf("Failed to get vertex 1: %v", err)
		}
		if vertex.GetId() != 1 {
			t.Errorf("Expected vertex ID 1, got %v", vertex.GetId())
		}

		// Test non-existent vertex
		_, err = graph.GetVertexById(999)
		if err == nil {
			t.Error("Expected error for non-existent vertex")
		}
	})

	t.Run("Get vertex by index", func(t *testing.T) {
		graph := &Graph[int, float64, string, bool]{
			vertices: []Vertex[int, float64]{
				{id: 1, customDataIndex: 0, edges: []Edge[int, float64]{}},
				{id: 2, customDataIndex: 1, edges: []Edge[int, float64]{}},
			},
			idToIndex: map[int]int{
				1: 0,
				2: 1,
			},
			customVertexData: []string{"vertex1", "vertex2"},
			customEdgeData:   []bool{},
			edgeCount:        0,
			biEdgeCount:      0,
		}

		vertex, err := graph.GetVertexByIndex(0)
		if err != nil {
			t.Errorf("Failed to get vertex at index 0: %v", err)
		}
		if vertex.GetId() != 1 {
			t.Errorf("Expected vertex ID 1, got %v", vertex.GetId())
		}

		// Test out of range index
		_, err = graph.GetVertexByIndex(999)
		if err == nil {
			t.Error("Expected error for out of range index")
		}

		_, err = graph.GetVertexByIndex(-1)
		if err == nil {
			t.Error("Expected error for negative index")
		}
	})

	t.Run("Get vertex data", func(t *testing.T) {
		graph := &Graph[int, float64, string, bool]{
			vertices: []Vertex[int, float64]{
				{id: 1, customDataIndex: 0, edges: []Edge[int, float64]{}},
			},
			idToIndex: map[int]int{
				1: 0,
			},
			customVertexData: []string{"vertex1"},
			customEdgeData:   []bool{},
			edgeCount:        0,
			biEdgeCount:      0,
		}

		vertex, _ := graph.GetVertexById(1)
		data, err := graph.GetVertexData(vertex)
		if err != nil {
			t.Errorf("Failed to get vertex data: %v", err)
		}
		if *data != "vertex1" {
			t.Errorf("Expected vertex data 'vertex1', got %v", *data)
		}

		// Test nil vertex
		_, err = graph.GetVertexData(nil)
		if err == nil {
			t.Error("Expected error for nil vertex")
		}
	})

	t.Run("Get edge data", func(t *testing.T) {
		graph := &Graph[int, float64, string, bool]{
			vertices: []Vertex[int, float64]{
				{id: 1, customDataIndex: 0, edges: []Edge[int, float64]{}},
			},
			idToIndex: map[int]int{
				1: 0,
			},
			customVertexData: []string{"vertex1"},
			customEdgeData:   []bool{true},
			edgeCount:        0,
			biEdgeCount:      0,
		}

		edge := &Edge[int, float64]{
			cost:            10.5,
			targetVertex:    &graph.vertices[0],
			customDataIndex: 0,
		}

		data, err := graph.GetEdgeData(edge)
		if err != nil {
			t.Errorf("Failed to get edge data: %v", err)
		}
		if *data != true {
			t.Errorf("Expected edge data true, got %v", *data)
		}

		// Test nil edge
		_, err = graph.GetEdgeData(nil)
		if err == nil {
			t.Error("Expected error for nil edge")
		}
	})

	t.Run("Get all vertices", func(t *testing.T) {
		graph := &Graph[int, float64, string, bool]{
			vertices: []Vertex[int, float64]{
				{id: 1, customDataIndex: 0, edges: []Edge[int, float64]{}},
				{id: 2, customDataIndex: 1, edges: []Edge[int, float64]{}},
			},
			idToIndex: map[int]int{
				1: 0,
				2: 1,
			},
			customVertexData: []string{"vertex1", "vertex2"},
			customEdgeData:   []bool{},
			edgeCount:        0,
			biEdgeCount:      0,
		}

		dtos := graph.GetAllVertices(func() VertexDto[int, string] {
			return &BasicVertexDto[int, string]{}
		})

		if len(dtos) != 2 {
			t.Errorf("Expected 2 vertex DTOs, got %d", len(dtos))
		}

		// Check first vertex
		if dtos[0].GetId() != 1 {
			t.Errorf("Expected first vertex ID 1, got %v", dtos[0].GetId())
		}
		if dtos[0].GetData() != "vertex1" {
			t.Errorf("Expected first vertex data 'vertex1', got %v", dtos[0].GetData())
		}

		// Check second vertex
		if dtos[1].GetId() != 2 {
			t.Errorf("Expected second vertex ID 2, got %v", dtos[1].GetId())
		}
		if dtos[1].GetData() != "vertex2" {
			t.Errorf("Expected second vertex data 'vertex2', got %v", dtos[1].GetData())
		}
	})

	t.Run("Get all edges", func(t *testing.T) {
		targetVertex := &Vertex[int, float64]{
			id:              2,
			customDataIndex: 1,
			edges:           []Edge[int, float64]{},
		}

		graph := &Graph[int, float64, string, bool]{
			vertices: []Vertex[int, float64]{
				{
					id:              1,
					customDataIndex: 0,
					edges: []Edge[int, float64]{
						{
							cost:            10.5,
							targetVertex:    targetVertex,
							customDataIndex: 0,
						},
					},
				},
				*targetVertex,
			},
			idToIndex: map[int]int{
				1: 0,
				2: 1,
			},
			customVertexData: []string{"vertex1", "vertex2"},
			customEdgeData:   []bool{true},
			edgeCount:        1,
			biEdgeCount:      1,
		}

		dtos := graph.GetAllEdges(func() EdgeDto[int, float64, bool] {
			return &BasicEdgeDto[int, float64, bool]{}
		})

		if len(dtos) != 1 {
			t.Errorf("Expected 1 edge DTO, got %d", len(dtos))
		}

		if dtos[0].GetOrigin() != 1 {
			t.Errorf("Expected edge origin 1, got %v", dtos[0].GetOrigin())
		}
		if dtos[0].GetTarget() != 2 {
			t.Errorf("Expected edge target 2, got %v", dtos[0].GetTarget())
		}
		if dtos[0].GetCost() != 10.5 {
			t.Errorf("Expected edge cost 10.5, got %v", dtos[0].GetCost())
		}
		if dtos[0].GetData() != true {
			t.Errorf("Expected edge data true, got %v", dtos[0].GetData())
		}
	})

	t.Run("Get all bidirectional edges", func(t *testing.T) {
		targetVertex := &Vertex[int, float64]{
			id:              2,
			customDataIndex: 1,
			edges:           []Edge[int, float64]{},
		}

		graph := &Graph[int, float64, string, bool]{
			vertices: []Vertex[int, float64]{
				{
					id:              1,
					customDataIndex: 0,
					edges: []Edge[int, float64]{
						{
							cost:            10.5,
							targetVertex:    targetVertex,
							customDataIndex: 0,
						},
					},
				},
				*targetVertex,
			},
			idToIndex: map[int]int{
				1: 0,
				2: 1,
			},
			customVertexData: []string{"vertex1", "vertex2"},
			customEdgeData:   []bool{true},
			edgeCount:        1,
			biEdgeCount:      1,
		}

		dtos := graph.GetAllBiEdges(func() EdgeDto[int, float64, bool] {
			return &BasicEdgeDto[int, float64, bool]{}
		})

		if len(dtos) != 1 {
			t.Errorf("Expected 1 bidirectional edge DTO, got %d", len(dtos))
		}

		// For bidirectional edges, the origin should be the smaller ID
		if dtos[0].GetOrigin() != 1 {
			t.Errorf("Expected bidirectional edge origin 1, got %v", dtos[0].GetOrigin())
		}
		if dtos[0].GetTarget() != 2 {
			t.Errorf("Expected bidirectional edge target 2, got %v", dtos[0].GetTarget())
		}
	})

	t.Run("Visit edges", func(t *testing.T) {
		targetVertex := &Vertex[int, float64]{
			id:              2,
			customDataIndex: 1,
			edges:           []Edge[int, float64]{},
		}

		graph := &Graph[int, float64, string, bool]{
			vertices: []Vertex[int, float64]{
				{
					id:              1,
					customDataIndex: 0,
					edges: []Edge[int, float64]{
						{
							cost:            10.5,
							targetVertex:    targetVertex,
							customDataIndex: 0,
						},
					},
				},
				*targetVertex,
			},
			idToIndex: map[int]int{
				1: 0,
				2: 1,
			},
			customVertexData: []string{"vertex1", "vertex2"},
			customEdgeData:   []bool{true},
			edgeCount:        1,
			biEdgeCount:      1,
		}

		visitedCount := 0
		graph.VisitEdges(func(vertex *Vertex[int, float64], edge *Edge[int, float64]) {
			visitedCount++
			if vertex.GetId() != 1 {
				t.Errorf("Expected visited vertex ID 1, got %v", vertex.GetId())
			}
			if edge.GetCost() != 10.5 {
				t.Errorf("Expected visited edge cost 10.5, got %v", edge.GetCost())
			}
		})

		if visitedCount != 1 {
			t.Errorf("Expected 1 edge visit, got %d", visitedCount)
		}
	})

	t.Run("Some edges", func(t *testing.T) {
		targetVertex := &Vertex[int, float64]{
			id:              2,
			customDataIndex: 1,
			edges:           []Edge[int, float64]{},
		}

		graph := &Graph[int, float64, string, bool]{
			vertices: []Vertex[int, float64]{
				{
					id:              1,
					customDataIndex: 0,
					edges: []Edge[int, float64]{
						{
							cost:            10.5,
							targetVertex:    targetVertex,
							customDataIndex: 0,
						},
					},
				},
				*targetVertex,
			},
			idToIndex: map[int]int{
				1: 0,
				2: 1,
			},
			customVertexData: []string{"vertex1", "vertex2"},
			customEdgeData:   []bool{true},
			edgeCount:        1,
			biEdgeCount:      1,
		}

		// Test predicate that returns true
		result := graph.SomeEdges(func(vertex *Vertex[int, float64], edge *Edge[int, float64]) bool {
			return edge.GetCost() == 10.5
		})
		if !result {
			t.Error("Expected SomeEdges to return true for matching predicate")
		}

		// Test predicate that returns false
		result = graph.SomeEdges(func(vertex *Vertex[int, float64], edge *Edge[int, float64]) bool {
			return edge.GetCost() == 999.0
		})
		if result {
			t.Error("Expected SomeEdges to return false for non-matching predicate")
		}
	})

	t.Run("Every edge", func(t *testing.T) {
		targetVertex := &Vertex[int, float64]{
			id:              2,
			customDataIndex: 1,
			edges:           []Edge[int, float64]{},
		}

		graph := &Graph[int, float64, string, bool]{
			vertices: []Vertex[int, float64]{
				{
					id:              1,
					customDataIndex: 0,
					edges: []Edge[int, float64]{
						{
							cost:            10.5,
							targetVertex:    targetVertex,
							customDataIndex: 0,
						},
					},
				},
				*targetVertex,
			},
			idToIndex: map[int]int{
				1: 0,
				2: 1,
			},
			customVertexData: []string{"vertex1", "vertex2"},
			customEdgeData:   []bool{true},
			edgeCount:        1,
			biEdgeCount:      1,
		}

		// Test predicate that returns true for all edges
		result := graph.EveryEdge(func(vertex *Vertex[int, float64], edge *Edge[int, float64]) bool {
			return edge.GetCost() > 0
		})
		if !result {
			t.Error("Expected EveryEdge to return true for all edges matching predicate")
		}

		// Test predicate that returns false for some edges
		result = graph.EveryEdge(func(vertex *Vertex[int, float64], edge *Edge[int, float64]) bool {
			return edge.GetCost() > 20.0
		})
		if result {
			t.Error("Expected EveryEdge to return false when some edges don't match predicate")
		}
	})

	t.Run("Visit vertices", func(t *testing.T) {
		graph := &Graph[int, float64, string, bool]{
			vertices: []Vertex[int, float64]{
				{id: 1, customDataIndex: 0, edges: []Edge[int, float64]{}},
				{id: 2, customDataIndex: 1, edges: []Edge[int, float64]{}},
			},
			idToIndex: map[int]int{
				1: 0,
				2: 1,
			},
			customVertexData: []string{"vertex1", "vertex2"},
			customEdgeData:   []bool{},
			edgeCount:        0,
			biEdgeCount:      0,
		}

		visitedCount := 0
		graph.VisitVertices(func(vertex *Vertex[int, float64]) {
			visitedCount++
		})

		if visitedCount != 2 {
			t.Errorf("Expected 2 vertex visits, got %d", visitedCount)
		}
	})

	t.Run("Some vertices", func(t *testing.T) {
		graph := &Graph[int, float64, string, bool]{
			vertices: []Vertex[int, float64]{
				{id: 1, customDataIndex: 0, edges: []Edge[int, float64]{}},
				{id: 2, customDataIndex: 1, edges: []Edge[int, float64]{}},
			},
			idToIndex: map[int]int{
				1: 0,
				2: 1,
			},
			customVertexData: []string{"vertex1", "vertex2"},
			customEdgeData:   []bool{},
			edgeCount:        0,
			biEdgeCount:      0,
		}

		// Test predicate that returns true
		result := graph.SomeVertices(func(vertex *Vertex[int, float64]) bool {
			return vertex.GetId() == 1
		})
		if !result {
			t.Error("Expected SomeVertices to return true for matching predicate")
		}

		// Test predicate that returns false
		result = graph.SomeVertices(func(vertex *Vertex[int, float64]) bool {
			return vertex.GetId() == 999
		})
		if result {
			t.Error("Expected SomeVertices to return false for non-matching predicate")
		}
	})

	t.Run("Every vertex", func(t *testing.T) {
		graph := &Graph[int, float64, string, bool]{
			vertices: []Vertex[int, float64]{
				{id: 1, customDataIndex: 0, edges: []Edge[int, float64]{}},
				{id: 2, customDataIndex: 1, edges: []Edge[int, float64]{}},
			},
			idToIndex: map[int]int{
				1: 0,
				2: 1,
			},
			customVertexData: []string{"vertex1", "vertex2"},
			customEdgeData:   []bool{},
			edgeCount:        0,
			biEdgeCount:      0,
		}

		// Test predicate that returns true for all vertices
		result := graph.EveryVertex(func(vertex *Vertex[int, float64]) bool {
			return vertex.GetId() > 0
		})
		if !result {
			t.Error("Expected EveryVertex to return true for all vertices matching predicate")
		}

		// Test predicate that returns false for some vertices
		result = graph.EveryVertex(func(vertex *Vertex[int, float64]) bool {
			return vertex.GetId() > 5
		})
		if result {
			t.Error("Expected EveryVertex to return false when some vertices don't match predicate")
		}
	})
}
