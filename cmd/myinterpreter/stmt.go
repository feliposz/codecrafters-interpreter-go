package main

type Stmt interface {
	Run() any
	Resolve()
}

type PrintStatement struct {
	Value Expr
}

type ExpressionStatement struct {
	Expr Expr
}

type VarStatement struct {
	Name        *Token
	Initializer Expr
}

type Block struct {
	Statements []Stmt
}

type IfStatement struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

type WhileStatement struct {
	Condition Expr
	Body      Stmt
}

type FunctionStatement struct {
	Name   *Token
	Params []*Token
	Body   []Stmt
}

type ReturnStatement struct {
	keyword *Token
	value   Expr
}
