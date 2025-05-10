package main

import "fmt"

type Environment struct {
	Enclosing *Environment
	Values    map[string]any
}

var globals *Environment = NewGlobalEnvironment()
var env *Environment = globals

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

func (e *Environment) Define(name string, value any) {
	e.Values[name] = value
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

func (e *Environment) GetAt(distance int, variable *Token) any {
	return e.Ancestor(distance).Get(variable)
}

func (e *Environment) GetByName(name string) any {
	curr := e
	for curr != nil {
		if value, found := curr.Values[name]; found {
			return value
		}
		curr = curr.Enclosing
	}
	runtimeError(nil, fmt.Sprintf("Undefined variable '%s'.", name))
	return nil
}

func (e *Environment) GetByNameAt(distance int, name string) any {
	return e.Ancestor(distance).GetByName(name)
}

func (e *Environment) AssignAt(distance int, name *Token, value any) {
	e.Ancestor(distance).Assign(name, value)
}

func (e *Environment) Ancestor(distance int) *Environment {
	result := e
	for i := 0; i < distance; i++ {
		result = result.Enclosing
	}
	return result
}
