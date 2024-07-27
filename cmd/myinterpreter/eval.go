package main

func (l *Literal) Evaluate() any {
	switch l.token.Type {
	case NIL:
		return nil
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
	return g.expr.Evaluate()
}

func (u *Unary) Evaluate() any {
	value := u.Expr.Evaluate()
	switch u.Op.Type {
	case MINUS:
		return -value.(float64)
	case BANG:
		if value == nil || value == false {
			return true
		}
		return false
	}
	loxError(u.Op, "invalid op")
	return nil
}

func (b *Binary) Evaluate() any {
	left, right := b.Left.Evaluate(), b.Right.Evaluate()
	switch b.Op.Type {
	case PLUS:
		switch left := left.(type) {
		case string:
			switch right := right.(type) {
			case string:
				return left + right
			}
		case float64:
			switch right := right.(type) {
			case float64:
				return left + right
			}
		}
	case MINUS:
		return left.(float64) - right.(float64)
	case STAR:
		return left.(float64) * right.(float64)
	case SLASH:
		return left.(float64) / right.(float64)
	case LESS:
		return left.(float64) < right.(float64)
	case GREATER:
		return left.(float64) > right.(float64)
	case LESS_EQUAL:
		return left.(float64) <= right.(float64)
	case GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	}
	loxError(b.Op, "not implemented")
	return nil
}
