package main

import "fmt"

type Expr interface {
	String() string
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
	expr Expr
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
