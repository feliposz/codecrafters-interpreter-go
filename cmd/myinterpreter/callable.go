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
	declaration   *FunctionDeclaration
	closure       *Environment
	isInitializer bool
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
	if f.isInitializer {
		return f.closure.Values["this"]
	}
	if result, ok := result.(ReturnValue); ok {
		return result.Value
	}
	return nil
}

func (f *LoxFunction) Bind(instance *LoxInstance) *LoxFunction {
	instanceEnv := NewEnvironent(f.closure)
	instanceEnv.Define("this", instance)
	return &LoxFunction{f.declaration, instanceEnv, f.isInitializer}
}

type LoxClass struct {
	name       string
	superclass *LoxClass
	methods    map[string]*LoxFunction
}

func (c *LoxClass) FindMethod(name string) *LoxFunction {
	if method, ok := c.methods[name]; ok {
		return method
	}
	if c.superclass != nil {
		return c.superclass.FindMethod(name)
	}
	return nil
}

func (c *LoxClass) Arity() int {
	if initializer := c.FindMethod("init"); initializer != nil {
		return initializer.Arity()
	}
	return 0
}

func (c *LoxClass) String() string {
	return c.name
}

func (c *LoxClass) Call(arguments []any) any {
	instance := &LoxInstance{c, make(map[string]any)}
	if initializer := c.FindMethod("init"); initializer != nil {
		initializer.Bind(instance).Call(arguments)
	}
	return instance
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
	if method := i.class.FindMethod(name.Str); method != nil {
		return method.Bind(i)
	}
	runtimeError(name, "Undefined property '"+name.Str+"'.")
	return nil
}

func (i *LoxInstance) Set(name *Token, value any) {
	i.fields[name.Str] = value
}
