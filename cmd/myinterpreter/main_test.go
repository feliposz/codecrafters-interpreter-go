package main

import "testing"

func TestNumbers(t *testing.T) {
	tokenizer([]byte("1234.1234\n.123\n456.\n123"))
}
