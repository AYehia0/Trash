package parser

import (
	"fmt"
	"testing"
	"trash/ast"
	"trash/lexer"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.Parse()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}
		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}
func testLetStatement(t *testing.T, stat ast.Statement, name string) bool {
	// checking the token literal
	if stat.TokenLiteral() != "let" {
		t.Fatalf("Expected statement.TokenLiteral to be let, got %s", stat.TokenLiteral())
	}

	letStat, ok := stat.(*ast.LetStatement)
	if !ok {
		t.Fatalf("stat not *stat.LetStatement")
		return false
	}

	if letStat.Name.Value != name {
		t.Fatalf("letStat.Name.Value not '%s'. got=%s", name, letStat.Name.Value)
		return false
	}

	if letStat.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStat.Name)
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("Parser has %d errors", len(errors))

	// printing
	for _, msg := range errors {
		t.Errorf("Parser error : %s", msg)
	}
	t.FailNow()
}

func TestReturnStatement(t *testing.T) {
	input := `
		return 69;
		return 420;
		return 2024;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Errorf("Expected to have 3 parsed return statemets, got %d", len(program.Statements))
	}

	//checking the tokens
	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
	}

}
func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.New(input)
	p := New(l)
	program := p.Parse()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected Program to have 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
			ident.TokenLiteral())
	}

}

func TestIntegerLiteralExpression(t *testing.T) {

	input := "9;"

	l := lexer.New(input)
	p := New(l)
	program := p.Parse()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected Program to have 1 statement, got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)

	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 9 {
		t.Errorf("literal.Value not %d. got=%d", 9, literal.Value)
	}
	if literal.TokenLiteral() != "9" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "9",
			literal.TokenLiteral())
	}
}

func TestStringLiteralExpression(t *testing.T) {

	input := `"string, Ma boi"`

	l := lexer.New(input)
	p := New(l)
	program := p.Parse()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected Program to have 1 statement, got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	literal, ok := stmt.Expression.(*ast.StringLiteral)

	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != "string, Ma boi" {
		t.Errorf("literal.Value not %s. got=%s", "string, Ma boi", literal.Value)
	}
}
func TestPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}
	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.Parse()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.TokenLiteral())
		return false
	}

	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}
	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.Parse()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)

		testInfixExpression(t, exp, tt.leftValue, tt.operator, tt.rightValue)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.Parse()
		checkParserErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	input := `
		true;
		false;
	`

	expected := []string{"true", "false"}
	l := lexer.New(input)
	p := New(l)
	program := p.Parse()

	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("Expected Program to have 2 statement, got=%d", len(program.Statements))
	}

	for i, stmt := range program.Statements {
		stmtExp := stmt.(*ast.ExpressionStatement)
		boolStmt, ok := stmtExp.Expression.(*ast.Boolean)
		if !ok {
			t.Errorf("stmt not *ast.Boolean got=%T", stmt)
			continue
		}
		if boolStmt.TokenLiteral() != expected[i] {
			t.Errorf("boolStmt.TokenLiteral not 'true or false', got %q",
				boolStmt.TokenLiteral())
		}
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, val string) bool {
	ident, ok := exp.(*ast.Identifier)

	if !ok {
		t.Fatalf("exp not *ast.Identifier, got v=%T", ident)
		return false
	}

	if ident.Value != val {
		t.Fatalf("ident.value != %s, got value=%s", val, ident.Value)
		return false
	}

	if ident.TokenLiteral() != val {
		t.Fatalf("ident.TokenLiteral != %s, got value=%s", val, ident.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testInfixExpression(
	t *testing.T,
	exp ast.Expression,
	left interface{},
	operator string,
	right interface{},
) bool {

	opExp, ok := exp.(*ast.InfixExpression)

	if !ok {
		t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}
	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}
	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			value, bo.TokenLiteral())
		return false
	}
	return true
}

func TestIfExpression(t *testing.T) {

	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
			stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}
	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
			stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	alt, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, alt.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	function, ok := stmt.Expression.(*ast.FunctionLiteral)

	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Expression)
	}
	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(function.Body.Statements))
	}
	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
			function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}
	for _, tt := range tests {

		l := lexer.New(tt.input)
		p := New(l)
		program := p.Parse()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(function.Parameters))
		}
		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

// extend it later
func TestParseListLiteral(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	l := lexer.New(input)
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ListLiteral)
	if !ok {
		t.Fatalf("exp not ast.ListLiteral. got=%T", stmt.Expression)
	}
	if len(array.Values) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Values))
	}
	testIntegerLiteral(t, array.Values[0], 1)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"
	l := lexer.New(input)
	p := New(l)

	program := p.Parse()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}
	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}
	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"
	l := lexer.New(input)
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}
	if !testIdentifier(t, exp.Function, "add") {
		return
	}
	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`
	l := lexer.New(input)
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}
	if len(hash.Store) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Store))
	}
	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	for key, value := range hash.Store {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}
		expectedValue := expected[literal.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

// arrhhhggggg!!!!!
func TestParsingEmptyHashLiteral(t *testing.T) {
	input := "{}"
	l := lexer.New(input)
	p := New(l)

	program := p.Parse()
	checkParserErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}
	if len(hash.Store) != 0 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Store))
	}
}
func TestParsingHashLiteralsWithExpressions(t *testing.T) {

	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`
	l := lexer.New(input)
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}
	for key, value := range hash.Store {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}
		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
			continue
		}
		testFunc(value)
	}
}
