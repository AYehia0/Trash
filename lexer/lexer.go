package lexer

import "trash/token"

type Lexer struct {
	input        string
	position     int  // current position in the input file.
	nextPosition int  // current reading position in input
	ch           byte // the current position char
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// give us the next character and advance our position in the input string
func (l *Lexer) readChar() {
	// reset the current position character to "NUL"
	if l.nextPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.nextPosition]
	}
	l.position = l.nextPosition
	l.nextPosition++
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

func (l *Lexer) NextToken() token.Token {
	var t token.Token

	switch l.ch {
	// operators
	case '=':
		t = newToken(token.ASSIGN, l.ch)
	case '+':
		t = newToken(token.PLUS, l.ch)
	case '-':
		t = newToken(token.NEG, l.ch)
	case '*':
		t = newToken(token.MUL, l.ch)
	case '/':
		t = newToken(token.DIV, l.ch)
	// delimiters
	case ';':
		t = newToken(token.SEMICOLON, l.ch)
	case ',':
		t = newToken(token.COMMA, l.ch)
	case '(':
		t = newToken(token.LEFT_PAREN, l.ch)
	case ')':
		t = newToken(token.RIGHT_PAREN, l.ch)
	case '{':
		t = newToken(token.LEFT_BRACE, l.ch)
	case '}':
		t = newToken(token.RIGHT_BRACE, l.ch)

	// end of the file
	case 0:
		t.Literal = ""
		t.Type = token.EOF
	}

	// advance to the next character
	l.readChar()

	return t
}
