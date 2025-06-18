package main

import (
	"encoding/json"
	"sync"
)

type GCounter struct {
	nodeID   string
	counters map[string]int
	mutex    sync.RWMutex
}

func NewGCounter(nodeID string) *GCounter {
	return &GCounter{
		nodeID:   nodeID,
		counters: make(map[string]int),
		mutex:    sync.RWMutex{},
	}
}

func (g *GCounter) Increment() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.counters[g.nodeID]++
}

func (g *GCounter) Decrement() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.counters[g.nodeID]--
}

func (g *GCounter) Value() int {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	total := 0
	for _, count := range g.counters {
		total += count
	}
	return total
}

func (g *GCounter) Merge(other *GCounter) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	for nodeID, count := range other.counters {
		g.counters[nodeID] = max(g.counters[nodeID], count)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (g *GCounter) GetState() map[string]int {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	state := make(map[string]int)
	for nodeID, count := range g.counters {
		state[nodeID] = count
	}
	return state
}

func (g *GCounter) SetState(state map[string]int) {
	other := &GCounter{counters: state}
	g.Merge(other)
}

func (g *GCounter) Marshal() ([]byte, error) {
	state := g.GetState()
	return json.Marshal(state)
}

func (g *GCounter) Unmarshal(data []byte) error {
	var state map[string]int
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}
	g.SetState(state)
	return nil
}
