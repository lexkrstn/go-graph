package graph

// Constants defining the bulk sizes for efficient memory allocation
const edgeBulkSize = 1000   // Number of edges to allocate in each bulk
const vertexBulkSize = 1000 // Number of vertices to allocate in each bulk

// edgeBulk represents a chunk of edge DTOs for efficient memory management.
// Uses a linked list structure to avoid large slice reallocations.
type edgeBulk[I Id, C Cost, E any] struct {
	edges []EdgeDto[I, C, E] // Slice of edge DTOs in this bulk
	next  *edgeBulk[I, C, E] // Pointer to the next bulk in the chain
}

// vertexBulk represents a chunk of vertex DTOs for efficient memory management.
// Uses a linked list structure to avoid large slice reallocations.
type vertexBulk[I Id, V any] struct {
	vertices []VertexDto[I, V] // Slice of vertex DTOs in this bulk
	next     *vertexBulk[I, V] // Pointer to the next bulk in the chain
}

// Builder constructs a Graph from DTOs in an efficient way.
// Uses bulk allocation to minimize memory allocations and improve performance.
// The generic types I, C, V, E represent Id, Cost, Vertex data, and Edge data respectively.
type Builder[I Id, C Cost, V any, E any] struct {
	firstEdgeBulk       *edgeBulk[I, C, E] // First bulk in the edge bulk chain
	edgeCount           int                // Total number of edges added
	freeEdgeSlotCount   int                // Number of free slots in the current edge bulk
	firstVertexBulk     *vertexBulk[I, V]  // First bulk in the vertex bulk chain
	vertexCount         int                // Total number of vertices added
	freeVertexSlotCount int                // Number of free slots in the current vertex bulk
}

// AddEdgeDto adds a directed edge using an EdgeDto.
// Automatically allocates new bulks when the current one is full.
// This method is the primary way to add edges to the builder.
func (b *Builder[I, C, V, E]) AddEdgeDto(dto EdgeDto[I, C, E]) {
	if b.freeEdgeSlotCount == 0 {
		newEdgeBulk := &edgeBulk[I, C, E]{
			edges: make([]EdgeDto[I, C, E], 0, edgeBulkSize),
			next:  b.firstEdgeBulk,
		}
		b.firstEdgeBulk = newEdgeBulk
		b.freeEdgeSlotCount = edgeBulkSize
	}
	b.firstEdgeBulk.edges = append(b.firstEdgeBulk.edges, dto)
	b.freeEdgeSlotCount--
	b.edgeCount++
}

// AddEdge adds a directed edge with the specified parameters.
// Creates a BasicEdgeDto internally and calls AddEdgeDto.
// This is a convenience method for adding edges without creating DTOs manually.
func (b *Builder[I, C, V, E]) AddEdge(origin I, target I, cost C, data E) {
	b.AddEdgeDto(&BasicEdgeDto[I, C, E]{origin, target, cost, data})
}

// AddBiEdge adds two directed edges in both directions.
// Creates edges (origin, target) and (target, origin) with the same cost and data.
// This is useful for creating undirected or bidirectional connections.
func (b *Builder[I, C, V, E]) AddBiEdge(origin I, target I, cost C, data E) {
	b.AddEdge(origin, target, cost, data)
	b.AddEdge(target, origin, cost, data)
}

// AddVertex adds a vertex with the specified parameters.
// Creates a BasicVertexDto internally and calls AddVertexDto.
// This is a convenience method for adding vertices without creating DTOs manually.
func (b *Builder[I, C, V, E]) AddVertex(id I, data V) {
	b.AddVertexDto(&BasicVertexDto[I, V]{id, data})
}

// AddVertexDto adds a vertex using a VertexDto.
// Automatically allocates new bulks when the current one is full.
// This method is the primary way to add vertices to the builder.
func (b *Builder[I, C, V, E]) AddVertexDto(dto VertexDto[I, V]) {
	if b.freeVertexSlotCount == 0 {
		newVertexBulk := &vertexBulk[I, V]{
			vertices: make([]VertexDto[I, V], 0, vertexBulkSize),
			next:     b.firstVertexBulk,
		}
		b.firstVertexBulk = newVertexBulk
		b.freeVertexSlotCount = vertexBulkSize
	}
	b.firstVertexBulk.vertices = append(b.firstVertexBulk.vertices, dto)
	b.freeVertexSlotCount--
	b.vertexCount++
}

// biEdgeKey is used for tracking unique bidirectional edges.
// Ensures consistent ordering of vertex pairs for deduplication.
type biEdgeKey[I Id] struct{ origin, target I }

// CountBiEdges returns the total number of unique bidirectional edges.
// Counts each pair of vertices {A, B} as one edge, regardless of direction.
// Uses a map to track unique vertex pairs and deduplicate bidirectional connections.
func (b *Builder[I, C, V, E]) CountBiEdges() int {
	existing := make(map[biEdgeKey[I]]struct{}, b.edgeCount)
	for bulk := b.firstEdgeBulk; bulk != nil; bulk = bulk.next {
		for i := range bulk.edges {
			key := biEdgeKey[I]{
				origin: bulk.edges[i].GetOrigin(),
				target: bulk.edges[i].GetTarget(),
			}
			if key.origin > key.target {
				key.origin, key.target = key.target, key.origin
			}
			existing[key] = struct{}{}
		}
	}
	return len(existing)
}

// predictVertexArrayLength calculates the number of unique vertices needed.
// Considers both explicitly added vertices and vertices referenced by edges.
// Returns the total number of unique vertex IDs to allocate space for.
func (b *Builder[I, C, V, E]) predictVertexArrayLength() int {
	ids := make(map[I]struct{}, b.vertexCount)
	for bulk := b.firstEdgeBulk; bulk != nil; bulk = bulk.next {
		for i := range bulk.edges {
			ids[bulk.edges[i].GetOrigin()] = struct{}{}
			ids[bulk.edges[i].GetTarget()] = struct{}{}
		}
	}
	for bulk := b.firstVertexBulk; bulk != nil; bulk = bulk.next {
		for i := range bulk.vertices {
			ids[bulk.vertices[i].GetId()] = struct{}{}
		}
	}
	return len(ids)
}

// countOutgoingEdges calculates the number of outgoing edges for each vertex.
// Returns a map associating vertex IDs with their outgoing edge counts.
// This information is used to pre-allocate edge slices for optimal performance.
func (b *Builder[I, C, V, E]) countOutgoingEdges() map[I]int {
	counters := make(map[I]int, b.vertexCount)
	for bulk := b.firstEdgeBulk; bulk != nil; bulk = bulk.next {
		for i := range bulk.edges {
			counters[bulk.edges[i].GetOrigin()]++
		}
	}
	return counters
}

// BuildDirected creates a directed graph from the collected DTOs.
// This method should only be called once per builder instance.
// It's unsafe to call multiple times as graphs would share data structures.
// Use Graph.Clone() to create multiple instances of the same graph.
// Returns a fully constructed Graph with all vertices and edges.
func (b *Builder[I, C, V, E]) BuildDirected() *Graph[I, C, V, E] {
	vertexCount := b.predictVertexArrayLength()
	g := &Graph[I, C, V, E]{
		vertices:         make([]Vertex[I, C], vertexCount),
		idToIndex:        make(map[I]int, vertexCount),
		customVertexData: make([]V, vertexCount),
		edgeCount:        b.edgeCount,
		biEdgeCount:      b.CountBiEdges(),
		customEdgeData:   make([]E, b.edgeCount),
	}
	vertIdxCnt := 0
	edgeIdxCnt := 0
	var originIdx, targetIdx int
	var exists bool
	outgoingEdgeCnt := b.countOutgoingEdges()
	for bulk := b.firstEdgeBulk; bulk != nil; bulk = bulk.next {
		for i := range bulk.edges {
			originId := bulk.edges[i].GetOrigin()
			if originIdx, exists = g.idToIndex[originId]; !exists {
				originIdx = vertIdxCnt
				vertIdxCnt++
				g.idToIndex[originId] = originIdx
				g.vertices[originIdx].id = originId
				g.vertices[originIdx].edges = make([]Edge[I, C], 0, outgoingEdgeCnt[originId])
				g.vertices[originIdx].customDataIndex = originIdx
			}
			targetId := bulk.edges[i].GetTarget()
			if targetIdx, exists = g.idToIndex[targetId]; !exists {
				targetIdx = vertIdxCnt
				vertIdxCnt++
				g.idToIndex[targetId] = targetIdx
				g.vertices[targetIdx].id = targetId
				g.vertices[targetIdx].edges = make([]Edge[I, C], 0, outgoingEdgeCnt[targetId])
				g.vertices[targetIdx].customDataIndex = targetIdx
			}
			g.vertices[originIdx].edges = append(g.vertices[originIdx].edges, Edge[I, C]{
				cost:            bulk.edges[i].GetCost(),
				targetVertex:    &g.vertices[targetIdx],
				customDataIndex: edgeIdxCnt,
			})
			g.customEdgeData[edgeIdxCnt] = bulk.edges[i].GetData()
			edgeIdxCnt++
		}
	}
	for bulk := b.firstVertexBulk; bulk != nil; bulk = bulk.next {
		for i := range bulk.vertices {
			vertexId := bulk.vertices[i].GetId()
			if originIdx, exists = g.idToIndex[vertexId]; !exists {
				originIdx = vertIdxCnt
				vertIdxCnt++
				g.idToIndex[vertexId] = originIdx
				g.vertices[originIdx].id = vertexId
				g.vertices[originIdx].edges = make([]Edge[I, C], 0)
				g.vertices[originIdx].customDataIndex = originIdx
			}
			g.customVertexData[originIdx] = bulk.vertices[i].GetData()
		}
	}
	return g
}
