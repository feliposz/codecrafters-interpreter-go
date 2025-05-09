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
		env.Define(param.Str, arguments[i])
	}
	result := runStatements(f.declaration.Body)
	env = prev
	if result, ok := result.(ReturnValue); ok {
		return result.Value
	}
	return nil
}

func (f *LoxFunction) Bind(instance *LoxInstance) *LoxFunction {
	instanceEnv := NewEnvironent(f.closure)
	instanceEnv.Define("this", instance)
	return &LoxFunction{f.declaration, instanceEnv}
}

type LoxClass struct {
	name    string
	methods map[string]*LoxFunction
}

func (c *LoxClass) FindMethod(str string) any {
	if method, ok := c.methods[str]; ok {
		return method
	}
	return nil
}

func (c *LoxClass) Arity() int {
	return 0
}

func (c *LoxClass) String() string {
	return c.name
}

func (c *LoxClass) Call(arguments []any) any {
	return &LoxInstance{c, make(map[string]any)}
}

type LoxInstance struct {
	class  *LoxClass
	fields map[string]any
}

func (i *LoxInstance) String() string {
	return i.class.String() + " instance"
}

func (i *LoxInstance) Get(name *Token) any {
	if value, ok := i.fields[name.Str]; ok {
		return value
	}
	if method := i.class.FindMethod(name.Str).(*LoxFunction); method != nil {
		return method.Bind(i)
	}
	runtimeError(name, "Undefined property '"+name.Str+"'.")
	return nil
}

func (i *LoxInstance) Set(name *Token, value any) {
	i.fields[name.Str] = value
}
