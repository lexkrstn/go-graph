# Performance Optimized Graph for Go

A high-performance, memory-efficient graph implementation for Go that can process graphs with millions of edges and vertices in milliseconds.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Type System](#type-system)
  - [Core Types](#core-types)
    - [`Id` - Vertex Identifier Types](#id---vertex-identifier-types)
    - [`Cost` - Edge Cost Types](#cost---edge-cost-types)
    - [`Graph[I, C, V, E]`](#graph-i-c-v-e)
- [Usage Examples](#usage-examples)
  - [1. Simple Graph with Integer IDs](#1-simple-graph-with-integer-ids)
  - [2. Graph with Custom Data](#2-graph-with-custom-data)
  - [3. Using DTOs for Data Transfer](#3-using-dtos-for-data-transfer)
  - [4. Graph Traversal](#4-graph-traversal)
- [Pathfinding Algorithms](#pathfinding-algorithms)
  - [Dijkstra's Algorithm](#dijkstras-algorithm)
    - [Basic Usage](#basic-usage)
    - [Performance Characteristics](#performance-characteristics)
  - [A* Algorithm](#a-algorithm)
    - [Basic Usage](#basic-usage-1)
    - [Heuristic Functions](#heuristic-functions)
    - [Performance Characteristics](#performance-characteristics-1)
  - [Bellman-Ford Algorithm](#bellman-ford-algorithm)
    - [Basic Usage](#basic-usage-2)
    - [Negative Cycle Detection](#negative-cycle-detection)
    - [Performance Characteristics](#performance-characteristics-2)
  - [Advanced Features](#advanced-features)
    - [Cost Amplification](#cost-amplification)
    - [Thread Safety](#thread-safety)
  - [Algorithm Comparison](#algorithm-comparison)
- [Performance Characteristics](#performance-characteristics-2)
- [TODO](#todo)
- [Contributing](#contributing)
- [License](#license)

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

## Pathfinding Algorithms

The library provides three powerful pathfinding algorithms optimized for performance and memory efficiency.

### Dijkstra's Algorithm

Dijkstra's algorithm finds the shortest path between two vertices in a weighted graph. It guarantees the optimal solution for non-negative edge weights.

#### Basic Usage

```go
// Create a graph
builder := &graph.Builder[string, float64, struct{}, struct{}]{}
builder.AddEdge("A", "B", 4.0, struct{}{})
builder.AddEdge("A", "C", 2.0, struct{}{})
builder.AddEdge("B", "C", 1.0, struct{}{})
builder.AddEdge("B", "D", 5.0, struct{}{})
builder.AddEdge("C", "D", 8.0, struct{}{})
builder.AddEdge("C", "E", 10.0, struct{}{})
builder.AddEdge("D", "E", 2.0, struct{}{})

g := builder.BuildDirected()

// Create Dijkstra instance
dijkstra := graph.NewDijkstra(g)

// Find shortest path from A to E
path := dijkstra.FindShortestPath("A", "E")
if path != nil {
    fmt.Printf("Shortest path: %v\n", path) // Output: [A C B D E]
} else {
    fmt.Println("No path found")
}
```

#### Performance Characteristics
- **Time Complexity**: O(E log V) where E is edges and V is vertices.
- **Space Complexity**: O(V) for vertex data storage.
- **Memory Efficient**: Reuses internal data structures between calls.
- **Thread Safety**: Not thread-safe for concurrent calls: use separate instances of the algorithm, but the graph itself can be safely shared as long as you don't modify it.

### A* Algorithm

A* is an informed search algorithm that uses heuristics to find the shortest path more efficiently than Dijkstra's algorithm. It's particularly effective when you have a good heuristic function.

#### Basic Usage

```go
// Create a graph with 2D coordinates
type Position struct {
    X, Y int
}

builder := &graph.Builder[string, float64, Position, struct{}]{}
builder.AddVertex("A", Position{0, 0})
builder.AddVertex("B", Position{1, 0})
builder.AddVertex("C", Position{1, 1})
builder.AddVertex("D", Position{2, 1})
builder.AddVertex("E", Position{2, 2})

// Add edges with Euclidean distances
builder.AddEdge("A", "B", 1.0, struct{}{})
builder.AddEdge("A", "C", 1.414, struct{}{}) // sqrt(2)
builder.AddEdge("B", "C", 1.0, struct{}{})
builder.AddEdge("B", "D", 1.0, struct{}{})
builder.AddEdge("C", "D", 1.0, struct{}{})
builder.AddEdge("C", "E", 1.414, struct{}{})
builder.AddEdge("D", "E", 1.0, struct{}{})

g := builder.BuildDirected()

// Define heuristic function (Euclidean distance)
heuristic := func(current *graph.Vertex[string, float64], goal *graph.Vertex[string, float64]) float64 {
    currentPos, _ := g.GetVertexData(current)
    goalPos, _ := g.GetVertexData(goal)
    
    dx := float64(currentPos.X - goalPos.X)
    dy := float64(currentPos.Y - goalPos.Y)
    return math.Sqrt(dx*dx + dy*dy)
}

// Create A* instance
astar := graph.NewAStar(g, heuristic)

// Find shortest path from A to E
path := astar.FindShortestPath("A", "E")
if path != nil {
    fmt.Printf("A* path: %v\n", path) // Output: [A C E]
} else {
    fmt.Println("No path found")
}
```

#### Heuristic Functions

A* requires a heuristic function that estimates the cost from any vertex to the goal. The heuristic must be **admissible** (never overestimate the actual cost) for A* to guarantee optimal solutions.

##### Common Heuristic Functions

**1. Zero Heuristic (Dijkstra's Algorithm)**
```go
zeroHeuristic := func(current *graph.Vertex[I, C], goal *graph.Vertex[I, C]) C {
    var zero C
    return zero
}
```

**2. Manhattan Distance (for grid-based problems)**
```go
manhattanHeuristic := func(current *graph.Vertex[string, float64], goal *graph.Vertex[string, float64]) float64 {
    currentPos, _ := g.GetVertexData(current)
    goalPos, _ := g.GetVertexData(goal)
    return math.Abs(float64(currentPos.X-goalPos.X)) + math.Abs(float64(currentPos.Y-goalPos.Y))
}
```

**3. Euclidean Distance (for 2D coordinates)**
```go
euclideanHeuristic := func(current *graph.Vertex[string, float64], goal *graph.Vertex[string, float64]) float64 {
    currentPos, _ := g.GetVertexData(current)
    goalPos, _ := g.GetVertexData(goal)
    dx := float64(currentPos.X - goalPos.X)
    dy := float64(currentPos.Y - goalPos.Y)
    return math.Sqrt(dx*dx + dy*dy)
}
```

#### Performance Characteristics
- **Time Complexity**: O(E log V) in worst case, often much better with good heuristics
- **Space Complexity**: O(V) for vertex data storage
- **Optimality**: Guarantees optimal solution with admissible heuristics
- **Efficiency**: Typically explores fewer vertices than Dijkstra's algorithm

### Bellman-Ford Algorithm

The Bellman-Ford algorithm finds the shortest path between two vertices in a weighted graph, even with negative edge weights. It can also detect negative cycles, making it more versatile than Dijkstra's algorithm for certain applications.

#### Basic Usage

```go
// Create a graph with some negative weights
builder := &graph.Builder[string, float64, struct{}, struct{}]{}
builder.AddEdge("A", "B", 4.0, struct{}{})
builder.AddEdge("A", "C", 2.0, struct{}{})
builder.AddEdge("B", "C", -1.0, struct{}{}) // Negative weight
builder.AddEdge("B", "D", 5.0, struct{}{})
builder.AddEdge("C", "D", 8.0, struct{}{})
builder.AddEdge("C", "E", 10.0, struct{}{})
builder.AddEdge("D", "E", 2.0, struct{}{})

g := builder.BuildDirected()

// Create Bellman-Ford instance
bellmanFord := graph.NewBellmanFord(g)

// Find shortest path from A to E
path := bellmanFord.FindShortestPath("A", "E")
if path != nil {
    fmt.Printf("Shortest path: %v\n", path) // Output: [A C B D E]
} else {
    fmt.Println("No path found or negative cycle detected")
}
```

#### Negative Cycle Detection

Bellman-Ford can detect negative cycles in the graph:

```go
// Check for negative cycles reachable from a specific vertex
hasNegativeCycle := bellmanFord.HasNegativeCycle("A")
if hasNegativeCycle {
    fmt.Println("Negative cycle detected!")
} else {
    fmt.Println("No negative cycle found")
}

// FindShortestPath returns nil if a negative cycle is detected
path := bellmanFord.FindShortestPath("A", "E")
if path == nil {
    fmt.Println("No path found due to negative cycle")
}
```

#### Performance Characteristics
- **Time Complexity**: O(VE) where E is edges and V is vertices
- **Space Complexity**: O(V) for vertex data storage
- **Memory Efficient**: Reuses internal data structures between calls
- **Thread Safety**: Not thread-safe for concurrent calls: use separate instances of the algorithm, but the graph itself can be safely shared as long as you don't modify it
- **Negative Weights**: Supports negative edge weights (unlike Dijkstra)
- **Cycle Detection**: Can detect negative cycles in the graph

### Advanced Features

#### Cost Amplification

All algorithms support cost amplification functions that can modify edge costs at runtime or disable edges:

```go
// Define a cost function that doubles the cost of certain edges
costAmplifier := func(vertex *graph.Vertex[string, float64], edge *graph.Edge[string, float64]) (float64, bool) {
    // Double the cost for edges with specific data
    if edge.GetData().IsHighway {
        return edge.GetCost() * 2, true
    }
    return edge.GetCost(), true
}

// Apply to Dijkstra
dijkstra.Amplifier = costAmplifier

// Apply to A*
astar.Amplifier = costAmplifier

// Apply to Bellman-Ford
bellmanFord.Amplifier = costAmplifier
```

#### Thread Safety

All algorithms are **not thread-safe** for concurrent calls, but the graph itself can be safely shared:

```go
// Safe: Multiple algorithms can use the same graph
graph1 := builder.BuildDirected()

dijkstra1 := graph.NewDijkstra(graph1)
dijkstra2 := graph.NewDijkstra(graph1)
bellmanFord1 := graph.NewBellmanFord(graph1)
bellmanFord2 := graph.NewBellmanFord(graph1)

// Safe: Different goroutines can use different algorithm instances
go func() {
    path1 := dijkstra1.FindShortestPath("A", "B")
    // Process path1
}()

go func() {
    path2 := bellmanFord1.FindShortestPath("C", "D")
    // Process path2
}()
```

### Algorithm Comparison

| Feature | Dijkstra's Algorithm | A* Algorithm | Bellman-Ford Algorithm |
|---------|---------------------|--------------|------------------------|
| **Optimality** | Always optimal | Optimal with admissible heuristics | Always optimal |
| **Performance** | O(E log V) | O(E log V), often faster with good heuristics | O(VE) |
| **Memory Usage** | O(V) | O(V) | O(V) |
| **Negative Weights** | ❌ No | ❌ No | ✅ Yes |
| **Negative Cycles** | ❌ Fails | ❌ Fails | ✅ Detects |
| **Use Case** | General shortest path | Pathfinding with spatial awareness | Graphs with negative weights |
| **Heuristic Required** | No | Yes (admissible) | No |
| **Best For** | Non-negative weights | Known goal location, spatial problems | Negative weights, cycle detection |

## Performance Characteristics

- **Memory Usage**: Optimized for large graphs with millions of vertices/edges
- **Vertex Lookup**: O(1) time complexity using hash maps
- **Edge Traversal**: O(1) per edge with little to no memory allocations
- **Graph Construction**: Bulk allocation reduces memory fragmentation
- **Thread Safety**: Immutable graph structure allows safe concurrent access

## TODO

- `Graph.Clone()`, `Vertex.Clone()`, etc.
- `Builder.BuildBiDirected()`
- `MakeBuilderFromGraph()` - use visitors to collect edges and vertices
- Algorithms:
  - `ConnectedComponents`

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
