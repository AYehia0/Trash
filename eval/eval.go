package eval

import (
	"trash/ast"
	"trash/object"
)

// instead of each time we encounter a new value we create one, instead we ref it.
var (
	NULL  = &object.Null{}
	TRUE  = &object.Bool{Value: true}
	FALSE = &object.Bool{Value: false}
)

func Eval(n ast.Node) object.Object {
	switch node := n.(type) {

	// statements
	case *ast.Program:
		return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.IntegerLiteral:
		return &object.Int{Value: node.Value}

	case *ast.Boolean:
		return mapBool(node.Value)
	}
	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var res object.Object

	for _, stmt := range stmts {
		// The return value of the outer call to Eval is the return value of the last call
		res = Eval(stmt)
	}
	return res
}

func mapBool(inp bool) *object.Bool {
	if inp {
		return TRUE
	}
	return FALSE
}
