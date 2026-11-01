// Package main is a negative test fixture: coverage is deliberately below the gate.
package main

// Add returns the sum of a and b.
func Add(a, b int) int {
	return a + b
}

// Sub returns the difference of a and b. It has no test coverage.
func Sub(a, b int) int {
	return a - b
}

func main() {}
