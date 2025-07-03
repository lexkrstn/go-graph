package graph

import (
	"testing"
)

func TestIntegration(t *testing.T) {
	t.Run("Social network graph", func(t *testing.T) {
		// Create a social network graph where vertices represent users
		// and edges represent friendships
		builder := &Builder[int, int, User, Friendship]{}

		// Add users
		builder.AddVertex(1, User{Name: "Alice", Age: 25})
		builder.AddVertex(2, User{Name: "Bob", Age: 30})
		builder.AddVertex(3, User{Name: "Charlie", Age: 28})
		builder.AddVertex(4, User{Name: "Diana", Age: 22})

		// Add friendships (bidirectional)
		builder.AddBiEdge(1, 2, 1, Friendship{Since: "2020-01-01", Strength: 8})
		builder.AddBiEdge(1, 3, 1, Friendship{Since: "2019-06-15", Strength: 6})
		builder.AddBiEdge(2, 3, 1, Friendship{Since: "2021-03-10", Strength: 9})
		builder.AddBiEdge(3, 4, 1, Friendship{Since: "2022-01-20", Strength: 7})

		graph := builder.BuildDirected()

		// Verify graph structure
		if graph.GetVertexCount() != 4 {
			t.Errorf("Expected 4 users, got %d", graph.GetVertexCount())
		}

		if graph.GetEdgeCount() != 8 { // 4 bidirectional edges = 8 directed edges
			t.Errorf("Expected 8 friendships, got %d", graph.GetEdgeCount())
		}

		// Test vertex data retrieval
		alice, err := graph.GetVertexById(1)
		if err != nil {
			t.Errorf("Failed to get Alice: %v", err)
		}

		aliceData, err := graph.GetVertexData(alice)
		if err != nil {
			t.Errorf("Failed to get Alice's data: %v", err)
		}

		if aliceData.Name != "Alice" {
			t.Errorf("Expected Alice's name 'Alice', got %s", aliceData.Name)
		}

		// Test edge traversal
		friendCount := 0
		graph.VisitEdges(func(vertex *Vertex[int, int], edge *Edge[int, int]) {
			if vertex.GetId() == 1 { // Alice's friends
				friendCount++
				friendship, _ := graph.GetEdgeData(edge)
				if friendship.Strength < 1 || friendship.Strength > 10 {
					t.Errorf("Invalid friendship strength: %d", friendship.Strength)
				}
			}
		})

		if friendCount != 2 { // Alice has 2 friends (Bob and Charlie)
			t.Errorf("Expected Alice to have 2 friends, got %d", friendCount)
		}

		// Test predicate functions
		hasStrongFriendships := graph.SomeEdges(func(vertex *Vertex[int, int], edge *Edge[int, int]) bool {
			friendship, _ := graph.GetEdgeData(edge)
			return friendship.Strength >= 8
		})

		if !hasStrongFriendships {
			t.Error("Expected to find strong friendships (strength >= 8)")
		}

		allFriendshipsValid := graph.EveryEdge(func(vertex *Vertex[int, int], edge *Edge[int, int]) bool {
			friendship, _ := graph.GetEdgeData(edge)
			return friendship.Strength >= 1 && friendship.Strength <= 10
		})

		if !allFriendshipsValid {
			t.Error("Expected all friendships to have valid strength values")
		}
	})

	t.Run("Road network graph", func(t *testing.T) {
		// Create a road network where vertices represent cities
		// and edges represent roads with distances
		builder := &Builder[string, float64, City, Road]{}

		// Add cities
		builder.AddVertex("NYC", City{Name: "New York City", Population: 8400000})
		builder.AddVertex("BOS", City{Name: "Boston", Population: 675000})
		builder.AddVertex("DC", City{Name: "Washington DC", Population: 689000})
		builder.AddVertex("PHL", City{Name: "Philadelphia", Population: 1600000})

		// Add roads (some bidirectional, some one-way)
		builder.AddBiEdge("NYC", "BOS", 215.0, Road{Name: "I-95", Lanes: 4, SpeedLimit: 65})
		builder.AddBiEdge("NYC", "DC", 225.0, Road{Name: "I-95", Lanes: 6, SpeedLimit: 70})
		builder.AddEdge("NYC", "PHL", 97.0, Road{Name: "I-95", Lanes: 4, SpeedLimit: 65})
		builder.AddEdge("PHL", "DC", 140.0, Road{Name: "I-95", Lanes: 4, SpeedLimit: 65})

		graph := builder.BuildDirected()

		// Verify graph structure
		if graph.GetVertexCount() != 4 {
			t.Errorf("Expected 4 cities, got %d", graph.GetVertexCount())
		}

		if graph.GetEdgeCount() != 6 { // 2 bidirectional + 2 one-way = 6 directed edges
			t.Errorf("Expected 6 roads, got %d", graph.GetEdgeCount())
		}

		// Test city data retrieval
		nyc, err := graph.GetVertexById("NYC")
		if err != nil {
			t.Errorf("Failed to get NYC: %v", err)
		}

		nycData, err := graph.GetVertexData(nyc)
		if err != nil {
			t.Errorf("Failed to get NYC's data: %v", err)
		}

		if nycData.Population != 8400000 {
			t.Errorf("Expected NYC population 8400000, got %d", nycData.Population)
		}

		// Test road traversal
		totalDistance := 0.0
		graph.VisitEdges(func(vertex *Vertex[string, float64], edge *Edge[string, float64]) {
			totalDistance += edge.GetCost()
		})

		expectedDistance := 215.0 + 215.0 + 225.0 + 225.0 + 97.0 + 140.0
		if totalDistance != expectedDistance {
			t.Errorf("Expected total distance %.1f, got %.1f", expectedDistance, totalDistance)
		}

		// Test finding high-speed roads
		hasHighSpeedRoads := graph.SomeEdges(func(vertex *Vertex[string, float64], edge *Edge[string, float64]) bool {
			road, _ := graph.GetEdgeData(edge)
			return road.SpeedLimit >= 70
		})

		if !hasHighSpeedRoads {
			t.Error("Expected to find high-speed roads (speed limit >= 70)")
		}

		// Test finding major cities (population > 1M)
		majorCities := 0
		graph.VisitVertices(func(vertex *Vertex[string, float64]) {
			city, _ := graph.GetVertexData(vertex)
			if city.Population > 1000000 {
				majorCities++
			}
		})

		if majorCities != 2 { // NYC and PHL
			t.Errorf("Expected 2 major cities, got %d", majorCities)
		}
	})

	t.Run("Task dependency graph", func(t *testing.T) {
		// Create a task dependency graph where vertices represent tasks
		// and edges represent dependencies with completion times
		builder := &Builder[int, int, Task, Dependency]{}

		// Add tasks
		builder.AddVertex(1, Task{Name: "Design", Duration: 5, Priority: "High"})
		builder.AddVertex(2, Task{Name: "Implement", Duration: 10, Priority: "Medium"})
		builder.AddVertex(3, Task{Name: "Test", Duration: 3, Priority: "High"})
		builder.AddVertex(4, Task{Name: "Deploy", Duration: 2, Priority: "Low"})

		// Add dependencies (Design -> Implement -> Test -> Deploy)
		builder.AddEdge(1, 2, 0, Dependency{Type: "Blocks", Critical: true})
		builder.AddEdge(2, 3, 0, Dependency{Type: "Blocks", Critical: true})
		builder.AddEdge(3, 4, 0, Dependency{Type: "Blocks", Critical: false})

		graph := builder.BuildDirected()

		// Verify graph structure
		if graph.GetVertexCount() != 4 {
			t.Errorf("Expected 4 tasks, got %d", graph.GetVertexCount())
		}

		if graph.GetEdgeCount() != 3 {
			t.Errorf("Expected 3 dependencies, got %d", graph.GetEdgeCount())
		}

		// Test task data retrieval
		design, err := graph.GetVertexById(1)
		if err != nil {
			t.Errorf("Failed to get Design task: %v", err)
		}

		designData, err := graph.GetVertexData(design)
		if err != nil {
			t.Errorf("Failed to get Design task data: %v", err)
		}

		if designData.Duration != 5 {
			t.Errorf("Expected Design duration 5, got %d", designData.Duration)
		}

		// Test dependency traversal
		criticalDependencies := 0
		graph.VisitEdges(func(vertex *Vertex[int, int], edge *Edge[int, int]) {
			dependency, _ := graph.GetEdgeData(edge)
			if dependency.Critical {
				criticalDependencies++
			}
		})

		if criticalDependencies != 2 {
			t.Errorf("Expected 2 critical dependencies, got %d", criticalDependencies)
		}

		// Test finding high priority tasks
		highPriorityTasks := 0
		graph.VisitVertices(func(vertex *Vertex[int, int]) {
			task, _ := graph.GetVertexData(vertex)
			if task.Priority == "High" {
				highPriorityTasks++
			}
		})

		if highPriorityTasks != 2 { // Design and Test
			t.Errorf("Expected 2 high priority tasks, got %d", highPriorityTasks)
		}

		// Test total project duration calculation
		totalDuration := 0
		graph.VisitVertices(func(vertex *Vertex[int, int]) {
			task, _ := graph.GetVertexData(vertex)
			totalDuration += task.Duration
		})

		if totalDuration != 20 { // 5 + 10 + 3 + 2
			t.Errorf("Expected total duration 20, got %d", totalDuration)
		}
	})
}

// Test data structures for integration tests

type User struct {
	Name string
	Age  int
}

type Friendship struct {
	Since    string
	Strength int
}

type City struct {
	Name       string
	Population int
}

type Road struct {
	Name       string
	Lanes      int
	SpeedLimit int
}

type Task struct {
	Name     string
	Duration int
	Priority string
}

type Dependency struct {
	Type     string
	Critical bool
}
