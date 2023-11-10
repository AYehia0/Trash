/*
This file contains all the token used for building our language, which are then used to make the tokenizer/lexer.
*/
package token

// using a constant since our language is going to really limited and small, while it's better to use a hashmap
const (
	// identifiers: let IDENTIFER = 4;
	IDENT  = "IDENT" // add, foobar, x, y
	INT    = "INT"   // for numbers, only supports integers for now
	STRING = "STRING"

	// operators: +, *, /, -
	ASSIGN    = "="
	PLUS      = "+"
	CONCAT    = ":"
	NEG       = "-"
	MUL       = "*"
	DIV       = "/"
	BANG      = "!"
	LT        = "<"
	GT        = ">"
	EQUAL     = "=="
	NOT_EQUAL = "!="

	// delimiters: (, ), {, }, ;, ,
	SEMICOLON     = ";"
	COMMA         = ","
	LEFT_PAREN    = "("
	RIGHT_PAREN   = ")"
	LEFT_BRACE    = "{"
	RIGHT_BRACE   = "}"
	LEFT_BRACKET  = "["
	RIGHT_BRACKET = "]"

	// keywords
	FUNC   = "FUNCTION"
	LET    = "LET"
	IF     = "IF"
	ELSE   = "ELSE"
	RETURN = "RETURN"
	TRUE   = "TRUE"
	FALSE  = "FALSE"

	// special types
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

// seperating user-defined identifiers from langauge keywords
var keywords = map[string]TokenType{
	"fn":     FUNC,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
}

func LookIdentifier(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT // default for all user-defined
}
