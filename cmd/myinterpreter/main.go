package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh <command> <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	switch command {
	case "tokenize":
		tokenizer(fileContents, true)
	case "parse":
		tokens := tokenizer(fileContents, false)
		parser := NewParser(tokens)
		expr := parser.expression()
		fmt.Println(expr)
	case "evaluate":
		tokens := tokenizer(fileContents, false)
		parser := NewParser(tokens)
		expr := parser.expression()
		result := expr.Evaluate()
		if result == nil {
			fmt.Println("nil")
		} else {
			fmt.Println(result)
		}
	case "run":
		tokens := tokenizer(fileContents, false)
		parser := NewParser(tokens)
		statements := parser.parse()
		for _, statement := range statements {
			statement.Run()
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}
