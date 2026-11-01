// Package main is a negative test fixture: it contains deliberate lint errors.
package main

import "fmt"

// unusedHelper is never called and must be flagged by the unused linter.
func unusedHelper() int {
	return 42
}

func main() {
	x := 1 // ineffectual assignment: overwritten before use
	x = 2
	fmt.Println(x)
}
