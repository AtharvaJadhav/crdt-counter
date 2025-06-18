package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"time"
)

type GossipMessage struct {
	FromNode string         `json:"from_node"`
	Counter  map[string]int `json:"counter"`
}

func (n *Node) startGossipListener() {
	addr, err := net.ResolveUDPAddr("udp", n.address)
	if err != nil {
		fmt.Printf("Error resolving address %s: %v\n", n.address, err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Printf("Error starting UDP listener on %s: %v\n", n.address, err)
		return
	}
	defer conn.Close()

	fmt.Printf("Node %s listening on %s\n", n.nodeID, n.address)

	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Error reading UDP message: %v\n", err)
			continue
		}

		go n.handleGossipMessage(buffer[:n], addr)

	}
}

func (n *Node) handleGossipMessage(data []byte, addr *net.UDPAddr) {
	var message GossipMessage
	err := json.Unmarshal(data, &message)
	if err != nil {
		fmt.Printf("Error unmarshaling gossip message: %v\n", err)
		return
	}

	// Update local counter with received state
	n.counter.SetState(message.Counter)

	fmt.Printf("Node %s received gossip from %s\n",
		n.nodeID, message.FromNode, n.counter.Value())
}

func (n *Node) gossipToRandomPeer() {
	if len(n.peers) == 0 {
		return
	}

	peer := n.peers[rand.Intn(len(n.peers))]
	message := GossipMessage{
		FromNode: n.nodeID,
		Counter:  n.counter.GetState(),
	}

	data, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("Error marshaling gossip message: %v\n", err)
		return
	}

	n.sendUDP(peer, data)
}

func (n *Node) sendUDP(address string, data []byte) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Printf("Error resolving address %s: %v\n", address, err)
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Printf("Error connecting to peer %s: %v\n", address, err)
		return
	}
	defer conn.Close()

	_, err = conn.Write(data)
	if err != nil {
		fmt.Printf("Error sending UDP message to %s: %v\n", address, err)
	}
}

func (n *Node) startGossipTimer() {
	ticker := time.NewTicker(2 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				if !n.partitioned {
					n.gossipToRandomPeer()
				}
			case <-n.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}
