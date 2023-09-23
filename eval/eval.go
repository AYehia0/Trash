package eval

import (
	"fmt"
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
		return evalProgram(node)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.ReturnStatement:
		returnVal := Eval(node.ReturnValue)
		return &object.ReturnValue{Value: returnVal}

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

	case *ast.BlockStatement:
		return evalBlockStatement(node)

	case *ast.IfExpression:
		return evalIfExpression(node)
	}
	return nil
}

func newErr(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	conditionVal := Eval(ie.Condition)
	if isTruthy(conditionVal) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	}
	return NULL
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case FALSE:
		return false
	case TRUE:
		return true
	default:
		return true
	}
}
func evalInfixExpression(left object.Object, op string, right object.Object) object.Object {
	// integars
	switch {
	case left.Type() == object.INT_OBJ && right.Type() == object.INT_OBJ:
		return evalIntInfixExpression(left, op, right)

	case op == "==":
		return mapBool(left == right)
	case op == "!=":
		return mapBool(left != right)
	case left.Type() != right.Type():
		return newErr("Type mismatch: %s %s %s", left.Type(), op, right.Type())
	default:
		return newErr("Unknown operator: %s %s %s", left.Type(), op, right.Type())
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
		return newErr("Unkown Error: %s%s", op, right)
	}
}

func evalMinusOpExpression(right object.Object) object.Object {
	// check if the expression is bool
	if right.Type() != object.INT_OBJ {
		return newErr("Unknown operator: -%s", right.Type())
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

func evalProgram(prog *ast.Program) object.Object {
	var res object.Object

	for _, stmt := range prog.Statements {
		// The return value of the outer call to Eval is the return value of the last call
		res = Eval(stmt)
		switch res := res.(type) {
		case *object.ReturnValue:
			return res.Value // unpack
		case *object.Error:
			return res
		}
	}
	return res
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var res object.Object

	for _, stmt := range block.Statements {
		// The return value of the outer call to Eval is the return value of the last call
		res = Eval(stmt)
		if res != nil && res.Type() == object.RETURN_OBJ {
			return res
		}
	}
	return res
}

func mapBool(inp bool) *object.Bool {
	if inp {
		return TRUE
	}
	return FALSE
}
