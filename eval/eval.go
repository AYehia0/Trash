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

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)

		return evalInfixExpression(left, node.Operator, right)
	}
	return nil
}

func evalInfixExpression(left object.Object, op string, right object.Object) object.Object {
	// integars
	switch {
	case left.Type() == object.INT_OBJ && right.Type() == object.INT_OBJ:
		return evalIntInfixExpression(left, op, right)
	default:
		return NULL
	}
}

func evalIntInfixExpression(left object.Object, op string, right object.Object) object.Object {
	leftVal := left.(*object.Int).Value
	rightVal := right.(*object.Int).Value

	switch op {

	// integer expressions
	case "+":
		return &object.Int{Value: leftVal + rightVal}
	case "-":
		return &object.Int{Value: leftVal - rightVal}
	case "*":
		return &object.Int{Value: leftVal * rightVal}
	// TODO: return float values after impl float vars
	case "/":
		return &object.Int{Value: leftVal / rightVal}

	// boolean expressions
	case "<":
		return mapBool(leftVal < rightVal)
	case ">":
		return mapBool(leftVal > rightVal)
	case "==":
		return mapBool(leftVal == rightVal)
	case "!=":
		return mapBool(leftVal != rightVal)
	default:
		return NULL
	}
}
func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOpExpression(right)
	case "-":
		return evalMinusOpExpression(right)
	default:
		return NULL
	}
}

func evalMinusOpExpression(right object.Object) object.Object {
	// check if the expression is bool
	if right.Type() != object.INT_OBJ {
		return NULL
	}
	val := right.(*object.Int).Value
	return &object.Int{Value: -val}
}

func evalBangOpExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
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
