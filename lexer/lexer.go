package lexer

import (
	"trash/token"
)

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

// check if the character is a word
// SUPPORT : ASCII only for now
func isLetter(ch byte) bool {
	if 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' {
		return true
	}
	return false
}

// check if the character is a digit
func isDigit(ch byte) bool {
	if '0' <= ch && ch <= '9' {
		return true
	}
	return false
}

func (l *Lexer) readInt() string {
	startPos := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[startPos:l.position]
}

func (l *Lexer) readAhead() byte {
	if l.nextPosition > len(l.input) {
		return 0
	}
	return l.input[l.nextPosition]
}

// read a complete word until the end, and update the position & nextPosition
// SUPPORT : ASCII only for now
func (l *Lexer) readIdentifer() string {

	// read until you find a " "
	startPos := l.position
	// can read letters with digits inside i think
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[startPos:l.position]
}

// skip tabs, whitespaces, ... etc
func (l *Lexer) skipSpaces() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

func (l *Lexer) NextToken() token.Token {
	var t token.Token

	// skip spaces
	l.skipSpaces()

	switch l.ch {
	// operators
	// TODO: refactor these branches
	case '=':
		if l.readAhead() == '=' {
			ch := l.ch
			l.readChar()
			t = token.Token{
				Type:    token.EQUAL,
				Literal: string(ch) + string(l.ch),
			}
		} else {
			t = newToken(token.ASSIGN, l.ch)
		}
	case '!':
		if l.readAhead() == '=' {
			ch := l.ch
			l.readChar()
			t = token.Token{
				Type:    token.NOT_EQUAL,
				Literal: string(ch) + string(l.ch),
			}
		} else {
			t = newToken(token.BANG, l.ch)
		}
	case '+':
		t = newToken(token.PLUS, l.ch)
	case '-':
		t = newToken(token.NEG, l.ch)
	case '*':
		t = newToken(token.MUL, l.ch)
	case '/':
		t = newToken(token.DIV, l.ch)
	case '>':
		t = newToken(token.GT, l.ch)
	case '<':
		t = newToken(token.LT, l.ch)
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

	// the default case is either: identifier, keyword, number or illeal
	default:
		// check i
		if isLetter(l.ch) {
			t.Literal = l.readIdentifer()
			t.Type = token.LookIdentifier(t.Literal)
			return t
		} else if isDigit(l.ch) {
			t.Literal = l.readInt()
			t.Type = token.INT
			return t
		} else {
			t = newToken(token.ILLEGAL, l.ch)
		}
	}

	// advance to the next character
	l.readChar()

	return t
}
