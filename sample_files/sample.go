package main

import (
	"fmt"
)

// This is an intentionally absurd implementation of FizzBuzz.  Its purpose is to
// have a real runnable program that exercises multiple language features and syntax.

type DivisorString struct {
	Divisor int
	String  string
}

type FizzBuzzMap struct {
	Map []DivisorString
}

func (fb FizzBuzzMap) String(i int) string {
	matched := false
	var result string
	for _, x := range fb.Map {
		if i%x.Divisor == 0 {
			matched = true
			result += x.String
		}
	}
	if matched {
		return result
	} else {
		return fmt.Sprintf("%d", i)
	}
}

func fizzBuzz(start int, end int, fb FizzBuzzMap) {
	for i := start; i <= end; i++ {
		s := fb.String(i)
		fmt.Println(s)
	}
}

func classicFizzBuzz(start int, end int) {
	for i := start; i <= end; i++ {
		if i%15 == 0 {
			s := "FizzBuzz"
			fmt.Println(s)
		} else if i%3 == 0 {
			s := "Fizz"
			fmt.Println(s)
		} else if i%5 == 0 {
			s := "Buzz"
			fmt.Println(s)
		} else {
			fmt.Println(i)
		}
	}
}

func hardcodedFizzBuzz() {
	for i := 1; i <= 15; i++ {
		switch i {
		case 15:
			s := "FizzBuzz"
			fmt.Println(s)
		case 3, 6, 9, 12:
			s := "Fizz"
			fmt.Println(s)
		case 5, 10:
			s := "Buzz"
			fmt.Println(s)
		default:
			fmt.Println(i)
		}
	}
}

func main() {
	fmt.Println("Classic:")
	classicFizzBuzz(1, 20)
	fmt.Println()
	fmt.Println("Brute force:")
	hardcodedFizzBuzz()
	fmt.Println()
	fmt.Println("Unnecessarily complicated:")
	fb := FizzBuzzMap{
		Map: []DivisorString{
			DivisorString{3, "Fizz"},
			DivisorString{5, "Buzz"},
			DivisorString{7, "Bang"},
		},
	}
	fizzBuzz(1, 35, fb)
}
