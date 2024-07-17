package main

import "testing"

func TestNumbers(t *testing.T) {
	tokenizer([]byte("1234.1234\n.123\n456.\n123"), true)
}

func TestIdentifiers(t *testing.T) {
	tokenizer([]byte("_123bar f00 6az bar 6ar"), true)
}
