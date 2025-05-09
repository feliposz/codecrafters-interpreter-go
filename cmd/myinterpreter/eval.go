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
	switch value := s.Value.Evaluate().(type) {
	case nil:
		fmt.Println("nil")
	case float64:
		if value == float64(int(value)) {
			fmt.Printf("%.0f\n", value)
		} else {
			fmt.Printf("%g\n", value)
		}
	default:
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
		return s.ThenBranch.Run()
	} else if s.ElseBranch != nil {
		return s.ElseBranch.Run()
	}
	return nil
}

func (w *WhileStatement) Run() any {
	for isTruthy(w.Condition.Evaluate()) {
		if returnValue, ok := w.Body.Run().(ReturnValue); ok {
			return returnValue
		}
	}
	return nil
}

func (b *Block) Run() any {
	prev := env
	env = NewEnvironent(prev)
	result := runStatements(b.Statements)
	env = prev
	return result
}

func runStatements(statements []Stmt) any {
	for _, statement := range statements {
		if returnValue, ok := statement.Run().(ReturnValue); ok {
			return returnValue
		}
	}
	return nil
}

func (f *FunctionDeclaration) Run() any {
	function := &LoxFunction{f, env}
	env.Define(f.Name, function)
	return nil
}

func (r *ReturnStatement) Run() any {
	var value any
	if r.value != nil {
		value = r.value.Evaluate()
	}
	return ReturnValue{value}
}

type ReturnValue struct {
	Value any
}

func (c *ClassDeclaration) Run() any {
	env.Define(c.Name, nil)
	class := &LoxClass{c.Name.Str}
	env.Assign(c.Name, class)
	return nil
}

func (v *Variable) Evaluate() any {
	return lookUpVariable(v)
}

func (a *Assign) Evaluate() any {
	value := a.Value.Evaluate()
	assignVariable(a.Name, value)
	return value
}

func (c *Call) Evaluate() any {
	callee := c.callee.Evaluate()
	arguments := make([]any, len(c.arguments))
	for i, arg := range c.arguments {
		arguments[i] = arg.Evaluate()
	}
	if function, ok := callee.(LoxCallable); ok {
		if len(c.arguments) != function.Arity() {
			runtimeError(c.paren, fmt.Sprintf("Expected %d arguments but got %d.", function.Arity(), len(c.arguments)))
		}
		return function.Call(arguments)
	}
	runtimeError(c.paren, "Can only call functions and classes.")
	return nil
}

func (g *Get) Evaluate() any {
	object := g.object.Evaluate()
	if object, ok := object.(*LoxInstance); ok {
		return object.Get(g.name)
	}
	runtimeError(g.name, "Only instances have properties.")
	return nil
}

func (s *Set) Evaluate() any {
	object := s.object.Evaluate()
	if object, ok := object.(*LoxInstance); ok {
		value := s.value.Evaluate()
		object.Set(s.name, value)
		return nil
	}
	runtimeError(s.name, "Only instances have fields.")
	return nil
}
