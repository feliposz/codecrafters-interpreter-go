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
	EQUAL
	EQUAL_EQUAL
	BANG
	BANG_EQUAL
	LESS
	LESS_EQUAL
	GREATER
	GREATER_EQUAL
	SLASH
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
	case EQUAL:
		return "EQUAL"
	case EQUAL_EQUAL:
		return "EQUAL_EQUAL"
	case BANG:
		return "BANG"
	case BANG_EQUAL:
		return "BANG_EQUAL"
	case LESS:
		return "LESS"
	case LESS_EQUAL:
		return "LESS_EQUAL"
	case GREATER:
		return "GREATER"
	case GREATER_EQUAL:
		return "GREATER_EQUAL"
	case SLASH:
		return "SLASH"
	}
	return "UNKNOWN"
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
	var tokenStr []byte
	lexicalErrors := false
	for i := 0; i < len(fileContents); i++ {
		ch := fileContents[i]
		tokenStr = fileContents[i : i+1]
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
		case '=':
			tt = EQUAL
			if i+1 < len(fileContents) && fileContents[i+1] == '=' {
				tt = EQUAL_EQUAL
				tokenStr = fileContents[i : i+2]
				i++
			}
		case '!':
			tt = BANG
			if i+1 < len(fileContents) && fileContents[i+1] == '=' {
				tt = BANG_EQUAL
				tokenStr = fileContents[i : i+2]
				i++
			}
		case '<':
			tt = LESS
			if i+1 < len(fileContents) && fileContents[i+1] == '=' {
				tt = LESS_EQUAL
				tokenStr = fileContents[i : i+2]
				i++
			}
		case '>':
			tt = GREATER
			if i+1 < len(fileContents) && fileContents[i+1] == '=' {
				tt = GREATER_EQUAL
				tokenStr = fileContents[i : i+2]
				i++
			}
		case '/':
			tt = SLASH
			if i+1 < len(fileContents) && fileContents[i+1] == '/' {
				tt = UNKNOWN
				for i < len(fileContents) && fileContents[i] != '\n' {
					i++
				}
			}
		case '\n':
			line++
		default:
			fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", line, ch)
			tt = UNKNOWN
			lexicalErrors = true
		}
		if tt != UNKNOWN {
			fmt.Printf("%v %s null\n", tt, tokenStr)
		}
	}
	fmt.Println("EOF  null")
	if lexicalErrors {
		os.Exit(65)
	}
}
