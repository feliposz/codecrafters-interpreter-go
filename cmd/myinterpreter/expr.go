package main

import (
	"fmt"
	"strings"
)

type Expr interface {
	String() string
	Evaluate() any
	Resolve()
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

type Assign struct {
	Name  *Variable
	Value Expr
}

func (a *Assign) String() string {
	return fmt.Sprintf("(= %s %s)", a.Name.String(), a.Value.String())
}

type Logical struct {
	left     Expr
	operator *Token
	right    Expr
}

func (l *Logical) String() string {
	return fmt.Sprintf("(%s %s %s)", l.operator.Str, l.left.String(), l.right.String())
}

type Call struct {
	callee    Expr
	paren     *Token
	arguments []Expr
}

func (c *Call) String() string {
	sb := strings.Builder{}
	for _, arg := range c.arguments {
		sb.WriteString(" ")
		sb.WriteString(arg.String())
	}
	return fmt.Sprintf("(call %s%s)", c.callee, sb.String())
}
