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
	if p.match(CLASS) {
		return p.classDeclaration()
	}
	if p.match(FUN) {
		return p.function("function")
	}
	if p.match(VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) classDeclaration() Stmt {
	name := p.consume(IDENTIFIER, "Expect class name.")
	var superclass *Variable
	if p.match(LESS) {
		identifier := p.consume(IDENTIFIER, "Expect superclass name.")
		superclass = &Variable{identifier}
	}
	p.consume(LEFT_BRACE, "Expect '{' before class body.")
	var methods []*FunctionDeclaration
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		methods = append(methods, p.function("method"))
	}
	p.consume(RIGHT_BRACE, "Expect '}' after class body.")
	return &ClassDeclaration{name, superclass, methods}
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

func (p *Parser) whileStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after condition.")
	body := p.statement()
	return &WhileStatement{condition, body}
}

func (p *Parser) forStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'for'.")
	var initializer Stmt
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}
	var condition Expr
	if !p.check(SEMICOLON) {
		condition = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after loop condition.")
	var increment Expr
	if !p.check(RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(RIGHT_PAREN, "Expect ')' after for clauses.")
	body := p.statement()
	if increment != nil {
		body = &Block{Statements: []Stmt{body, &ExpressionStatement{increment}}}
	}
	if condition == nil {
		condition = &Literal{&Token{Type: TRUE}}
	}
	body = &WhileStatement{condition, body}
	if initializer != nil {
		body = &Block{Statements: []Stmt{initializer, body}}
	}
	return body
}

func (p *Parser) statement() Stmt {
	if p.match(FOR) {
		return p.forStatement()
	}
	if p.match(IF) {
		return p.ifStatement()
	}
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(RETURN) {
		return p.returnStatement()
	}
	if p.match(WHILE) {
		return p.whileStatement()
	}
	if p.match(LEFT_BRACE) {
		return &Block{p.block()}
	}
	return p.expressionStatement()
}

func (p *Parser) ifStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after if condition.")
	thenBranch := p.statement()
	var elseBranch Stmt
	if p.match(ELSE) {
		elseBranch = p.statement()
	}
	return &IfStatement{condition, thenBranch, elseBranch}
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

func (p *Parser) function(kind string) *FunctionDeclaration {
	name := p.consume(IDENTIFIER, "Expect "+kind+" name.")
	p.consume(LEFT_PAREN, "Expect '(' after "+kind+" name.")
	parameters := []*Token{}
	if !p.check(RIGHT_PAREN) {
		for {
			if len(parameters) > 255 {
				loxError(p.peek(), "Can't have more than 255 parameters.")
			}
			parameters = append(parameters, p.consume(IDENTIFIER, "Expect parameter name."))
			if !p.match(COMMA) {
				break
			}
		}
	}
	p.consume(RIGHT_PAREN, "Expect ')' after parameters.")
	p.consume(LEFT_BRACE, "Expect '{' before "+kind+" body.")
	body := p.block()
	return &FunctionDeclaration{name, parameters, body}
}

func (p *Parser) block() []Stmt {
	statements := []Stmt{}
	for !p.isAtEnd() && !p.check(RIGHT_BRACE) {
		statements = append(statements, p.declaration())
	}
	p.consume(RIGHT_BRACE, "Expect '}' after block.")
	return statements
}

func (p *Parser) returnStatement() Stmt {
	keyword := p.previous()
	var value Expr
	if !p.check(SEMICOLON) {
		value = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after return value.")
	return &ReturnStatement{keyword, value}
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	expr := p.or()
	if p.match(EQUAL) {
		equals := p.previous()
		value := p.assignment()
		if name, ok := expr.(*Variable); ok {
			return &Assign{name, value}
		} else if get, ok := expr.(*Get); ok {
			return &Set{get.object, get.name, value}
		}
		loxError(equals, "Invalid assignment target.")
	}
	return expr
}

func (p *Parser) or() Expr {
	expr := p.and()
	for p.match(OR) {
		operator := p.previous()
		right := p.and()
		expr = &Logical{expr, operator, right}
	}
	return expr
}

func (p *Parser) and() Expr {
	expr := p.equality()
	for p.match(AND) {
		operator := p.previous()
		right := p.and()
		expr = &Logical{expr, operator, right}
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
	if p.match(THIS) {
		return &This{p.previous()}
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
	return p.call()
}

func (p *Parser) call() Expr {
	expr := p.primary()
	for {
		if p.match(LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else if p.match(DOT) {
			name := p.consume(IDENTIFIER, "Expect property name after '.'.")
			expr = &Get{expr, name}
		} else {
			break
		}
	}
	return expr
}

func (p *Parser) finishCall(callee Expr) Expr {
	arguments := []Expr{}
	if !p.check(RIGHT_PAREN) {
		for {
			if len(arguments) > 255 {
				loxError(p.peek(), "Can't have more than 255 arguments.")
			}
			arguments = append(arguments, p.expression())
			if !p.match(COMMA) {
				break
			}
		}
	}
	paren := p.consume(RIGHT_PAREN, "Expect ')' after arguments.")
	return &Call{callee, paren, arguments}
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
