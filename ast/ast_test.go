package ast

import (
	"testing"
	"trash/token"
)

// check if the AST == String() (AST after forming)
func TestString(t *testing.T) {
	// let x = y;
	ast := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				// right side
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "x"},
					Value: "x",
				},
				// left side of =
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "y"},
					Value: "y",
				},
			},
		},
	}
	if ast.String() != "let x = y;" {
		t.Errorf("program.String() returned wrong value, expected: 'let x = y;' got %s",
			ast.String())
	}
}
