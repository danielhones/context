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
func visitAllNodes(cur *sitter.TreeCursor, f func(n *sitter.Node, h History), hist History) {
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

func processFile(path string, search Search, opts Options) {
	parser := sitter.NewParser()
	parser.SetLanguage(python.GetLanguage())

	contents, err := ioutil.ReadFile(path)

	srcLines := bytes.Split(contents, []byte("\n"))

	// Get list of lines in the source file that have what we're
	// looking for - either string match or specific line numbers:
	matchingLines := make(map[int]struct{})
	if search.IsRegexMatch() {
		for i, s := range srcLines {
			// TODO: Use a regex match here:
			if strings.Contains(string(s), search.Val) {
				matchingLines[i] = struct{}{}
			}
		}
	} else {
		for _, x := range search.ValInts {
			matchingLines[x] = struct{}{}
		}
	}

	// The matcher just checks whether the current Node starts on
	// one of the matching lines
	matcher := func(n *sitter.Node) bool {
		_, found := matchingLines[int(n.StartPoint().Row)]
		return found
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %v: %v\n", path, err)
		return
	}

	tree := parser.Parse(nil, contents)
	lines := context(contents, tree, matcher)
	tree.Close()

	if len(lines) == 0 {
		return
	}

	fmt.Printf("\n%s\n\n", path)

	// We want line numbers to be right justified, with leading space before them,
	// so get the width of the longest number, and use that to build a format string
	// like "%4d" that we use when printing the line number:
	maxNumWidth := len(fmt.Sprintf("%d", lines[len(lines)-1]))
	fs := fmt.Sprintf("%%%dd:", maxNumWidth)

	for _, x := range lines {
		if opts.PrintNums {
			fmt.Printf(fs, x)
		}
		fmt.Println(string(srcLines[x]))
	}

	fmt.Println()
}

func usage() {
	w := flag.CommandLine.Output()
	fmt.Fprintf(
		w,
		`Usage: context [options] <search> [file1 file2 ...]

Find lines in a source code file and print the lines in the syntax tree 
leading up to them.  The "search" command line argument is required, and
so is at least one file argument.  By default, the search value is read
as an integer and searches for a line number

Options
`,
	)
	flag.PrintDefaults()
}

func main() {
	opts := Options{}

	flag.Usage = usage
	matchRegex := flag.Bool("e", false, "Search by regex instead of line number")
	flag.BoolVar(&opts.PrintNums, "n", false, "Include line numbers in output")
	flag.StringVar(&opts.Language,
		"l",
		"",
		"Specify the language to parse.  Omitting this will autodetect language from file extension.",
	)
	flag.Parse()

	if len(flag.Args()) < 2 {
		flag.Usage()
		os.Exit(2)
	}

	searchArg := flag.Args()[0]
	files := flag.Args()[1:]

	search := NewSearch()

	if *matchRegex {
		search.Val = searchArg
		search.SetRegexMatch()
	} else {
		// TODO: Accept and parse comma-separated list of integers, and maybe ranges
		//	 like 23-30:
		searchInt, err := strconv.Atoi(searchArg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse %s as int\n", searchArg)
			os.Exit(1)
		}
		search.ValInts = []int{searchInt}
	}

	for _, path := range files {
		processFile(path, search, opts)
	}
}
