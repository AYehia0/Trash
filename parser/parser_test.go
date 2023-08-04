package parser

import (
	"fmt"
	"testing"
	"trash/ast"
	"trash/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 8;
	let foobar = 29393;
	let foobar = x + y;
	`
	l := lexer.New(input)
	p := New(l)

	ast := p.Parse()
	checkParserErrors(t, p)
	if ast == nil {
		t.Fatalf("Parse() returned an error;")
	}

	tests := []struct {
		expectedIndentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
		{"foobar"},
	}

	if len(ast.Statements) < 3 {
		t.Fatalf("Expected to have length of %d, got %d", len(tests), len(ast.Statements))
	}

	for i, tt := range tests {
		stmt := ast.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIndentifier) {
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

func TestPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
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

		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
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
