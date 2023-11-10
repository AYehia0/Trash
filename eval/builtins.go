/*
Instead of attaching the builtins to the environment store, it's better to have it in a separete place.
*/
package eval

import (
	"trash/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Func: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErr(`Builtin "len": wrong number of args. got=%d, expected=1`, len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Int{Value: int64(len(arg.Value))}
			default:
				return newErr(`Builtin "len" doesn't take %s args`, arg.Type())
			}
		},
	},
}
