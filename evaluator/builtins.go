package evaluator

import (
	"fmt"
	"monkey/object"
)

var builtins = map[string]*object.Builtin{
	"len": {

		Fn: func(args ...object.Object) object.Object {

			if len(args) != 1 {
				return newError(fmt.Sprintf("len: wrong number of args. expected=1. got=%d", len(args)))
			}

			switch arg := args[0].(type) {

			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}

			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}

			default:
				return newError(fmt.Sprintf("len: Unsupported argument. expected=STRING. got=%s", arg.Type()))

			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {

			if len(args) != 1 {
				return newError("first: Got wrong number of args. Expected=%d. Got=%d", 1, len(args))
			}

			if args[0].Type() != object.ArrayObj {
				return newError("first: No implementation for argument type %T. Expected=%s", args[0], object.ArrayObj)
			}

			arr := args[0].(*object.Array)

			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return Null

		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {

			if len(args) != 1 {
				return newError("last: Got wrong number of args. Expected=%d. Got=%d", 1, len(args))
			}

			if args[0].Type() != object.ArrayObj {
				return newError("last: No implementation for argument type %T. Expected=%s", args[0], object.ArrayObj)
			}

			arr := args[0].(*object.Array)

			if len(arr.Elements) > 0 {
				return arr.Elements[len(arr.Elements)-1]
			}

			return Null

		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {

			if len(args) != 1 {
				return newError("rest: Got wrong number of args. Expected=%d. Got=%d", 1, len(args))
			}

			if args[0].Type() != object.ArrayObj {
				return newError("rest: No implementation for argument type %T. Expected=%s", args[0], object.ArrayObj)
			}

			arr := args[0].(*object.Array)

			if length := len(arr.Elements); length > 0 {

				newElements := make([]object.Object, length-1, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Array{Elements: newElements}

			}

			return Null

		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {

			if len(args) != 2 {
				return newError("push: Got wrong number of args. Expected=%d. Got=%d", 2, len(args))
			}

			if args[0].Type() != object.ArrayObj {
				return newError("push: No implementation for argument type %T. Expected=%s", args[0], object.ArrayObj)
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)

			newElements := make([]object.Object, length+1, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]

			return &object.Array{Elements: newElements}

		},
	},
}
