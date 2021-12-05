package main

import (
	"testing"
)

func TestPopFromEmpty(t *testing.T) {
	h := History{}
	_, found := h.Pop()
	if found != false {
		t.Fatalf("Expected false, got %v", found)
	}
}
