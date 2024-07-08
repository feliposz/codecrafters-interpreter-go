package main

import (
	"fmt"
	"os"
)

type TokenType uint8

const (
	UNKNOWN TokenType = iota
	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	STAR
)

func (tt TokenType) String() string {
	switch tt {
	case LEFT_PAREN:
		return "LEFT_PAREN"
	case RIGHT_PAREN:
		return "RIGHT_PAREN"
	case LEFT_BRACE:
		return "LEFT_BRACE"
	case RIGHT_BRACE:
		return "RIGHT_BRACE"
	case COMMA:
		return "COMMA"
	case DOT:
		return "DOT"
	case MINUS:
		return "MINUS"
	case PLUS:
		return "PLUS"
	case SEMICOLON:
		return "SEMICOLON"
	case STAR:
		return "STAR"
	}
	return "unknown"
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	line := 1
	var tt TokenType
	lexicalErrors := false
	for _, ch := range fileContents {
		switch ch {
		case '(':
			tt = LEFT_PAREN
		case ')':
			tt = RIGHT_PAREN
		case '{':
			tt = LEFT_BRACE
		case '}':
			tt = RIGHT_BRACE
		case ',':
			tt = COMMA
		case '.':
			tt = DOT
		case '-':
			tt = MINUS
		case '+':
			tt = PLUS
		case ';':
			tt = SEMICOLON
		case '*':
			tt = STAR
		case '\n':
			line++
		default:
			fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", line, ch)
			tt = UNKNOWN
			lexicalErrors = true
		}
		if tt != UNKNOWN {
			fmt.Printf("%v %c null\n", tt, ch)
		}
	}
	fmt.Println("EOF  null")
	if lexicalErrors {
		os.Exit(65)
	}
}
