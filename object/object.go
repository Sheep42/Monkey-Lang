package object

import (
	"bytes"
	"fmt"
	"monkey/ast"
	"strings"
)

const (
	StringObj      = "STRING"
	IntegerObj     = "INTEGER"
	NullObj        = "NULL"
	BooleanObj     = "BOOLEAN"
	ReturnValueObj = "RETURN_VALUE"
	ErrorObj       = "ERROR"
	FunctionObj    = "FUNCTION"
)

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {

	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}

}

func NewEnclosedEnvironment(outer *Environment) *Environment {

	env := NewEnvironment()
	env.outer = outer
	return env

}

func (e *Environment) Get(name string) (Object, bool) {

	obj, ok := e.store[name]

	if !ok && e.outer != nil {

		obj, ok = e.outer.Get(name)

	}

	return obj, ok

}

func (e *Environment) Set(name string, val Object) Object {

	e.store[name] = val
	return val

}

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Null struct{}

func (n *Null) Type() ObjectType { return NullObj }
func (n *Null) Inspect() string  { return fmt.Sprintf("null") }

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return IntegerObj }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return StringObj }
func (s *String) Inspect() string  { return s.Value }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BooleanObj }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Type() ObjectType { return ReturnValueObj }
func (r *ReturnValue) Inspect() string  { return r.Value.Inspect() }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ErrorObj }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FunctionObj }
func (f *Function) Inspect() string {

	var out bytes.Buffer

	params := []string{}

	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()

}
