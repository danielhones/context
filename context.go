package main

import (
	// "bytes"
	"fmt"
	"io/ioutil"
	"os"
	// "strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/python"
	// "github.com/smacker/go-tree-sitter/ruby"
)

func visitAllNodes(cur *sitter.TreeCursor, f func(n *sitter.Node)) {
	if cur.GoToFirstChild() {
		visitAllNodes(cur, f)
	} else {
		// The current node has no children, so we visit it now:
		f(cur.CurrentNode())

		// Then move onto each sibling:
		for {
			if cur.GoToNextSibling() {
				visitAllNodes(cur, f)
			} else if cur.GoToParent() {
				f(cur.CurrentNode())
			} else {
				// There's no sibling and no parent so we must be at the root
				return
			}
		}
	}
}

// Find matching lines and the branches in the tree that lead to them.
// Returns list of line numbers.
func context(src []byte, tree *sitter.Tree, lookFor string) []int {
	lineMap := make(map[int]struct{})

	cur := sitter.NewTreeCursor(tree.RootNode())
	defer cur.Close()

	visitAllNodes(
		cur,
		func(n *sitter.Node) {
			fmt.Println("line:", n.StartPoint().Row, n.Type())
		},
	)

	// TODO: Sort the line numbers before returning:
	lines := make([]int, len(lineMap))
	j := 0
	for k := range lineMap {
		lines[j] = k
	}
	return lines
}

func main() {
	parser := sitter.NewParser()
	parser.SetLanguage(python.GetLanguage())

	for _, path := range os.Args[1:] {
		contents, err := ioutil.ReadFile(path)

		// srcLines := bytes.Split(contents, []byte("\n"))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading %v: %v\n", path, err)
			continue
		}

		fmt.Println("File:", path)

		tree := parser.Parse(nil, contents)
		lines := context(contents, tree, "some_variable")
		tree.Close()

		fmt.Println("Lines:", lines)
	}
}
