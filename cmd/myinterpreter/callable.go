package main

import (
	"fmt"
	"time"
)

type LoxCallable interface {
	Call(arguments []any) any
	Arity() int
	String() string
}

type FunctionClock struct{}

func (f *FunctionClock) Arity() int {
	return 0
}

func (f *FunctionClock) String() string {
	return "<native fn>"
}

func (f *FunctionClock) Call(arguments []any) any {
	return float64(time.Now().Unix())
}

type LoxFunction struct {
	declaration *FunctionDeclaration
	closure     *Environment
}

func (f *LoxFunction) Arity() int {
	return len(f.declaration.Params)
}

func (f *LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.Name.Str)
}

func (f *LoxFunction) Call(arguments []any) any {
	prev := env
	env = NewEnvironent(f.closure)
	for i, param := range f.declaration.Params {
		env.Define(param, arguments[i])
	}
	result := runStatements(f.declaration.Body)
	env = prev
	if result, ok := result.(ReturnValue); ok {
		return result.Value
	}
	return nil
}

type LoxClass struct {
	name string
}

func (c *LoxClass) Arity() int {
	return 0
}

func (c *LoxClass) String() string {
	return c.name
}

func (c *LoxClass) Call(arguments []any) any {
	return &LoxInstance{c}
}

type LoxInstance struct {
	class *LoxClass
}

func (i *LoxInstance) String() string {
	return i.class.String() + " instance"
}
