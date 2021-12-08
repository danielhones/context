package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func assertStringEqual(t *testing.T, expected string, result string, msg string) {
	if result != expected {
		t.Fatalf("%s\nActual:\n<<<%v>>>\n\nExpected:\n<<<%v>>>", msg, result, expected)
	}
}

func assertInString(t *testing.T, expected string, result string, msg string) {
	if !strings.Contains(result, expected) {
		t.Fatalf("%s\nActual:\n<<<%v>>>\n\nExpected:\n<<<%v>>>", msg, result, expected)
	}
}

func assertEqual(t *testing.T, expected interface{}, actual interface{}, msg string) {
	if expected != actual {
		t.Fatalf("%s\nActual: %v\nExpected: %v", msg, actual, expected)
	}
}

func assertNotEqual(t *testing.T, expected interface{}, actual interface{}, msg string) {
	if expected == actual {
		t.Fatalf("%s\nActual: %v\nExpected: %v", msg, actual, expected)
	}
}

func TestRunUsage(t *testing.T) {
	out := bytes.Buffer{}
	errOut := bytes.Buffer{}
	in := bytes.Buffer{}
	// This represents an invalid use, with no arguments passed.
	// Make sure it prints our custom usage message:
	exitCode := run(&out, &errOut, &in, []string{})
	assertEqual(t, exitCode, 2, "")
	expected := `Usage: context [options] <search> [file1 file2 ...]

Find lines in a source code file and print the lines in the syntax tree 
leading up to them.  The <search> argument is required, but there can be
any number of file arguments passed.  If there are no file arguments, it
will read from stdin.  By default, the search value is read as an integer
and searches for a line number.

Options:
`
	assertInString(t, expected, errOut.String(), "")

	// Check that intentionally using -h flag prints usage and exits 0:
	out.Reset()
	errOut.Reset()
	exitCode = run(&out, &errOut, &in, []string{"-h"})
	assertEqual(t, exitCode, 0, "")
	assertInString(t, expected, errOut.String(), "")
}

func TestNoMatchingLines(t *testing.T) {
	out := bytes.Buffer{}
	errOut := bytes.Buffer{}
	in := bytes.Buffer{}
	exitCode := run(&out, &errOut, &in, []string{"-n", "999999", "sample_files/sample.go"})
	assertEqual(t, exitCode, 0, "")
	assertStringEqual(t, "", out.String(), "")
}

func TestUnsupportedLanguage(t *testing.T) {
	out := bytes.Buffer{}
	errOut := bytes.Buffer{}
	in := bytes.Buffer{}
	exitCode := run(&out, &errOut, &in, []string{"-n", "-l", "fake", "foo", "bar"})
	assertEqual(t, exitCode, 1, "")
	assertStringEqual(t, "Language is not supported: fake\n", errOut.String(), "")
}

func TestNonexistentFile(t *testing.T) {
	out := bytes.Buffer{}
	errOut := bytes.Buffer{}
	in := bytes.Buffer{}

	exitCode := run(&out, &errOut, &in, []string{"-n", "12", "/not/a/real/file.py"})
	// Exit is still zero because a failure of one file in a list of many shouldn't
	// cause the entire thing to fail:
	assertEqual(t, exitCode, 0, "")
	assertInString(t, "Error processing /not/a/real/file.py", errOut.String(), "")
}

func TestUnparsableInt(t *testing.T) {
	out := bytes.Buffer{}
	errOut := bytes.Buffer{}
	in := bytes.Buffer{}

	exitCode := run(&out, &errOut, &in, []string{"-n", "noninteger", "sample_files/sample.go"})
	assertEqual(t, exitCode, 1, "")
	assertStringEqual(t, "Could not parse noninteger as int\n", errOut.String(), "")
}

func TestColorization(t *testing.T) {
	out := bytes.Buffer{}
	errOut := bytes.Buffer{}
	in := bytes.Buffer{}

	run(&out, &errOut, &in, []string{"-n", "-c", "19", "sample_files/sample.go"})
	expected := fmt.Sprintf(`
sample_files/sample.go

%s 1:%spackage main
%s19:%s%sfunc (fb FizzBuzzMap) String(i int) string {%s

`, BLUE, END_COLOR, BLUE, END_COLOR, GREEN, END_COLOR)
	assertStringEqual(t, expected, out.String(), "")

	// Test with string match coloring:
	out.Reset()
	errOut.Reset()

	run(&out, &errOut, &in, []string{"-n", "-c", "-e", "(fb FizzBuzzMap)", "sample_files/sample.go"})
	expected = fmt.Sprintf(`
sample_files/sample.go

%s 1:%spackage main
%s19:%sfunc %s(fb FizzBuzzMap)%s String(i int) string {

`, BLUE, END_COLOR, BLUE, END_COLOR, GREEN, END_COLOR)
	assertStringEqual(t, expected, out.String(), "")

}

func TestReadFromStdin(t *testing.T) {
	out := bytes.Buffer{}
	errOut := bytes.Buffer{}
	in := strings.NewReader(`package main
import "fmt"
func main() {
	if false {
		fmt.Println("false")
	}
	fmt.Println("true")
}
`)
	exitStatus := run(&out, &errOut, in, []string{"-n", "-l", "go", "7"})
	assertEqual(t, exitStatus, 0, "")
	expected := `
<stdin>

1:package main
3:func main() {
7:	fmt.Println("true")

`
	assertStringEqual(t, expected, out.String(), "")
}

func TestProcessFileUnsupportedLanguage(t *testing.T) {
	search := Search{}
	opts := Options{Language: "foobar"}
	err := processFile("sample_files/sample.go", search, opts)
	assertNotEqual(t, err, nil, "")
	assertInString(t, "Unknown language", err.Error(), "")
}
