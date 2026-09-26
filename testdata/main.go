// Package main is a test fixture for the Ultimate Go CI Suite action.
package main

import (
	"fmt"
	"log"
)

func main() {
	result := Foo()
	log.Printf("Result: %v\n", result)
}

// Foo returns a greeting message.
func Foo() string {
	return "Hello, World!"
}

// Bar performs a simple calculation.
func Bar(x, y int) int {
	return x + y
}

// Baz demonstrates error handling.
func Baz(n int) (int, error) {
	if n < 0 {
		return 0, fmt.Errorf("negative number not allowed: %d", n)
	}
	return n * 2, nil
}
