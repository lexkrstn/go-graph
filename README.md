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
  - [Connected Components Algorithm](#connected-components-algorithm)
    - [Basic Usage](#basic-usage-3)
    - [Performance Characteristics](#performance-characteristics-3)
  - [Depth-First Search (DFS) Algorithm](#depth-first-search-dfs-algorithm)
    - [Basic Usage](#basic-usage-4)
    - [Advanced Usage with Callbacks](#advanced-usage-with-callbacks)
    - [Cycle Detection](#cycle-detection)
    - [Use Cases](#use-cases)
    - [Performance Characteristics](#performance-characteristics-4)
    - [Algorithm Comparison](#algorithm-comparison)
  - [Advanced Features](#advanced-features)
    - [Cost Amplification](#cost-amplification)
    - [Thread Safety](#thread-safety)
- [Performance Characteristics](#performance-characteristics-2)
- [TODO](#todo)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Minimal Allocations**: Algorithms reuse memory between operations, and the Builder optimizes large graph loading.
- **Heap-based Priority Queues**: Optimized search algorithms with efficient priority queue implementations.
- **Graph Analysis**: Connected components algorithm for understanding graph connectivity and structure.
- **Flexible Traversal**: Depth-First Search (DFS) algorithm with callback support for custom vertex and edge processing.
- **Efficient Custom Data**: Custom data can be attached to vertices and edges and can be accessed with O(1) time complexity during runtime. The implementation relies on the [SparseSet data structure](https://
medium.com/gitconnected/
fast-ecs-from-scratch-in-rust-for-your-game-engine-d7
de8f23cd4a#a6b8) to efficiently store the associated data, making data access as fast as reading from a slice by a known index (which is much faster than using a map, particularly for large graphs).
- **Runtime Weight Computation**: Edge weights can be overwritten dynamically based on vertex/edge data. Moreover some transitions can be completely disabled at runtime, which is useful for example when you want to simulate a traffic jam or any other dynamic situation without having to rebuild the graph.
- **Thread Safety**: Static graph data is immutable and can be safely shared between goroutines. Separating the algorithms' data from the graph data saves a lot of memory, especially in case of large graphs.
- **Generic Design**: Type-safe implementation with flexible ID, cost, vertex and edge data types.

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

    // Create Dijkstra instance
    dijkstra := graph.NewDijkstra(g)

    // Find shortest path from A to C
    path := dijkstra.FindShortestPath("A", "C")
    if path != nil {
        fmt.Printf("Shortest path: %v\n", path) // Output: [A B C]
    } else {
        fmt.Println("No path found")
    }
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

The library provides four powerful pathfinding algorithms and one graph analysis algorithm, all optimized for performance and memory efficiency.

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

### Connected Components Algorithm

The Connected Components algorithm finds all groups of vertices that are reachable from each other in a graph. It's essential for understanding graph connectivity and identifying isolated subgraphs.

#### Basic Usage

```go
// Create a graph with multiple connected components
builder := &graph.Builder[string, float64, struct{}, struct{}]{}
// Component 1: A-B-C
builder.AddEdge("A", "B", 1.0, struct{}{})
builder.AddEdge("B", "C", 1.0, struct{}{})
// Component 2: D-E
builder.AddEdge("D", "E", 2.0, struct{}{})
// Component 3: F (isolated)
builder.AddVertex("F", struct{}{})

g := builder.BuildDirected()

// Find all connected components
cc := graph.FindConnectedComponents(g)

// Get all components
components := cc.GetComponents()
fmt.Printf("Found %d components:\n", len(components))
for i, component := range components {
    fmt.Printf("Component %d: %v\n", i+1, component)
}
// Output:
// Found 3 components:
// Component 1: [A B C]
// Component 2: [D E]
// Component 3: [F]

// Check if graph is connected
if cc.IsConnected() {
    fmt.Println("Graph is connected")
} else {
    fmt.Println("Graph is not connected")
}

// Get component count
count := cc.GetComponentCount()
fmt.Printf("Component count: %d\n", count) // Output: 3

// Find component containing a specific vertex
componentA := cc.GetComponentForVertex("A")
fmt.Printf("Component containing 'A': %v\n", componentA) // Output: [A B C]

componentF := cc.GetComponentForVertex("F")
fmt.Printf("Component containing 'F': %v\n", componentF) // Output: [F]
```

#### Performance Characteristics
- **Time Complexity**: O(V + E) where V is vertices and E is edges (computed once)
- **Space Complexity**: O(V) for vertex data storage
- **Query Performance**: O(1) for most operations after initial computation
- **Memory Efficient**: Components are computed once and cached for fast subsequent queries
- **Thread Safety**: Not thread-safe for concurrent calls: use separate instances of the algorithm, but the graph itself can be safely shared as long as you don't modify it
- **Directed Graphs**: Handles directed graphs by considering both incoming and outgoing edges

### Depth-First Search (DFS) Algorithm

The Depth-First Search algorithm provides flexible graph traversal capabilities with support for pathfinding, reachability checking, and custom vertex/edge processing. It uses an iterative approach to avoid stack overflow issues with large graphs.

#### Basic Usage

```go
// Create a graph
builder := &graph.Builder[string, float64, struct{}, struct{}]{}
builder.AddEdge("A", "B", 1.0, struct{}{})
builder.AddEdge("B", "C", 2.0, struct{}{})
builder.AddEdge("A", "D", 3.0, struct{}{})
builder.AddEdge("D", "E", 4.0, struct{}{})

g := builder.BuildDirected()

// Create DFS instance
dfs := graph.NewDFS(g)

// Find a path between two vertices
path := dfs.FindPath("A", "E")
if path != nil {
    fmt.Printf("Path from A to E: %v\n", path) // Output: [A D E]
} else {
    fmt.Println("No path found")
}

// Check if one vertex is reachable from another
if dfs.IsReachable("A", "C") {
    fmt.Println("C is reachable from A")
} else {
    fmt.Println("C is not reachable from A")
}

// Get all vertices reachable from a starting vertex
reachable := dfs.GetAllReachable("A")
fmt.Printf("All reachable from A: %v\n", reachable) // Output: [A B C D E]

// Check if the graph contains any cycles
if dfs.HasCycle() {
    fmt.Println("Graph contains a cycle")
} else {
    fmt.Println("Graph is acyclic")
}

// Find all cycles in the graph
cycles := dfs.FindCycles()
if len(cycles) > 0 {
    fmt.Printf("Found %d cycles:\n", len(cycles))
    for i, cycle := range cycles {
        fmt.Printf("Cycle %d: %v\n", i+1, cycle)
    }
} else {
    fmt.Println("No cycles found")
}
```

#### Advanced Usage with Callbacks

The DFS algorithm supports callback-based traversal for custom processing:

```go
// Traverse with custom callback for each vertex and edge
var visitedVertices []string
var edgeCosts []float64

dfs.TraverseFrom("A", func(vertex *graph.Vertex[string, float64], edge *graph.Edge[string, float64]) {
    visitedVertices = append(visitedVertices, vertex.GetId())
    
    if edge != nil {
        edgeCosts = append(edgeCosts, edge.GetCost())
        fmt.Printf("Visited %s via edge with cost %.1f\n", vertex.GetId(), edge.GetCost())
    } else {
        fmt.Printf("Starting from %s\n", vertex.GetId())
    }
})

fmt.Printf("Visited vertices: %v\n", visitedVertices)
fmt.Printf("Edge costs: %v\n", edgeCosts)
```

#### Cycle Detection

The DFS algorithm includes built-in cycle detection capabilities:

```go
// Create a graph with a cycle
builder := &graph.Builder[string, float64, struct{}, struct{}]{}
builder.AddEdge("A", "B", 1.0, struct{}{})
builder.AddEdge("B", "C", 2.0, struct{}{})
builder.AddEdge("C", "A", 3.0, struct{}{}) // Creates a cycle

g := builder.BuildDirected()
dfs := graph.NewDFS(g)

// Detect cycles
if dfs.HasCycle() {
    fmt.Println("Graph contains a cycle")
} else {
    fmt.Println("Graph is acyclic (DAG)")
}

// Find all cycles for detailed analysis
cycles := dfs.FindCycles()
if len(cycles) > 0 {
    fmt.Printf("Found %d cycles:\n", len(cycles))
    for i, cycle := range cycles {
        fmt.Printf("Cycle %d: %v\n", i+1, cycle)
    }
}

// Example: Validate DAG before topological sorting
if !dfs.HasCycle() {
    fmt.Println("Safe to perform topological sort")
} else {
    fmt.Println("Cannot perform topological sort - cycle detected")
    fmt.Printf("Cycles found: %v\n", cycles)
}
```

**Cycle Detection Features:**
- **HasCycle()**: Quick boolean check for cycle existence
- **FindCycles()**: Returns detailed cycle information as vertex ID arrays
- **Directed Graphs**: Detects directed cycles (A→B→C→A)
- **Self-Loops**: Detects self-loops (A→A)
- **Multiple Components**: Finds cycles in any connected component
- **Duplicate Prevention**: Avoids reporting the same cycle multiple times
- **Performance**: O(V + E) time complexity
- **DAG Validation**: Perfect for validating Directed Acyclic Graphs

#### Use Cases

- **Path Finding**: Find any path between two vertices (not necessarily shortest)
- **Reachability Analysis**: Check if one vertex can reach another
- **Graph Traversal**: Visit all reachable vertices from a starting point
- **Custom Processing**: Process vertices and edges during traversal
- **Cycle Detection**: Detect cycles in directed graphs using `HasCycle()` and `FindCycles()`
- **DAG Validation**: Verify that a graph is a Directed Acyclic Graph (DAG)
- **Cycle Analysis**: Get detailed information about all cycles in the graph
- **Component Analysis**: Find all vertices in a connected component

#### Performance Characteristics

- **Time Complexity**: O(V + E) where V is vertices and E is edges
- **Space Complexity**: O(V) for vertex data storage
- **Iterative Implementation**: Uses explicit stack to avoid recursion and stack overflow
- **Memory Efficient**: Reuses internal data structures between calls
- **Thread Safety**: Not thread-safe for concurrent calls: use separate instances of the algorithm, but the graph itself can be safely shared as long as you don't modify it
- **Large Graph Support**: Can handle very deep graphs without stack overflow issues

#### Algorithm Comparison

| Algorithm | Use Case | Time Complexity | Space Complexity | Optimal Path |
|-----------|----------|----------------|------------------|--------------|
| **DFS** | Path finding, traversal, reachability | O(V + E) | O(V) | No |
| **Dijkstra** | Shortest path (non-negative weights) | O(E log V) | O(V) | Yes |
| **A*** | Shortest path with heuristics | O(E log V) | O(V) | Yes |
| **Bellman-Ford** | Shortest path (negative weights) | O(VE) | O(V) | Yes |

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

### Graph Analysis Algorithms

| Feature | Connected Components Algorithm |
|---------|-------------------------------|
| **Purpose** | Find groups of connected vertices |
| **Performance** | O(V + E) (computed once) |
| **Memory Usage** | O(V) |
| **Query Performance** | O(1) for most operations |
| **Directed Graphs** | ✅ Yes (considers both directions) |
| **Use Case** | Graph connectivity analysis, subgraph identification |
| **Best For** | Understanding graph structure, finding isolated components |

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
