package main

import (
	"fmt"
	"time"
)

type Node struct {
	nodeID      string
	address     string
	peers       []string
	counter     *GCounter
	partitioned bool
	stopChan    chan bool
}

func NewNode(nodeID, address string, peers []string) *Node {
	return &Node{
		nodeID:      nodeID,
		address:     address,
		peers:       peers,
		counter:     NewGCounter(nodeID),
		partitioned: false,
		stopChan:    make(chan bool),
	}
}

func (n *Node) Start() {
	fmt.Printf("Starting node %s on %s\n", n.nodeID, n.address)

	go n.startGossipListener()

	time.Sleep(100 * time.Millisecond)

	n.startGossipTimer()
}

func (n *Node) Stop() {
	close(n.stopChan)
}

func (n *Node) Increment() {
	n.counter.Increment()
	fmt.Printf("Node %s incremented, local value: %d\n",
		n.nodeID, n.counter.Value())
}

// GetValue returns current counter value
func (n *Node) GetValue() int {
	return n.counter.Value()
}

func (n *Node) SimulatePartition(duration time.Duration) {
	fmt.Printf("Node %s entering partition for %s\n", n.nodeID, duration)
	n.partitioned = true

	time.AfterFunc(duration, func() {
		n.partitioned = false
		fmt.Printf("Node %s partition healed\n", n.nodeID)
	})
}

func (n *Node) GetDetailedState() map[string]int {
	return n.counter.GetState()
}
