package main

type Stmt interface {
	Run() any
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

type Function struct {
	Name   *Token
	Params []*Token
	Body   Stmt
}
