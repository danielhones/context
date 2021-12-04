package main

//
type History struct {
	lines []uint32 // List of line numbers for the nodes in our current path in the tree
}

func (h *History) Push(x uint32) {
	(*h).lines = append((*h).lines, x)
}

// Returns (uint32, true) if h is not empty, otherwise returns (0, false)
func (h *History) Pop() (uint32, bool) {
	length := len((*h).lines)
	if length < 1 {
		return 0, false
	}
	x := (*h).lines[length-1]
	(*h).lines = (*h).lines[:length-1]
	return x, true
}

func (h History) Lines() []uint32 {
	return h.lines
}
