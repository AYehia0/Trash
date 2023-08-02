package parser

import (
	"fmt"
	"trash/ast"
	"trash/lexer"
	"trash/token"
)

type Parser struct {
	l         *lexer.Lexer
	currToken token.Token
	peekToken token.Token
	errors    []string
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{
		l:      lexer,
		errors: []string{},
	}
	// read 2 tokens so current and next token are set
	p.nextToken()
	p.nextToken()

	return p
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
		return nil
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
	if p.currToken.Type != token.SEMICOLON {
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
