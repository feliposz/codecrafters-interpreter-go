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
		switch value := value.(type) {
		case float64:
			return -value
		}
		runtimeError("Operand must be a number.")
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
		runtimeError("Operands must be two numbers or two strings.")
	case MINUS:
		switch left := left.(type) {
		case float64:
			switch right := right.(type) {
			case float64:
				return left - right
			}
		}
		runtimeError("Operands must be numbers.")
	case STAR:
		switch left := left.(type) {
		case float64:
			switch right := right.(type) {
			case float64:
				return left * right
			}
		}
		runtimeError("Operands must be numbers.")
	case SLASH:
		switch left := left.(type) {
		case float64:
			switch right := right.(type) {
			case float64:
				return left / right
			}
		}
		runtimeError("Operands must be numbers.")
	case LESS:
		switch left := left.(type) {
		case float64:
			switch right := right.(type) {
			case float64:
				return left < right
			}
		}
		runtimeError("Operands must be numbers.")
	case GREATER:
		switch left := left.(type) {
		case float64:
			switch right := right.(type) {
			case float64:
				return left > right
			}
		}
		runtimeError("Operands must be numbers.")
	case LESS_EQUAL:
		switch left := left.(type) {
		case float64:
			switch right := right.(type) {
			case float64:
				return left <= right
			}
		}
		runtimeError("Operands must be numbers.")
	case GREATER_EQUAL:
		switch left := left.(type) {
		case float64:
			switch right := right.(type) {
			case float64:
				return left >= right
			}
		}
		runtimeError("Operands must be numbers.")
	case EQUAL_EQUAL:
		return left == right
	case BANG_EQUAL:
		return left != right
	}
	loxError(b.Op, "not implemented")
	return nil
}
