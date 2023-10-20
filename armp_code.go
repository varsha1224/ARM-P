package main

import (
    "fmt"
    "math/rand"
    "time"
)

const (
    numNodes = 5
)

type message struct {
    sender  int
    payload string
}

type node struct {
    id           int
    received     map[int]bool
    nextSequence int
}

func (n *node) multicast(m message, nodes []node) {
    for i := range nodes {
        if i == n.id {
            continue
        }
        if rand.Float32() < 0.5 {
            nodes[i].deliver(m)
        }
    }
}

func (n *node) deliver(m message) {
    if !n.received[m.sender] && m.payload != "" {
        fmt.Printf("Node %d received message %q from node %d\n", n.id, m.payload, m.sender)
        n.received[m.sender] = true
        n.nextSequence++
        n.multicast(message{n.id, fmt.Sprintf("%d", n.nextSequence)}, nodes)
    }
}

var nodes []node

func main() {
    rand.Seed(time.Now().UnixNano())

    nodes = make([]node, numNodes)
    for i := range nodes {
        nodes[i].id = i
        nodes[i].received = make(map[int]bool)
    }

    initialMessage := message{-1, "initial"}
    for i := range nodes {
        nodes[i].deliver(initialMessage)
    }

    for {
        time.Sleep(time.Second)
        for i := range nodes {
            nodes[i].multicast(message{i, fmt.Sprintf("msg %d", i)}, nodes)
        }
    }
}
