package main

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens, 0}
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) peek() *Token {
	return &p.tokens[p.current]
}

func (p *Parser) previous() *Token {
	return &p.tokens[p.current-1]
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) parse() Expr {
	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(NIL, TRUE, FALSE, NUMBER) {
		return &Literal{p.previous()}
	}
	return nil
}
