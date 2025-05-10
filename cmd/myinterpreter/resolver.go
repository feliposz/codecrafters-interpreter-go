package main

type FunctionType uint8

const (
	FT_NONE FunctionType = iota
	FT_FUNCTION
	FT_INITIALIZER
	FT_METHOD
)

type ClassType uint8

const (
	CT_NONE ClassType = iota
	CT_CLASS
	CT_SUBCLASS
)

var currentFunction = FT_NONE
var currentClass = CT_NONE
var scopes []map[string]bool
var localsResolver = make(map[Expr]int, 0)

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

func resolveLocalVariable(variable Expr, token *Token) {
	for i := len(scopes) - 1; i >= 0; i-- {
		if _, found := scopes[i][token.Str]; found {
			depth := len(scopes) - 1 - i
			localsResolver[variable] = depth
			return
		}
	}
}

func lookUpVariable(variable Expr, token *Token) any {
	if distance, found := localsResolver[variable]; found {
		return env.GetAt(distance, token)
	} else {
		return globals.Get(token)
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

func (f *FunctionDeclaration) Resolve() {
	declare(f.Name)
	define(f.Name)
	resolveFunction(f, FT_FUNCTION)
}

func resolveFunction(f *FunctionDeclaration, functionType FunctionType) {
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
		if currentFunction == FT_INITIALIZER {
			loxError(r.keyword, "Can't return a value from an initializer.")
		}
		r.value.Resolve()
	}
}

func (w *WhileStatement) Resolve() {
	w.Condition.Resolve()
	w.Body.Resolve()
}

func (c *ClassDeclaration) Resolve() {
	enclosingClass := currentClass
	currentClass = CT_CLASS
	declare(c.Name)
	define(c.Name)
	if c.Superclass != nil {
		currentClass = CT_SUBCLASS
		if c.Name.Str == c.Superclass.Name.Str {
			loxError(c.Superclass.Name, "A class can't inherit from itself.")
			return
		}
		c.Superclass.Resolve()
	}
	if c.Superclass != nil {
		beginScope()
		currentScope()["super"] = true
	}
	beginScope()
	currentScope()["this"] = true
	for _, method := range c.Methods {
		functionType := FT_METHOD
		if method.Name.Str == "init" {
			functionType = FT_INITIALIZER
		}
		resolveFunction(method, functionType)
	}
	endScope()
	if c.Superclass != nil {
		endScope()
	}
	currentClass = enclosingClass
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

func (g *Get) Resolve() {
	g.object.Resolve()
}

func (s *Set) Resolve() {
	s.value.Resolve()
	s.object.Resolve()
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

func (t *This) Resolve() {
	if currentClass == CT_NONE {
		loxError(t.keyword, "Can't use 'this' outside of a class.")
		return
	}
	resolveLocalVariable(t, t.keyword)
}

func (s *Super) Resolve() {
	if currentClass == CT_NONE {
		loxError(s.keyword, "Can't use 'super' outside of a class.")
		return
	}
	if currentClass != CT_SUBCLASS {
		loxError(s.keyword, "Can't use 'super' in a class with no superclass.")
		return
	}
	resolveLocalVariable(s, s.keyword)
}
