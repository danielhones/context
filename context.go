package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
)

const HELP_TEXT = `Usage: context [options] <search> [file1 file2 ...]

Find lines in a source code file and print the lines in the syntax tree 
leading up to them.  The <search> argument is required.  By default, the
search value is read as an integer and searches for a line number.

There can be any number of file arguments passed.  If there are none, it
will read from stdin.  When reading from stdin, you must include a -l flag
to indicate what language to use for parsing.
`

// Special value for the processFile function to indicate we want to read from stdin:
const READ_FROM_STDIN = "-"

// Some ANSI color codes for terminal output:
const BLUE string = "\033[34m"
const GREEN string = "\033[32m"
const END_COLOR string = "\033[0m"

// Used in the Search struct to indicate match types:
const MATCH_TYPE_REGEX = "regex"
const MATCH_TYPE_LINE = "line"

// This struct stores information about options for processing files.  It's primarily
// used as an argument to processFile
type Options struct {
	PrintNums bool      // whether to print line numbers in the output
	Language  string    // string indicating which language to use when parsing
	Colorize  bool      // whether to colorize the output
	Out       io.Writer // where to write the results
	Err       io.Writer // where to write errors
	In        io.Reader // where to read input
}

func (o Options) AutoDetect() bool {
	return o.Language == ""
}

// This struct holds information about what to look for in the source code, and what type
// of match it is.  Primarily used as an argument to processFile
type Search struct {
	MatchType string
	// If a regex/string match, this is the value to search for.  If not a regex/string
	// match, then this value is empty and meaningless:
	Val string
	// If a line match, these are the line numbers to search for.  If not a line match,
	// then this value is empty and meaningless:
	ValInts []int
}

func (s Search) IsRegexMatch() bool {
	return s.MatchType == MATCH_TYPE_REGEX
}

func (s *Search) SetRegexMatch() {
	(*s).MatchType = MATCH_TYPE_REGEX
}

func NewSearch() Search {
	return Search{MatchType: MATCH_TYPE_LINE}
}

// Visit every node in the tree, in a depth-first left-to-right traversal. Call the f function
// on each node.
func visitAllNodes(cur *sitter.TreeCursor, f func(n *sitter.Node, h History), hist History) {
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

// Read a source code file at the given path and print any matching lines to opts.Out.
// If path == "-", then this will read from opts.In instead.
func processFile(path string, search Search, opts Options) error {
	var contents []byte
	var err error

	if path == READ_FROM_STDIN {
		contents, err = ioutil.ReadAll(opts.In)
	} else {
		contents, err = ioutil.ReadFile(path)
	}
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
			// If x is not < len(srcLines), it's looking for a line number that's
			// not in the file, so we can omit it
			if x < len(srcLines) {
				matchingLines[x] = struct{}{}
			}
		}
	}

	if len(matchingLines) == 0 {
		// No need to parse anything, nothing will match
		return nil
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

	if path == READ_FROM_STDIN {
		fmt.Fprintf(opts.Out, "\n<stdin>\n\n")
	} else {
		fmt.Fprintf(opts.Out, "\n%s\n\n", path)
	}

	// We want line numbers to be right justified, with leading space before them.
	// So we get the width of the longest number first, and use that to build a
	// format string like "%4d" that we use when printing the line number:
	maxNumWidth := len(fmt.Sprintf("%d", lines[len(lines)-1]))
	fs := fmt.Sprintf("%%%dd:", maxNumWidth)

	for _, x := range lines {
		if !opts.Colorize {
			if opts.PrintNums {
				fmt.Fprintf(opts.Out, fs, x+1)
			}
			fmt.Fprintln(opts.Out, string(srcLines[x]))
			continue
		}

		if opts.PrintNums {
			fmt.Fprintf(opts.Out, BLUE+fs+END_COLOR, x+1)
		}

		_, wasMatch := matchingLines[x]
		if wasMatch {
			var s string
			if search.IsRegexMatch() {
				// TODO: This will need updating once we're using actual regex matching:
				// Colorize only the part of the line that matches the search:
				s = strings.Replace(string(srcLines[x]), search.Val, GREEN+search.Val+END_COLOR, -1)
			} else {
				// Colorize the entire matching line:
				s = GREEN + string(srcLines[x]) + END_COLOR
			}
			fmt.Fprintln(opts.Out, s)
		} else {
			fmt.Fprintln(opts.Out, string(srcLines[x]))
		}
	}

	fmt.Fprintln(opts.Out)
	return nil
}

func usage(fs *flag.FlagSet) func() {
	// Return a function so we can write to specified output and print the right defaults:
	return func() {
		w := fs.Output()
		fmt.Fprintf(w, HELP_TEXT)
		fmt.Fprintf(w, "\nOptions:\n")
		fs.PrintDefaults()

		fmt.Fprintf(w, "\nLanguages:\n")
		for _, v := range LANGUAGES {
			for _, ext := range v.Exts {
				fmt.Fprintf(w, "  %v\t%v\n", ext, v.Name)
			}
		}
	}
}

// This is effectively the main() function but with stdOut, StdErr, and command-line
// args passed in as arguments to make testing easier. The return value is the integer
// that main should exit with
func run(stdOut io.Writer, stdErr io.Writer, stdIn io.Reader, args []string) int {
	opts := Options{Out: stdOut, Err: stdErr, In: stdIn}

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.SetOutput(stdErr)
	fs.Usage = usage(fs)
	// TODO: This will eventually be a regex match, but currently is just string.Contains()
	//	 When regex is implemented, update this help message:
	matchRegex := fs.Bool("e", false, "Search by string instead of line number")
	fs.BoolVar(&opts.PrintNums, "n", false, "Include line numbers in output")
	fs.BoolVar(&opts.Colorize, "c", false, "Colorize line numbers and matches in output")
	fs.StringVar(&opts.Language,
		"l",
		"",
		"Language to parse, in shorthand/file extension form.  See below for list.  Omitting this will detect language from filename.",
	)

	err := fs.Parse(args)
	if err == flag.ErrHelp {
		// Help message was requested with the -h flag
		return 0
	} else if err != nil {
		// Usage was already printed here by FlagSet.failf, so just return:
		return 2
	} else if len(fs.Args()) < 1 {
		// The 1 remaining required arg is the search value, hence the check for len < 1.
		fs.Usage()
		return 2
	}

	// Make sure language, if specified, is valid:
	if !opts.AutoDetect() && !LangIsSupported(opts.Language) {
		fmt.Fprintf(stdErr, "Language is not supported: %v\n", opts.Language)
		return 1
	}

	searchArg := fs.Args()[0]
	files := fs.Args()[1:]

	search := NewSearch()
	if *matchRegex {
		search.Val = searchArg
		search.SetRegexMatch()
	} else {
		// TODO: Accept and parse comma-separated list of integers, and maybe ranges
		//	 like 23-30:
		searchInt, err := strconv.Atoi(searchArg)
		if err != nil {
			fmt.Fprintf(stdErr, "Could not parse %s as int\n", searchArg)
			return 1
		}
		searchInt = searchInt - 1 // convert to the zero-based index we need
		search.ValInts = []int{searchInt}
	}

	if len(files) == 0 {
		files = []string{READ_FROM_STDIN}
	}

	for _, path := range files {
		err := processFile(path, search, opts)
		if err != nil {
			fmt.Fprintf(stdErr, "Error processing %v: %v\n", path, err)
		}
	}

	return 0
}

func main() {
	x := run(os.Stdout, os.Stderr, os.Stdin, os.Args[1:])
	os.Exit(x)
}
