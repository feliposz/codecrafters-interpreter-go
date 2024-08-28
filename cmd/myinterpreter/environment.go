package main

import "fmt"

type Environment struct {
	Enclosing *Environment
	Values    map[string]any
}

var env *Environment = NewEnvironent(nil)

func NewEnvironent(enclosing *Environment) *Environment {
	return &Environment{
		enclosing,
		make(map[string]any),
	}
}

func (e *Environment) Set(variable *Token, value any) {
	e.Values[variable.String()] = value
}

func (e *Environment) Get(variable *Token) any {
	curr := e
	name := variable.String()
	for curr != nil {
		if value, found := curr.Values[name]; found {
			return value
		}
		curr = curr.Enclosing
	}
	runtimeError(variable, fmt.Sprintf("Undefined variable '%s'.", name))
	return nil
}
