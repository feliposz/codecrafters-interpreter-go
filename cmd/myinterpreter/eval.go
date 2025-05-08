package main

import "fmt"

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

func (l *Logical) Evaluate() any {
	left := l.left.Evaluate()
	switch l.operator.Type {
	case OR:
		if isTruthy(left) {
			return left
		}
	case AND:
		if !isTruthy(left) {
			return left
		}
	default:
		loxError(l.operator, "unknow operator")
		return nil
	}
	return l.right.Evaluate()
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
		runtimeError(u.Op, "Operand must be a number.")
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
		runtimeError(b.Op, "Operands must be two numbers or two strings.")
	case MINUS:
		switch left := left.(type) {
		case float64:
			switch right := right.(type) {
			case float64:
				return left - right
			}
		}
		runtimeError(b.Op, "Operands must be numbers.")
	case STAR:
		switch left := left.(type) {
		case float64:
			switch right := right.(type) {
			case float64:
				return left * right
			}
		}
		runtimeError(b.Op, "Operands must be numbers.")
	case SLASH:
		switch left := left.(type) {
		case float64:
			switch right := right.(type) {
			case float64:
				return left / right
			}
		}
		runtimeError(b.Op, "Operands must be numbers.")
	case LESS:
		switch left := left.(type) {
		case float64:
			switch right := right.(type) {
			case float64:
				return left < right
			}
		}
		runtimeError(b.Op, "Operands must be numbers.")
	case GREATER:
		switch left := left.(type) {
		case float64:
			switch right := right.(type) {
			case float64:
				return left > right
			}
		}
		runtimeError(b.Op, "Operands must be numbers.")
	case LESS_EQUAL:
		switch left := left.(type) {
		case float64:
			switch right := right.(type) {
			case float64:
				return left <= right
			}
		}
		runtimeError(b.Op, "Operands must be numbers.")
	case GREATER_EQUAL:
		switch left := left.(type) {
		case float64:
			switch right := right.(type) {
			case float64:
				return left >= right
			}
		}
		runtimeError(b.Op, "Operands must be numbers.")
	case EQUAL_EQUAL:
		return left == right
	case BANG_EQUAL:
		return left != right
	}
	loxError(b.Op, "not implemented")
	return nil
}

func (s *PrintStatement) Run() any {
	value := s.Value.Evaluate()
	if value == nil {
		fmt.Println("nil")
	} else {
		fmt.Println(value)
	}
	return nil
}

func (s *ExpressionStatement) Run() any {
	return s.Expr.Evaluate()
}

func (s *VarStatement) Run() any {
	var value any
	if s.Initializer != nil {
		value = s.Initializer.Evaluate()
	}
	env.Define(s.Name, value)
	return nil
}

func isTruthy(condition any) bool {
	switch condition := condition.(type) {
	case bool:
		return condition
	case nil:
		return false
	default:
		return true
	}
}

func (s *IfStatement) Run() any {
	condition := s.Condition.Evaluate()
	if isTruthy(condition) {
		s.ThenBranch.Run()
	} else if s.ElseBranch != nil {
		s.ElseBranch.Run()
	}
	return nil
}

func (b *Block) Run() any {
	prev := env
	env = NewEnvironent(prev)
	for _, statement := range b.Statements {
		statement.Run()
	}
	env = prev
	return nil
}

func (v *Variable) Evaluate() any {
	return env.Get(v.Name)
}

func (a *Assign) Evaluate() any {
	value := a.Value.Evaluate()
	env.Assign(a.Name.Name, value)
	return value
}
