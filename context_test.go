package main

import (
	"bytes"
	"testing"
)

func assertOutput(t *testing.T, expected string, result string, msg string) {
	if result != expected {
		t.Fatalf("%s\nActual:\n<<<%v>>>\n\nExpected:\n<<<%v>>>'", msg, result, expected)
	}
}

func TestProcessFileGoLineNumber(t *testing.T) {
	out := bytes.Buffer{}
	opts := Options{
		PrintNums: true,
		Out:       &out,
	}

	// Test case 1, nothing special:
	search := NewSearch()
	search.ValInts = []int{24}
	expected := `
sample_files/sample.go

 1:package main
19:func (fb FizzBuzzMap) String(i int) string {
22:	for _, x := range fb.Map {
23:		if i%x.Divisor == 0 {
25:			result += x.String

`
	processFile("sample_files/sample.go", search, opts)
	assertOutput(t, expected, out.String(), "")

	// Test case 2, multiple if/elseif clauses:
	search.ValInts = []int{53}
	expected = `
sample_files/sample.go

 1:package main
42:func classicFizzBuzz(start int, end int) {
43:	for i := start; i <= end; i++ {
44:		if i%15 == 0 {
47:		} else if i%3 == 0 {
50:		} else if i%5 == 0 {
53:		} else {
54:			fmt.Println(i)

`
	out.Reset()
	processFile("sample_files/sample.go", search, opts)
	assertOutput(t, expected, out.String(), "Failed if/elseif line number case")

	// Test case 3, multiple switch/case clauses:
	search.ValInts = []int{71}
	expected = `
sample_files/sample.go

 1:package main
59:func hardcodedFizzBuzz() {
60:	for i := 1; i <= 15; i++ {
61:		switch i {
62:		case 15:
65:		case 3, 6, 9, 12:
68:		case 5, 10:
71:		default:
72:			fmt.Println(i)

`
	out.Reset()
	processFile("sample_files/sample.go", search, opts)
	assertOutput(t, expected, out.String(), "")
}

func TestProcessFileGoStringMatch(t *testing.T) {
	out := bytes.Buffer{}
	opts := Options{
		PrintNums: true,
		Out:       &out,
	}

	search := NewSearch()
	search.SetRegexMatch()
	search.Val = "\"Buzz\""
	expected := `
sample_files/sample.go

 1:package main
42:func classicFizzBuzz(start int, end int) {
43:	for i := start; i <= end; i++ {
44:		if i%15 == 0 {
47:		} else if i%3 == 0 {
50:		} else if i%5 == 0 {
51:			s := "Buzz"
59:func hardcodedFizzBuzz() {
60:	for i := 1; i <= 15; i++ {
61:		switch i {
62:		case 15:
65:		case 3, 6, 9, 12:
68:		case 5, 10:
69:			s := "Buzz"
77:func main() {
85:	fb := FizzBuzzMap{
86:		Map: []DivisorString{
88:			DivisorString{5, "Buzz"},

`
	processFile("sample_files/sample.go", search, opts)
	assertOutput(t, expected, out.String(), "")
}

func TestNoMatchingLines(t *testing.T) {
	out := bytes.Buffer{}
	opts := Options{
		PrintNums: true,
		Out:       &out,
	}
	search := NewSearch()
	search.ValInts = []int{999999}

	processFile("sample_files/sample.go", search, opts)
	assertOutput(t, "", out.String(), "")
}

func TestUnsupportedLanguage(t *testing.T) {
	out := bytes.Buffer{}
	opts := Options{
		PrintNums: true,
		Out:       &out,
		Language:  "fake",
	}
	search := NewSearch()
	search.ValInts = []int{53}

	err := processFile("sample_files/sample.go", search, opts)
	if err.Error() != "Unknown language for \'fake\'" {
		t.Fatalf("Expected error for nonexistent file, got %q", err.Error())
	}
}

func TestNonexistentFile(t *testing.T) {
	out := bytes.Buffer{}
	opts := Options{
		PrintNums: true,
		Out:       &out,
	}
	search := NewSearch()

	err := processFile("not/a/real/file.py", search, opts)
	if err == nil {
		t.Fatalf("Expected an error for nonexistent file, got nil")
	}
}
