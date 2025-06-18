# CRDT Distributed Counter

A Go implementation of Conflict-free Replicated Data Types (CRDTs) for distributed counting.

## What This Demonstrates

- **Conflict-free merging**: Multiple nodes can increment simultaneously without conflicts
- **Eventual consistency**: All nodes converge to the same value
- **Partition tolerance**: System continues working during network splits
- **Automatic conflict resolution**: No manual intervention needed

## How to Run

```bash
go run *.go