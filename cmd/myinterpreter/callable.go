package main

import "time"

type LoxCallable interface {
	Call(arguments []any) any
	Arity() int
	String() string
}

type FunctionClock struct{}

func (c *FunctionClock) Arity() int {
	return 0
}

func (c *FunctionClock) String() string {
	return "<native fn>"
}

func (c *FunctionClock) Call(arguments []any) any {
	return float64(time.Now().Unix())
}
