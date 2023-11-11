/*
Instead of attaching the builtins to the environment store, it's better to have it in a separete place.
*/
package eval

import (
	"fmt"
	"os"
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
			case *object.List:
				return &object.Int{Value: int64(len(arg.Values))}
			default:
				return newErr(`Builtin "len" doesn't take %s args`, arg.Type())
			}
		},
	},
	// exit with status code
	"exit": {
		Func: func(args ...object.Object) object.Object {
			if len(args) > 1 {
				return newErr(`Builtin "len": wrong number of args. got=%d, expected=1`, len(args))
			}

			statusCode := 0
			if len(args) == 1 {
				statusArg, ok := args[0].(*object.Int)
				if !ok {
					return newErr("Builtin 'exit': argument must be an integer")
				}
				statusCode = int(statusArg.Value)
			}

			os.Exit(statusCode)
			return NULL
		},
	},
	"print": {
		Func: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
}
