package eval

import (
	"fmt"
	"trash/ast"
	"trash/object"
	"trash/token"
)

// instead of each time we encounter a new value we create one, instead we ref it.
var (
	NULL  = &object.Null{}
	TRUE  = &object.Bool{Value: true}
	FALSE = &object.Bool{Value: false}
)

func Eval(n ast.Node, env *object.Env) object.Object {
	switch node := n.(type) {

	// statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.ReturnStatement:
		returnVal := Eval(node.ReturnValue, env)
		if isErr(returnVal) {
			return returnVal
		}
		return &object.ReturnValue{Value: returnVal}

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Params: params, Body: body, Env: env}

	case *ast.CallExpression:
		function := Eval(node.Function, env)

		if isErr(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 {
			if isErr(args[0]) {
				return args[0]
			}
		}
		return getObjectFunction(function, args)
	case *ast.IntegerLiteral:
		return &object.Int{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.Boolean:
		return mapBool(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isErr(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isErr(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isErr(right) {
			return right
		}
		return evalInfixExpression(left, node.Operator, right)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.Identifier:
		return evalIdenterifer(node, env)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isErr(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	// TODO: really duplicate code from LetStatement
	case *ast.AssignExpression:
		val := Eval(node.Value, env)
		if isErr(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	}
	return nil
}

func evalExpressions(exps []ast.Expression, env *object.Env) []object.Object {
	var result []object.Object
	for _, obj := range exps {
		evaluted := Eval(obj, env)
		if isErr(evaluted) {
			return []object.Object{evaluted}
		}
		result = append(result, evaluted)
	}
	return result
}

func getObjectFunction(function object.Object, args []object.Object) object.Object {
	fn, ok := function.(*object.Function)
	if !ok {
		return newErr("%s is not a function", function.Type())
	}
	// mismatching args with given
	// TODO: add optional args
	if len(args) != len(fn.Params) {
		return newErr("Error: missing args to the function: %s", function.Inspect())
	}
	expandedEnv := expandFunctionEnv(fn, args)
	evaluated := Eval(fn.Body, expandedEnv)
	if returnVal, ok := evaluated.(*object.ReturnValue); ok {
		return returnVal.Value
	}
	return evaluated
}

func expandFunctionEnv(function *object.Function, args []object.Object) *object.Env {
	env := object.NewEnclosedEnv(function.Env)
	for i, param := range function.Params {
		env.Set(param.Value, args[i])
	}
	return env
}
func newErr(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func evalIdenterifer(node *ast.Identifier, env *object.Env) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newErr("Identifier not found: %s", node.Value)
	}
	return val
}

func evalIfExpression(ie *ast.IfExpression, env *object.Env) object.Object {
	conditionVal := Eval(ie.Condition, env)
	if isErr(conditionVal) {
		return conditionVal
	}
	if isTruthy(conditionVal) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
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

	// string concat
	case left.Type() == object.STR_OBJ && right.Type() == object.STR_OBJ:
		return evalStringConcat(left, op, right)

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

func evalStringConcat(left object.Object, op string, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	if op != token.CONCAT {
		return newErr("Unknown concat operator: '%s', use %s", op, token.CONCAT)
	}

	return &object.String{Value: leftVal + rightVal}
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
		return newErr("Unknown Error: %s%s", op, right)
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

func evalProgram(prog *ast.Program, env *object.Env) object.Object {
	var res object.Object

	for _, stmt := range prog.Statements {
		// The return value of the outer call to Eval is the return value of the last call
		res = Eval(stmt, env)
		switch res := res.(type) {
		case *object.ReturnValue:
			return res.Value // unpack
		case *object.Error:
			return res
		}
	}
	return res
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Env) object.Object {
	var res object.Object

	for _, stmt := range block.Statements {
		// The return value of the outer call to Eval is the return value of the last call
		res = Eval(stmt, env)
		if res != nil {
			resType := res.Type()
			if resType == object.RETURN_OBJ || resType == object.ERROR_OBJ {
				return res
			}
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
func isErr(obj object.Object) bool {
	if obj != nil {
		if obj.Type() == object.ERROR_OBJ {
			return true
		}
	}
	return false
}
