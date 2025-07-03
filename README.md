# Performance Optimized Graph for Go

A high-performance, memory-efficient graph implementation for Go that can process graphs with millions of edges and vertices in milliseconds.

## Features

- **Memory Optimization**: Uses the [flyweight pattern](https://refactoring.guru/design-patterns/flyweight) for efficient memory usage
- **Minimal Allocations**: Algorithms reuse memory between operations, and the Builder optimizes large graph loading
- **Heap-based Priority Queues**: Optimized search algorithms with efficient priority queue implementations
- **Runtime Weight Computation**: Edge weights can be computed dynamically based on vertex/edge data
- **Thread Safety**: Static graph data is immutable and can be safely shared between goroutines
- **Generic Design**: Type-safe implementation with flexible ID and cost types

## Installation

```bash
go get github.com/lexkrstn/go-graph
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/lexkrstn/go-graph"
)

func main() {
    // Create a builder for a simple graph with string IDs and float costs
    builder := &graph.Builder[string, float64, struct{}, struct{}]{}
    
    // Add vertices (optional - they'll be created automatically from edges)
    builder.AddVertex("A", struct{}{})
    builder.AddVertex("B", struct{}{})
    builder.AddVertex("C", struct{}{})
    
    // Add directed edges
    builder.AddEdge("A", "B", 5.0, struct{}{})
    builder.AddEdge("B", "C", 3.0, struct{}{})
    builder.AddEdge("A", "C", 10.0, struct{}{})
    
    // Build the graph
    g := builder.BuildDirected()
    
    fmt.Printf("Graph has %d vertices and %d edges\n", 
        g.GetVertexCount(), g.GetEdgeCount())
}
```

## Type System

### Core Types

The library uses Go generics with the following type constraints:

#### `Id` - Vertex Identifier Types
```go
type Id interface {
    SInt | UInt | string | rune
}
```
- **SInt**: `int`, `int8`, `int16`, `int32`, `int64`
- **UInt**: `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- **string**: String identifiers
- **rune**: Unicode character identifiers

#### `Cost` - Edge Cost Types
```go
type Cost interface {
    SInt | UInt | Float
}
```
- **SInt**: `int`, `int8`, `int16`, `int32`, `int64`
- **UInt**: `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- **Float**: `float32`, `float64`

#### `Graph[I, C, V, E]`
The main graph type with generic parameters:
- `I`: Vertex ID type (must satisfy `Id` constraint)
- `C`: Edge cost type (must satisfy `Cost` constraint)
- `V`: Custom vertex data type
- `E`: Custom edge data type

## Usage Examples

### 1. Simple Graph with Integer IDs

```go
// Create a simple graph with integer vertices and float costs
builder := &graph.Builder[int, float64, struct{}, struct{}]{}

// Add edges (vertices created automatically)
builder.AddEdge(1, 2, 5.0, struct{}{})
builder.AddEdge(2, 3, 3.0, struct{}{})
builder.AddEdge(1, 3, 10.0, struct{}{})
builder.AddEdge(3, 4, 2.0, struct{}{})

// Build the graph
g := builder.BuildDirected()

fmt.Printf("Vertices: %d, Edges: %d\n", g.GetVertexCount(), g.GetEdgeCount())

// Get a vertex by ID
vertex, err := g.GetVertexById(1)
if err != nil {
    log.Fatal(err)
}

// Access vertex edges
edges := vertex.GetEdges()
for _, edge := range edges {
    fmt.Printf("Edge to %d with cost %.1f\n", 
        edge.GetTargetVertex().GetId(), edge.GetCost())
}
```

### 2. Graph with Custom Data

```go
// Define custom data types
type CityData struct {
    Name     string
    Population int
}

type RoadData struct {
    Name        string
    SpeedLimit  int
    IsHighway   bool
}

// Create graph with custom data
builder := &graph.Builder[string, int, CityData, RoadData]{}

// Add cities with data
builder.AddVertex("NYC", CityData{Name: "New York City", Population: 8336817})
builder.AddVertex("LA", CityData{Name: "Los Angeles", Population: 3979576})
builder.AddVertex("CHI", CityData{Name: "Chicago", Population: 2693976})

// Add roads with data
builder.AddEdge("NYC", "LA", 2789, RoadData{Name: "I-80", SpeedLimit: 70, IsHighway: true})
builder.AddEdge("LA", "CHI", 2004, RoadData{Name: "I-40", SpeedLimit: 75, IsHighway: true})
builder.AddEdge("CHI", "NYC", 787, RoadData{Name: "I-90", SpeedLimit: 65, IsHighway: true})

g := builder.BuildDirected()

// Access custom data
vertex, _ := g.GetVertexById("NYC")
cityData, _ := g.GetVertexData(vertex)
fmt.Printf("City: %s, Population: %d\n", cityData.Name, cityData.Population)
```

### 3. Using DTOs for Data Transfer

```go
// Create edge DTOs
edge1 := &graph.BasicEdgeDto[int, float64, string]{
    Origin: 1,
    Target: 2,
    Cost:   5.0,
    Data:   "Main Street",
}

edge2 := &graph.BasicEdgeDto[int, float64, string]{
    Origin: 2,
    Target: 3,
    Cost:   3.0,
    Data:   "Side Street",
}

// Create vertex DTOs
vertex1 := &graph.BasicVertexDto[int, string]{
    Id:   1,
    Data: "Downtown",
}

// Add to builder
builder := &graph.Builder[int, float64, string, string]{}
builder.AddEdgeDto(edge1)
builder.AddEdgeDto(edge2)
builder.AddVertexDto(vertex1)

g := builder.BuildDirected()

// Export all data as DTOs
vertices := g.GetAllVertices(func() graph.VertexDto[int, string] {
    return &graph.BasicVertexDto[int, string]{}
})

edges := g.GetAllEdges(func() graph.EdgeDto[int, float64, string] {
    return &graph.BasicEdgeDto[int, float64, string]{}
})

// DTOs can be easily serialized to JSON
```

### 4. Graph Traversal

```go
builder := &graph.Builder[int, float64, string, bool]{}
builder.AddEdge(1, 2, 5.0, true)
builder.AddEdge(2, 3, 3.0, false)
builder.AddEdge(1, 3, 10.0, true)

g := builder.BuildDirected()

// Visit all vertices
g.VisitVertices(func(vertex *graph.Vertex[int, float64]) {
    fmt.Printf("Vertex %d has %d outgoing edges\n", 
        vertex.GetId(), len(vertex.GetEdges()))
})

// Visit all edges
g.VisitEdges(func(vertex *graph.Vertex[int, float64], edge *graph.Edge[int, float64]) {
    fmt.Printf("Edge from %d to %d with cost %.1f\n", 
        vertex.GetId(), edge.GetTargetVertex().GetId(), edge.GetCost())
})

// Check if any edge satisfies a condition
hasExpensiveEdge := g.SomeEdges(func(vertex *graph.Vertex[int, float64], edge *graph.Edge[int, float64]) bool {
    return edge.GetCost() > 8.0
})

// Check if all edges satisfy a condition
allEdgesPositive := g.EveryEdge(func(vertex *graph.Vertex[int, float64], edge *graph.Edge[int, float64]) bool {
    return edge.GetCost() > 0
})
```

## Performance Characteristics

- **Memory Usage**: Optimized for large graphs with millions of vertices/edges
- **Vertex Lookup**: O(1) time complexity using hash maps
- **Edge Traversal**: O(1) per edge with minimal memory allocations
- **Graph Construction**: Bulk allocation reduces memory fragmentation
- **Thread Safety**: Immutable graph structure allows safe concurrent access

## TODO

- `Graph.Clone()`, `Vertex.Clone()`, etc.
- `Builder.BuildBiDirected()`
- `MakeBuilderFromGraph()` - use visitors to collect edges and vertices
- Algorithms:
  - `Dijkstra`
  - `AStar`
  - `BellmanFord`
  - `ConnectedComponents`

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
