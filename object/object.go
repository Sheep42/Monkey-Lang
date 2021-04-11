package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/Sheep42/Monkey-Lang/ast"
)

const (
	StringObj      = "STRING"
	IntegerObj     = "INTEGER"
	NullObj        = "NULL"
	BooleanObj     = "BOOLEAN"
	ReturnValueObj = "RETURN_VALUE"
	ErrorObj       = "ERROR"
	FunctionObj    = "FUNCTION"
	BuiltinObj     = "BUILTIN"
	ArrayObj       = "ARRAY"
	HashObj        = "HASH"
)

type BuiltinFn func(args ...Object) Object

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

type Hashable interface {
	HashKey() HashKey
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

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ArrayObj }
func (a *Array) Inspect() string {

	var out bytes.Buffer

	elements := []string{}

	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()

}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HashObj }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}

	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type Builtin struct {
	Fn BuiltinFn
}

func (b *Builtin) Type() ObjectType { return BuiltinObj }
func (b *Builtin) Inspect() string  { return "builtin function" }
