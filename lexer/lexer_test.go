package lexer

import (
	"testing"
	"trash/token"
)

func TestNextToken(t *testing.T) {
	input := `
		let six = 6;
		let se7enty = 70;
		let add = fn(x, y) {
			x + y;
		};
		let result = add(six, se7enty);
	`
	expectedTests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "six"},
		{token.ASSIGN, "="},
		{token.INT, "6"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "se7enty"},
		{token.ASSIGN, "="},
		{token.INT, "70"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNC, "fn"},
		{token.LEFT_PAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RIGHT_PAREN, ")"},
		{token.LEFT_BRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RIGHT_BRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LEFT_PAREN, "("},
		{token.IDENT, "six"},
		{token.COMMA, ","},
		{token.IDENT, "se7enty"},
		{token.RIGHT_PAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range expectedTests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
