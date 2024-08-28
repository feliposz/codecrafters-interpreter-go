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

func (p *Parser) consume(t TokenType, msg string) *Token {
	if !p.check(t) {
		loxError(p.peek(), msg)
	}
	return p.advance()
}

func (p *Parser) parse() []Stmt {
	statements := []Stmt{}
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements
}

func (p *Parser) declaration() Stmt {
	if p.match(VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) varDeclaration() Stmt {
	name := p.consume(IDENTIFIER, "Expect variable name.")

	var initializer Expr
	if p.match(EQUAL) {
		initializer = p.expression()
	}

	p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	return &VarStatement{name, initializer}
}

func (p *Parser) statement() Stmt {
	if p.match(PRINT) {
		return p.printStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) printStatement() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Expect ';' after value.")
	return &PrintStatement{expr}
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Expect ';' after expression.")
	return &ExpressionStatement{expr}
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	expr := p.equality()

	if p.match(EQUAL) {
		equals := p.previous()
		value := p.assignment()
		if name, ok := expr.(*Variable); ok {
			return &Assign{name, value}
		}
		loxError(equals, "Invalid assignment target.")
	}
	return expr
}

func (p *Parser) primary() Expr {
	if p.match(NIL, TRUE, FALSE, NUMBER, STRING) {
		return &Literal{p.previous()}
	}
	if p.match(LEFT_PAREN) {
		paren := p.previous()
		expr := p.expression()
		if expr == nil {
			loxError(p.peek(), "Expected expression.")
		}
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return &Grouping{paren, expr}
	}
	if p.match(IDENTIFIER) {
		return &Variable{p.previous()}
	}
	loxError(p.peek(), "Expected expression.")
	return nil
}

func (p *Parser) unary() Expr {
	if p.match(MINUS, BANG) {
		op := p.previous()
		expr := p.unary()
		return &Unary{op, expr}
	}
	return p.primary()
}

func (p *Parser) factor() Expr {
	left := p.unary()
	for p.match(STAR, SLASH) {
		op := p.previous()
		right := p.unary()
		left = &Binary{op, left, right}
	}
	return left
}

func (p *Parser) term() Expr {
	left := p.factor()
	for p.match(PLUS, MINUS) {
		op := p.previous()
		right := p.factor()
		left = &Binary{op, left, right}
	}
	return left
}

func (p *Parser) comparison() Expr {
	left := p.term()
	for p.match(LESS, GREATER, LESS_EQUAL, GREATER_EQUAL) {
		op := p.previous()
		right := p.term()
		left = &Binary{op, left, right}
	}
	return left
}

func (p *Parser) equality() Expr {
	left := p.comparison()
	for p.match(EQUAL_EQUAL, BANG_EQUAL) {
		op := p.previous()
		right := p.comparison()
		left = &Binary{op, left, right}
	}
	return left
}
