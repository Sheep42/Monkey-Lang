package object

import "fmt"

const (
	IntegerObj     = "INTEGER"
	NullObj        = "NULL"
	BooleanObj     = "BOOLEAN"
	ReturnValueObj = "RETURN_VALUE"
)

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
