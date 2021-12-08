package main

import (
	"bytes"
	"testing"
)

func TestGoLineNumberBasic(t *testing.T) {
	out := bytes.Buffer{}
	errOut := bytes.Buffer{}
	in := bytes.Buffer{}

	run(&out, &errOut, &in, []string{"-n", "25", "sample_files/sample.go"})
	expected := `
sample_files/sample.go

 1:package main
19:func (fb FizzBuzzMap) String(i int) string {
22:	for _, x := range fb.Map {
23:		if i%x.Divisor == 0 {
25:			result += x.String

`
	assertStringEqual(t, expected, out.String(), "")
}

func TestGoLineNumberMultipleIf(t *testing.T) {
	out := bytes.Buffer{}
	errOut := bytes.Buffer{}
	in := bytes.Buffer{}

	run(&out, &errOut, &in, []string{"-n", "54", "sample_files/sample.go"})
	expected := `
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
	assertStringEqual(t, expected, out.String(), "Failed if/elseif line number case")
}

func TestGoLineNumberSwitchCase(t *testing.T) {
	out := bytes.Buffer{}
	errOut := bytes.Buffer{}
	in := bytes.Buffer{}

	run(&out, &errOut, &in, []string{"-n", "72", "sample_files/sample.go"})
	expected := `
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
	assertStringEqual(t, expected, out.String(), "")
}

func TestGoStringMatch(t *testing.T) {
	out := bytes.Buffer{}
	errOut := bytes.Buffer{}
	in := bytes.Buffer{}

	run(&out, &errOut, &in, []string{"-n", "-e", "\"Buzz\"", "sample_files/sample.go"})

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
	assertStringEqual(t, expected, out.String(), "")
}
