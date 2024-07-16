package main

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
