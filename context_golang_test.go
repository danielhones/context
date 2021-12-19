package main

import (
	"bytes"
	"testing"
)

func TestGoLineNumberBasic(t *testing.T) {
	out := bytes.Buffer{}
	errOut := bytes.Buffer{}
	in := bytes.Buffer{}

	run(&out, &errOut, &in, []string{"-n", "29", "sample_files/sample.go"})
	expected := `
sample_files/sample.go

 1:package main
23:func (fb FizzBuzzMap) String(i int) string {
26:	for _, x := range fb.Map {
27:		if i%x.Divisor == 0 {
29:			result += x.String

`
	assertStringEqual(t, expected, out.String(), "")
}

func TestGoLineNumberMultipleIf(t *testing.T) {
	out := bytes.Buffer{}
	errOut := bytes.Buffer{}
	in := bytes.Buffer{}

	run(&out, &errOut, &in, []string{"-n", "73", "sample_files/sample.go"})
	expected := `
sample_files/sample.go

 1:package main
46:func classicFizzBuzz(start int, end int) {
47:	for i := start; i <= end; i++ {
48:		if i%15 == 0 {
57:		} else if i%3 == 0 {
63:		} else if i%5 == 0 {
69:		} else {
73:			fmt.Println(i)

`
	assertStringEqual(t, expected, out.String(), "Failed if/elseif line number case")
}

func TestGoLineNumberSwitchCase(t *testing.T) {
	out := bytes.Buffer{}
	errOut := bytes.Buffer{}
	in := bytes.Buffer{}

	run(&out, &errOut, &in, []string{"-n", "91", "sample_files/sample.go"})
	expected := `
sample_files/sample.go

 1:package main
78:func hardcodedFizzBuzz() {
79:	for i := 1; i <= 15; i++ {
80:		switch i {
81:		case 15:
84:		case 3, 6, 9, 12:
87:		case 5, 10:
90:		default:
91:			fmt.Println(i)

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
 46:func classicFizzBuzz(start int, end int) {
 47:	for i := start; i <= end; i++ {
 48:		if i%15 == 0 {
 57:		} else if i%3 == 0 {
 63:		} else if i%5 == 0 {
 67:			s := "Buzz"
 78:func hardcodedFizzBuzz() {
 79:	for i := 1; i <= 15; i++ {
 80:		switch i {
 81:		case 15:
 84:		case 3, 6, 9, 12:
 87:		case 5, 10:
 88:			s := "Buzz"
 96:func main() {
104:	fb := FizzBuzzMap{
105:		Map: []DivisorString{
107:			DivisorString{5, "Buzz"},

`
	assertStringEqual(t, expected, out.String(), "")
}
