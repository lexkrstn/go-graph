package graph

// The data that is attached to the vertices by the Bellman-Ford algorithm.
type bellmanFordVertexData[I Id, C Cost] struct {
	previous *Vertex[I, C]
	cost     C
}
