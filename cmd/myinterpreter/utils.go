package main

import (
	"fmt"
	"os"
)

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func loxError(token *Token, msg string) {
	fmt.Fprint(os.Stderr, msg)
	os.Exit(65)
}
