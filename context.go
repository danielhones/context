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
)

// Visit every node in the tree, in a depth-first left-to-right traversal
func visitAllNodes(cur *sitter.TreeCursor, f func(n *sitter.Node, h History), hist History) {
	// hist.Push(cur.CurrentNode().StartPoint().Row)
	hist.Push(cur.CurrentNode())

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
func context(src []byte, tree *sitter.Tree, matcher func(n *sitter.Node) bool, lang LanguageInfo) []int {
	lineMap := make(map[uint32]struct{}) // Use a map like a set to track the line numbers we accumulate

	cur := sitter.NewTreeCursor(tree.RootNode())
	defer cur.Close()

	// The visitor function does most of the work here.  It's called once on every node
	// in the tree.  It calls the matcher function on the node, and if it returns true
	// then it adds the current node line number and all the line numbers in the history
	// to the lineMap.  Then, for each node in the history, it looks back for any other
	// odes it needs to store as well.  This is how we accumulate all the if/elifs or
	// switch/cases leading up to a final matching line.
	visitor := func(n *sitter.Node, h History) {
		if !matcher(n) {
			return
		}

		// If the node matched, then we add it and the current history to
		// lineMap.  For any "multi-branch" nodes, eg else blocks or elif blocks,
		// we add their previous siblings as well.  This way we can see the
		// full set of conditons that lead to a given line executing
		lineMap[n.StartPoint().Row] = struct{}{}
		for _, x := range h.Nodes() {
			lineMap[x.StartPoint().Row] = struct{}{}
			// fmt.Println("History node:", x.StartPoint().Row, x.Type())
			if IsMultiBranchNode(x, lang) {
				// Add line numbers for all previous siblings of this node
				prev := x.PrevSibling()
				for {
					if prev == nil {
						break
					}
					// fmt.Println("PrevSib node:", prev.StartPoint().Row, prev.Type())
					if IsMultiBranchNode(prev, lang) {
						lineMap[prev.StartPoint().Row] = struct{}{}
					}
					prev = prev.PrevSibling()
				}
			}
		}
	}

	visitAllNodes(cur, visitor, History{})

	// Turn our "set" of line numbers into a list, sort it, then return it:
	lines := make([]int, len(lineMap))
	i := 0
	for k := range lineMap {
		lines[i] = int(k)
		i++
	}
	sort.Ints(lines)
	return lines
}

func processFile(path string, search Search, opts Options) error {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

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

	// Then the matcher just needs to check whether the current Node starts
	// on one of the matching lines we found:
	matcher := func(n *sitter.Node) bool {
		_, found := matchingLines[int(n.StartPoint().Row)]
		return found
	}

	var langInfo LanguageInfo
	if opts.AutoDetect() {
		langInfo, err = LangFromFilename(path)
	} else {
		langInfo, err = LangFromString(opts.Language)
	}
	if err != nil {
		return err
	}

	parser := sitter.NewParser()
	parser.SetLanguage(langInfo.Lang)
	tree := parser.Parse(nil, contents)
	lines := context(contents, tree, matcher, langInfo)
	tree.Close()

	if len(lines) == 0 {
		return nil
	}

	fmt.Fprintf(os.Stdout, "\n%s\n\n", path)

	// We want line numbers to be right justified, with leading space before them,
	// so get the width of the longest number, and use that to build a format string
	// like "%4d" that we use when printing the line number:
	maxNumWidth := len(fmt.Sprintf("%d", lines[len(lines)-1]))
	fs := fmt.Sprintf("%%%dd:", maxNumWidth)

	for _, x := range lines {
		if opts.PrintNums {
			fmt.Fprintf(opts.Out, fs, x+1)
		}
		fmt.Fprintln(opts.Out, string(srcLines[x]))
	}

	fmt.Fprintln(opts.Out)
	return nil
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

Options:
`,
	)
	flag.PrintDefaults()

	fmt.Fprintf(w, "\nLanguages:\n")
	for _, v := range LANGUAGES {
		for _, ext := range v.Exts {
			fmt.Fprintf(w, "  %v\t%v\n", ext, v.Name)
		}
	}
}

func main() {
	opts := NewOptions()

	flag.Usage = usage
	// TODO: This will eventually be a regex match, but currently is just string.Contains()
	//	 When regex is implemented, update this help message:
	matchRegex := flag.Bool("e", false, "Search by string instead of line number")
	flag.BoolVar(&opts.PrintNums, "n", false, "Include line numbers in output")
	flag.StringVar(&opts.Language,
		"l",
		"",
		"Language to parse, in shorthand/file extension form.  See below for list.  Omitting this will detect language from filename.",
	)
	flag.Parse()

	if len(flag.Args()) < 2 {
		flag.Usage()
		os.Exit(2)
	}

	// Make sure language, if specified, is valid:
	if !opts.AutoDetect() && !LangIsSupported(opts.Language) {
		fmt.Fprintf(os.Stderr, "Language is not supported: %v\n", opts.Language)
		os.Exit(1)
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
		searchInt = searchInt - 1 // convert to the zero-based index we need
		search.ValInts = []int{searchInt}
	}

	for _, path := range files {
		err := processFile(path, search, opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %v: %v\n", path, err)
		}
	}
}
