package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("🚀 Starting CRDT Counter Demo")
	fmt.Println("===============================")

	// Create 3 nodes
	node1 := NewNode("alice", ":8001", []string{":8002", ":8003"})
	node2 := NewNode("bob", ":8002", []string{":8001", ":8003"})
	node3 := NewNode("charlie", ":8003", []string{":8001", ":8002"})

	// Start all nodes
	go node1.Start()
	go node2.Start()
	go node3.Start()

	// Wait for nodes to start up
	time.Sleep(500 * time.Millisecond)

	fmt.Println("\n📊 Initial state:")
	printStatus(node1, node2, node3)

	// Phase 1: Concurrent increments
	fmt.Println("\n🔄 Phase 1: Concurrent increments")
	go simulateClicks(node1, 3)
	go simulateClicks(node2, 2)
	go simulateClicks(node3, 4)

	// Wait for gossip to spread
	time.Sleep(3 * time.Second)
	fmt.Println("\nAfter gossip convergence:")
	printStatus(node1, node2, node3)

	// Phase 2: Network partition
	fmt.Println("\n🔌 Phase 2: Network partition (isolating Alice)")
	node1.SimulatePartition(5 * time.Second)

	// Increments during partition
	go simulateClicks(node1, 3) // Alice isolated
	go simulateClicks(node2, 2) // Bob and Charlie connected
	go simulateClicks(node3, 2)

	time.Sleep(3 * time.Second)
	fmt.Println("\nDuring partition:")
	printStatus(node1, node2, node3)

	// Wait for partition to heal
	time.Sleep(3 * time.Second)
	// After partition healed:
	fmt.Println("\nAfter partition healed:")
	printStatus(node1, node2, node3)

	// Add this:
	fmt.Println("\n⏳ Waiting for full convergence...")
	time.Sleep(4 * time.Second) // Give Alice time to receive full gossip
	fmt.Println("Final state:")
	printStatus(node1, node2, node3)

	// Phase 3: Show detailed state
	fmt.Println("\n🔍 Detailed counter state:")
	fmt.Printf("Alice's view: %v\n", node1.GetDetailedState())
	fmt.Printf("Bob's view:   %v\n", node2.GetDetailedState())
	fmt.Printf("Charlie's view: %v\n", node3.GetDetailedState())

	fmt.Println("\n✅ Demo complete! All nodes converged.")
}

// simulateClicks increments the counter multiple times with random delays
func simulateClicks(node *Node, count int) {
	for i := 0; i < count; i++ {
		node.Increment()
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
}

// printStatus shows current value for all nodes
func printStatus(node1, node2, node3 *Node) {
	fmt.Printf("Alice: %d, Bob: %d, Charlie: %d\n",
		node1.GetValue(),
		node2.GetValue(),
		node3.GetValue())
}
