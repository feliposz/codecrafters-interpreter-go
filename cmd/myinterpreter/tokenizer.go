package main

import (
	"fmt"
	"os"
	"strconv"
)

type TokenType uint8

const (
	UNKNOWN TokenType = iota
	EOF
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
	STRING
	NUMBER
	IDENTIFIER
	COMMENT
	AND
	CLASS
	ELSE
	FALSE
	FOR
	FUN
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE
)

func (tt TokenType) String() string {
	switch tt {
	case EOF:
		return "EOF"
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
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case IDENTIFIER:
		return "IDENTIFIER"
	case COMMENT:
		return "COMMENT"
	case AND:
		return "AND"
	case CLASS:
		return "CLASS"
	case ELSE:
		return "ELSE"
	case FALSE:
		return "FALSE"
	case FOR:
		return "FOR"
	case FUN:
		return "FUN"
	case IF:
		return "IF"
	case NIL:
		return "NIL"
	case OR:
		return "OR"
	case PRINT:
		return "PRINT"
	case RETURN:
		return "RETURN"
	case SUPER:
		return "SUPER"
	case THIS:
		return "THIS"
	case TRUE:
		return "TRUE"
	case VAR:
		return "VAR"
	case WHILE:
		return "WHILE"
	}
	return "UNKNOWN"
}

var reservedKeywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

type Token struct {
	Type    TokenType
	Str     string
	Content any
	Line    int
}

func (t Token) String() string {
	switch t.Type {
	case STRING:
		return fmt.Sprintf("%v %s %s", t.Type, t.Str, t.Content)
	case NUMBER:
		value, _ := t.Content.(float64)
		return fmt.Sprintf("%v %s %s", t.Type, t.Str, FloatFormat(value))
	case COMMENT, UNKNOWN:
		return ""
	default:
		return fmt.Sprintf("%v %s null", t.Type, t.Str)
	}
}

func FloatFormat(value float64) string {
	if value == float64(int(value)) {
		return fmt.Sprintf("%.1f", value)
	} else {
		return fmt.Sprintf("%g", value)
	}
}

func tokenizer(fileContents []byte, print bool) []Token {
	line := 1
	result := []Token{}
	var tt TokenType
	var tokenStr []byte
	lexicalErrors := false
	for i := 0; i < len(fileContents); i++ {
		ch := fileContents[i]
		var content any = "null"
		if ch == ' ' || ch == '\t' {
			continue
		}
		if ch == '\n' {
			line++
			continue
		}
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
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			tt = NUMBER
			j := i
			hasDot := false
			for j < len(fileContents) && (fileContents[j] == '.' || isDigit(fileContents[j])) {
				if fileContents[j] == '.' {
					if j+1 >= len(fileContents) || !isDigit(fileContents[j+1]) {
						break
					}
					if hasDot {
						break
					} else {
						hasDot = true
					}
				}
				j++
			}
			if j <= len(fileContents) {
				j--
			}
			tokenStr = fileContents[i : j+1]
			i = j
			content, _ = strconv.ParseFloat(string(tokenStr), 64)
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
				tt = COMMENT
				for i < len(fileContents) && fileContents[i] != '\n' {
					i++
				}
				if i < len(fileContents) && fileContents[i] == '\n' {
					line++
				}
			}
		case '"':
			j := i + 1
			for j < len(fileContents) && fileContents[j] != '"' {
				j++
			}
			if j < len(fileContents) && fileContents[j] == '"' {
				tt = STRING
				tokenStr = fileContents[i : j+1]
				content = string(fileContents[i+1 : j])
			} else {
				fmt.Fprintf(os.Stderr, "[line %d] Error: Unterminated string.\n", line)
				tt = UNKNOWN
				lexicalErrors = true
			}
			i = j
		default:
			if ch == '_' || isLetter(ch) {
				tt = IDENTIFIER
				j := i
				for j < len(fileContents) && (fileContents[j] == '_' || isLetter(fileContents[j]) || isDigit(fileContents[j])) {
					j++
				}
				if j <= len(fileContents) {
					j--
				}
				tokenStr = fileContents[i : j+1]
				i = j
				if kwType, found := reservedKeywords[string(tokenStr)]; found {
					tt = kwType
				}
			} else {
				fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", line, ch)
				tt = UNKNOWN
				lexicalErrors = true
			}
		}
		if tt == UNKNOWN || tt == COMMENT {
			continue
		}
		token := Token{
			tt,
			string(tokenStr),
			content,
			line,
		}
		result = append(result, token)
		if print {
			fmt.Println(token)
		}
	}
	eofToken := Token{
		EOF,
		"",
		"",
		line,
	}
	result = append(result, eofToken)
	if print {
		fmt.Println(eofToken)
	}
	if lexicalErrors {
		os.Exit(65)
	}
	return result
}
