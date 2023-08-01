package parser

import (
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
