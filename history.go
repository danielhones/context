package main

import (
	sitter "github.com/smacker/go-tree-sitter"
)

//
type History struct {
	nodes []*sitter.Node // List of line nodes in our current path in the tree
}

func (h *History) Push(x *sitter.Node) {
	(*h).nodes = append((*h).nodes, x)
}

// Returns (*sitter.Node, true) if h is not empty, otherwise returns (nil, false)
func (h *History) Pop() (*sitter.Node, bool) {
	length := len((*h).nodes)
	if length < 1 {
		return nil, false
	}
	x := (*h).nodes[length-1]
	(*h).nodes = (*h).nodes[:length-1]
	return x, true
}

func (h History) Nodes() []*sitter.Node {
	return h.nodes
}
