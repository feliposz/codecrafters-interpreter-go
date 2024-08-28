package main

import "fmt"

type Expr interface {
	String() string
	Evaluate() any
}

type Literal struct {
	token *Token
}

func (l *Literal) String() string {
	switch l.token.Type {
	case NIL:
		return "nil"
	case TRUE:
		return "true"
	case FALSE:
		return "false"
	case NUMBER:
		return FloatFormat(l.token.Content.(float64))
	case STRING:
		return l.token.Content.(string)
	}
	return "?"
}

type Grouping struct {
	token *Token
	expr  Expr
}

func (g *Grouping) String() string {
	return fmt.Sprintf("(group %s)", g.expr.String())
}

type Unary struct {
	Op   *Token
	Expr Expr
}

func (u *Unary) String() string {
	if u.Op.Type == MINUS {
		return fmt.Sprintf("(- %s)", u.Expr.String())
	} else {
		return fmt.Sprintf("(! %s)", u.Expr.String())
	}
}

type Binary struct {
	Op    *Token
	Left  Expr
	Right Expr
}

func (b *Binary) String() string {
	op := "?"
	switch b.Op.Type {
	case PLUS:
		op = "+"
	case MINUS:
		op = "-"
	case STAR:
		op = "*"
	case SLASH:
		op = "/"
	case LESS:
		op = "<"
	case GREATER:
		op = ">"
	case LESS_EQUAL:
		op = "<="
	case GREATER_EQUAL:
		op = ">="
	case EQUAL_EQUAL:
		op = "=="
	case BANG_EQUAL:
		op = "!="
	}
	return fmt.Sprintf("(%s %s %s)", op, b.Left.String(), b.Right.String())
}

type Variable struct {
	Name *Token
}

func (v *Variable) String() string {
	return fmt.Sprintf("(var %s)", v.Name.String())
}
