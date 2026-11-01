// Package main is a negative test fixture: it calls a symbol reachable to a
// known vulnerability (golang.org/x/text v0.3.7, GO-2021-0113).
package main

import (
	"fmt"

	"golang.org/x/text/language"
)

func main() {
	_, _, err := language.ParseAcceptLanguage("en-US")
	if err != nil {
		fmt.Println(err)
	}
}
