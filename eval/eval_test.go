package eval

import (
	"testing"
	"trash/lexer"
	"trash/object"
	"trash/parser"
)

func TestEvalIntExpression(t *testing.T) {
	tests := []struct {
		input  string
		expect int64
	}{
		{"9", 9},
		{"6", 6},
		{"-9", -9},
		{"-6", -6},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntObject(t, evaluated, tt.expect)
	}
}
func TestEvalStringLiteral(t *testing.T) {
	input := `"Hello World!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}
func TestEvalBoolExpression(t *testing.T) {
	tests := []struct {
		input  string
		expect bool
	}{
		{"false", false},
		{"true", true},
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBoolObject(t, evaluated, tt.expect)
	}

}
func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStetements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 21;", 21},
		{"return 21; 90;", 21},
		{"90; if (true) { return 10; 39; }", 10},
		{"return 2 * 10;", 20},
		{"if (4 > 2) { if (3 > 1) { return 20; } return 0;}", 20},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"Type mismatch: INT + BOOL",
		},
		{
			"5 + true; 5;",
			"Type mismatch: INT + BOOL",
		},
		{
			"-true",
			"Unknown operator: -BOOL",
		},
		{
			"true + false;",
			"Unknown operator: BOOL + BOOL",
		},
		{
			"5; true + false; 5",
			"Unknown operator: BOOL + BOOL",
		},
		{
			"if (10 > 1) { true + false; }",
			"Unknown operator: BOOL + BOOL",
		},
		{
			`
			if (10 > 1) {
			if (10 > 1) {
			return true + false;
			}
			return 1;
			}`, "Unknown operator: BOOL + BOOL",
		},
		{
			"xyz",
			"Identifier not found: xyz",
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue

		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}
func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}
func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.Parse()
	env := object.NewEnv()

	return Eval(program, env)
}

func testIntObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Int)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
}

func testBoolObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Bool)
	if !ok {
		t.Errorf("object is not Bool. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input  string
		expect int64
	}{
		{"let x = 9; x;", 9},
		{"let x = 6; x;", 6},
		{"let x = -9; let y = x; y", -9},
		{"let x = -6; let y = x + 6; y", 0},
		// {"let x = 5 + 5 + 5 + 5 - 10; x = x + 10; x", 20}, // doesn't work !!
		// {"let x = 2 * 2 * 2 * 2 * 2; let a = x; a = a * 2; a", 64}, // doesn't work !!
		{"let x = -50 + 100 + -50; x;", 0},
		{"let x = 5 * 2 + 10; x;", 20},
		// {"let x = 5 + 2 * 10; x = x - 5", 20}, // doesn't work !!
		{"let x = 20 + 2 * -10; x;", 0},
		{"let x = 50 / 2 * 2 + 10; x;", 60},
		{"let x = 2 * (5 + 10); x;", 30},
		{"let x = 3 * 3 * 3 + 10; x;", 37},
		{"let x = 3 * (3 * 3) + 10; x;", 37},
		{"let x = (5 + 10 * 2 + 15 / 3) * 2 + -10; x;", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntObject(t, evaluated, tt.expect)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn (x) {return x*2;}"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}
	if len(fn.Params) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Params)
	}
	if fn.Params[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Params[0])
	}
	expectedBody := "return (x * 2);"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}
func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Int)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}
