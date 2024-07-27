package main

func (l *Literal) Evaluate() any {
	switch l.token.Type {
	case NIL:
		return "nil"
	case TRUE:
		return true
	case FALSE:
		return false
	case NUMBER:
		return l.token.Content.(float64)
	case STRING:
		return l.token.Content.(string)
	}
	loxError(l.token, "unknow type")
	return nil
}

func (g *Grouping) Evaluate() any {
	loxError(g.token, "not implemented")
	return nil
}

func (u *Unary) Evaluate() any {
	loxError(u.Op, "not implemented")
	return nil
}

func (b *Binary) Evaluate() any {
	loxError(b.Op, "not implemented")
	return nil
}
