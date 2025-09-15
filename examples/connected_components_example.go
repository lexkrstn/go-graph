package main

import (
	"fmt"

	"github.com/lexkrstn/go-graph"
)

func main() {
	// Create a graph with multiple connected components
	builder := &graph.Builder[string, float64, string, string]{}

	// Component 1: A-B-C
	builder.AddVertex("A", "Node A")
	builder.AddVertex("B", "Node B")
	builder.AddVertex("C", "Node C")
	builder.AddEdge("A", "B", 1.0, "A-B")
	builder.AddEdge("B", "C", 1.0, "B-C")

	// Component 2: D-E
	builder.AddVertex("D", "Node D")
	builder.AddVertex("E", "Node E")
	builder.AddEdge("D", "E", 2.0, "D-E")

	// Component 3: F (isolated)
	builder.AddVertex("F", "Node F")

	// Build the graph
	g := builder.BuildDirected()

	// Find all connected components
	cc := graph.FindConnectedComponents(g)
	components := cc.GetComponents()
	fmt.Printf("Found %d connected components:\n", len(components))
	for i, component := range components {
		fmt.Printf("Component %d: %v\n", i+1, component)
	}

	// Check if the graph is connected
	if cc.IsConnected() {
		fmt.Println("The graph is connected")
	} else {
		fmt.Println("The graph is not connected")
	}

	// Get component count
	fmt.Printf("Component count: %d\n", cc.GetComponentCount())

	// Find the component containing a specific vertex
	componentA := cc.GetComponentForVertex("A")
	fmt.Printf("Component containing 'A': %v\n", componentA)

	componentF := cc.GetComponentForVertex("F")
	fmt.Printf("Component containing 'F': %v\n", componentF)
}
