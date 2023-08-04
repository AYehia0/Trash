/*
A Pratt parser’s main idea is the association of parsing functions (which Pratt calls “semantic
code”) with token types. Whenever this token type is encountered, the parsing functions are
called to parse the appropriate expression and return an AST node that represents it. Each
token type can have up to two parsing functions associated with it, depending on whether the
token is found in a prefx or an infx position.
*/
package parser

import (
	"fmt"
	"strconv"
	"trash/ast"
	"trash/lexer"
	"trash/token"
)

type (
	prefixParseFn func() ast.Expression               // --x
	infixParseFn  func(ast.Expression) ast.Expression // 6 * 9 (left side that's being parsed)
)

// precedences of the parser
/*
What we want out of these constants is to later be able to answer:
- Does the * operator have a higher precedence than the == operator?
- Does a prefx operator have a higher preference than a call expression?
*/
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // < >
	SUM         // +
	PRODUCT     // *
	PREFIX      // -x or !x
	CALL        // myfunc(x)
)

type Parser struct {
	l         *lexer.Lexer
	currToken token.Token
	peekToken token.Token
	errors    []string

	// check if the appropriate map (infx or prefx) has a parsing function associated with currToken.Type
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{
		l:      lexer,
		errors: []string{},
	}
	// associate tokens to the parser
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)

	// read 2 tokens so current and next token are set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(tk token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", tk, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// the parser returns the AST
func (p *Parser) Parse() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currToken.Type != token.EOF {
		stmt := p.parseStatement()

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		// advance to the next token
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// check if we have a parsing function associated to the current token, if yes call it (parse it according to its type)
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currToken.Type]

	if prefix == nil {
		return nil
	}
	// parse the prefix
	leftExp := prefix()

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}
}

// let <identifier> = <expression>
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{
		Token: p.currToken,
	}

	if !p.expectNextToken(token.IDENT) {
		return nil
	}
	// the left side
	stmt.Name = &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}

	if !p.expectNextToken(token.ASSIGN) {
		return nil
	}
	// TODO: finish the expression, for now skip till you find a semicolon
	// stmt.Value = ast.Expression{
	//
	// }
	// parse to the end
	for !p.TokenIs(p.currToken, token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) TokenIs(token token.Token, tt token.TokenType) bool {
	return token.Type == tt
}

func (p *Parser) expectNextToken(tokenType token.TokenType) bool {
	if p.TokenIs(p.peekToken, tokenType) {
		p.nextToken()
		return true
	}
	p.peekError(tokenType)
	return false
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {

	stmt := &ast.ReturnStatement{
		Token: p.currToken,
	}

	p.nextToken()

	// TODO: finish the expression, for now skip till you find a semicolon
	if !p.TokenIs(p.currToken, token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Token: p.currToken,
	}
	stmt.Expression = p.parseExpression(LOWEST)

	// optional semicolons
	if p.TokenIs(p.peekToken, token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{
		Token: p.currToken,
	}

	intValue, err := strconv.ParseInt(p.currToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("Couldn't parse %s as integer", p.currToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = intValue

	return lit
}
