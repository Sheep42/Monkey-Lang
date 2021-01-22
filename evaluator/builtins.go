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

			default:
				return newError(fmt.Sprintf("len: Unsupported argument. expected=STRING. got=%s", arg.Type()))

			}
		},
	},
}
