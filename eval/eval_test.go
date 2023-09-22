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
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBoolObject(t, evaluated, tt.expect)
	}

}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.Parse()

	return Eval(program)
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
