package main

import (
	"fmt"
	"testing"
)

func TestGrouping(t *testing.T) {
	tokens := tokenizer([]byte("()"), false)
	parser := NewParser(tokens)
	expr := parser.parse()
	fmt.Println(expr)
}
