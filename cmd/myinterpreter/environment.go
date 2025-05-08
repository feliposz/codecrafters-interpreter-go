package main

import "fmt"

type Environment struct {
	Enclosing *Environment
	Values    map[string]any
}

var env *Environment = NewGlobalEnvironment()

func NewGlobalEnvironment() *Environment {
	e := NewEnvironent(nil)
	e.Values["clock"] = &FunctionClock{}
	return e
}

func NewEnvironent(enclosing *Environment) *Environment {
	return &Environment{
		enclosing,
		make(map[string]any),
	}
}

func (e *Environment) Define(variable *Token, value any) {
	e.Values[variable.Str] = value
}

func (e *Environment) Assign(variable *Token, value any) {
	curr := e
	name := variable.Str
	for curr != nil {
		if _, found := curr.Values[name]; found {
			curr.Values[name] = value
			return
		}
		curr = curr.Enclosing
	}
	runtimeError(variable, fmt.Sprintf("Undefined variable '%s'.", name))
}

func (e *Environment) Get(variable *Token) any {
	curr := e
	name := variable.Str
	for curr != nil {
		if value, found := curr.Values[name]; found {
			return value
		}
		curr = curr.Enclosing
	}
	runtimeError(variable, fmt.Sprintf("Undefined variable '%s'.", name))
	return nil
}
