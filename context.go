package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/python"
	// "github.com/smacker/go-tree-sitter/ruby"
)

// Visit every node in the tree, in a depth-first left-to-right traversal
func ORIGvisitAllNodes(cur *sitter.TreeCursor, f func(n *sitter.Node)) {
	if cur.GoToFirstChild() {
		ORIGvisitAllNodes(cur, f)
	} else {
		// The current node has no children, so we visit it now:
		f(cur.CurrentNode())

		// Then move onto each sibling:
		for {
			if cur.GoToNextSibling() {
				ORIGvisitAllNodes(cur, f)
			} else if cur.GoToParent() {
				f(cur.CurrentNode())
			} else {
				// There's no sibling and no parent so we must be at the root
				return
			}
		}
	}
}

// Visit every node in the tree, in a depth-first left-to-right traversal
func visitAllNodes(cur *sitter.TreeCursor, f func(n *sitter.Node, h History), hist History) {
	// added := 1
	hist.Push(cur.CurrentNode().StartPoint().Row)

	if cur.GoToFirstChild() {
		visitAllNodes(cur, f, hist)
	} else {
		// The current node has no children, so we visit it now:
		f(cur.CurrentNode(), hist)
		hist.Pop()

		// Then move onto each sibling:
		for {
			if cur.GoToNextSibling() {
				visitAllNodes(cur, f, hist)
			} else if cur.GoToParent() {
				hist.Pop()
				f(cur.CurrentNode(), hist)
			} else {
				// There's no sibling and no parent so we must be at the root
				return
			}
		}
	}
}

// Find matching lines and the branches in the tree that lead to them.
// Returns list of line numbers.
func context(src []byte, tree *sitter.Tree, matcher func(n *sitter.Node) bool) []int {
	lineMap := make(map[uint32]struct{})

	cur := sitter.NewTreeCursor(tree.RootNode())
	defer cur.Close()

	visitor := func(n *sitter.Node, h History) {
		if matcher(n) {
			lineMap[n.StartPoint().Row] = struct{}{}
			for _, x := range h.Lines() {

				lineMap[x] = struct{}{}
			}
		}
	}

	visitAllNodes(cur, visitor, History{})

	lines := make([]int, len(lineMap))
	i := 0
	for k := range lineMap {
		lines[i] = int(k)
		i++
	}
	sort.Ints(lines)
	return lines
}

func usage() {
	w := flag.CommandLine.Output()
	fmt.Fprintf(
		w,
		`Usage: context [options] <search> [file1 file2 ...]

Find lines in a source code file and print the lines in the syntax tree 
leading up to them.  The "search" command line argument is required, and
so is at least one file path.

Options
`,
	)
	flag.PrintDefaults()
}

func main() {
	parser := sitter.NewParser()
	parser.SetLanguage(python.GetLanguage())

	flag.Usage = usage
	matchRegex := flag.Bool("e", false, "Search by regex instead of line number")
	printNums := flag.Bool("n", false, "Include line numbers in output")
	flag.Parse()

	if len(flag.Args()) < 2 {
		flag.Usage()
		os.Exit(2)
	}

	lookFor := flag.Args()[0]
	files := flag.Args()[1:]

	lookForInt, err := strconv.Atoi(lookFor)

	if !*matchRegex && err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse %s as int\n", lookFor)
		os.Exit(1)
	}

	for _, path := range files {
		contents, err := ioutil.ReadFile(path)

		srcLines := bytes.Split(contents, []byte("\n"))

		// Get list of lines in the source file that have what we're
		// looking for - either string match or specific line numbers:
		matchingLines := make(map[int]struct{})
		if *matchRegex {
			for i, s := range srcLines {
				// TODO: Use a regex match here:
				if strings.Contains(string(s), lookFor) {
					matchingLines[i] = struct{}{}
				}
			}
		} else {
			matchingLines[lookForInt] = struct{}{}
		}

		// The matcher just checks whether the current Node starts on
		// one of the matching lines
		matcher := func(n *sitter.Node) bool {
			_, found := matchingLines[int(n.StartPoint().Row)]
			return found
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading %v: %v\n", path, err)
			continue
		}

		tree := parser.Parse(nil, contents)
		lines := context(contents, tree, matcher)
		tree.Close()

		if len(lines) == 0 {
			continue
		}

		// Now print the result, starting with filename if there were multiple files:
		if len(files) > 1 {
			fmt.Printf("\n%s\n\n", path)
		} else {
			fmt.Printf("\n")
		}

		// We want line numbers to be right justified, with leading space before them,
		// so get the width of the longest number, and use that to build a format string
		// like "%4d" that we use when printing the line number:
		maxNumWidth := len(fmt.Sprintf("%d", lines[len(lines)-1]))
		fs := fmt.Sprintf("%%%dd:", maxNumWidth)

		for _, x := range lines {
			if *printNums {
				fmt.Printf(fs, x)
			}
			fmt.Println(string(srcLines[x]))
		}

		fmt.Println()
	}
}