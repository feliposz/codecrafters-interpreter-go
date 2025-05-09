package main

type FunctionType uint8

const (
	FT_NONE FunctionType = iota
	FT_FUNCTION
)

var currentFunction = FT_NONE
var scopes []map[string]bool
var localsResolver = make(map[*Variable]int, 0)

func beginScope() {
	scopes = append(scopes, make(map[string]bool))
}

func endScope() {
	scopes = scopes[:len(scopes)-1]
}

func currentScope() map[string]bool {
	if len(scopes) == 0 {
		return nil
	}
	return scopes[len(scopes)-1]
}

func declare(token *Token) {
	if scope := currentScope(); scope != nil {
		if _, found := scope[token.Str]; found {
			loxError(token, "Already a variable with this name in this scope.")
		}
		scope[token.Str] = false
	}
}

func define(token *Token) {
	if scope := currentScope(); scope != nil {
		scope[token.Str] = true
	}
}

func resolveLocalVariable(variable *Variable, token *Token) {
	for i := len(scopes) - 1; i >= 0; i-- {
		if _, found := scopes[i][token.Str]; found {
			depth := len(scopes) - 1 - i
			localsResolver[variable] = depth
			return
		}
	}
}

func lookUpVariable(variable *Variable) any {
	if distance, found := localsResolver[variable]; found {
		return env.GetAt(distance, variable.Name)
	} else {
		return globals.Get(variable.Name)
	}
}

func assignVariable(variable *Variable, value any) {
	if distance, found := localsResolver[variable]; found {
		env.AssignAt(distance, variable.Name, value)
	} else {
		globals.Assign(variable.Name, value)
	}
}

func resolveStatements(statements []Stmt) {
	for _, statement := range statements {
		statement.Resolve()
	}
}

func (b *Block) Resolve() {
	beginScope()
	resolveStatements(b.Statements)
	endScope()
}

func (s *VarStatement) Resolve() {
	declare(s.Name)
	if s.Initializer != nil {
		s.Initializer.Resolve()
	}
	define(s.Name)
}

func (v *Variable) Resolve() {
	if scope := currentScope(); scope != nil {
		if initialized, found := scope[v.Name.Str]; found && !initialized {
			loxError(v.Name, "Can't read local variable in its own initializer.")
		}
	}
	resolveLocalVariable(v, v.Name)
}

func (a *Assign) Resolve() {
	a.Value.Resolve()
	resolveLocalVariable(a.Name, a.Name.Name)
}

func (f *FunctionStatement) Resolve() {
	declare(f.Name)
	define(f.Name)
	resolveFunction(f, FT_FUNCTION)
}

func resolveFunction(f *FunctionStatement, functionType FunctionType) {
	enclosingFunction := currentFunction
	currentFunction = functionType
	beginScope()
	for _, param := range f.Params {
		declare(param)
		define(param)
	}
	resolveStatements(f.Body)
	endScope()
	currentFunction = enclosingFunction
}

func (s *ExpressionStatement) Resolve() {
	s.Expr.Resolve()
}

func (s *IfStatement) Resolve() {
	s.Condition.Resolve()
	s.ThenBranch.Resolve()
	if s.ElseBranch != nil {
		s.ElseBranch.Resolve()
	}
}

func (s *PrintStatement) Resolve() {
	s.Value.Resolve()
}

func (r *ReturnStatement) Resolve() {
	if currentFunction == FT_NONE {
		loxError(r.keyword, "Can't return from top-level code.")
	}
	if r.value != nil {
		r.value.Resolve()
	}
}

func (w *WhileStatement) Resolve() {
	w.Condition.Resolve()
	w.Body.Resolve()
}

func (b *Binary) Resolve() {
	b.Left.Resolve()
	b.Right.Resolve()
}

func (c *Call) Resolve() {
	c.callee.Resolve()
	for _, arg := range c.arguments {
		arg.Resolve()
	}
}

func (g *Grouping) Resolve() {
	g.expr.Resolve()
}

func (l *Literal) Resolve() {
	// nothing to to
}

func (l *Logical) Resolve() {
	l.left.Resolve()
	l.right.Resolve()
}

func (u *Unary) Resolve() {
	u.Expr.Resolve()
}
