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
	INDEX       // list[index]: the highest precedence
)

var precedences = map[token.TokenType]int{
	token.MUL:          PRODUCT,
	token.DIV:          PRODUCT,
	token.EQUAL:        EQUALS,
	token.NOT_EQUAL:    EQUALS,
	token.GT:           LESSGREATER,
	token.LT:           LESSGREATER,
	token.CONCAT:       SUM,
	token.PLUS:         SUM,
	token.NEG:          SUM,
	token.LEFT_PAREN:   CALL,
	token.LEFT_BRACKET: INDEX,
}

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
	// prefix
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.NEG, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBooleanExpression)
	p.registerPrefix(token.FALSE, p.parseBooleanExpression)
	p.registerPrefix(token.LEFT_BRACKET, p.parseListLiteral)

	// infix
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.CONCAT, p.parseInfixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.NEG, p.parseInfixExpression)
	p.registerInfix(token.MUL, p.parseInfixExpression)
	p.registerInfix(token.DIV, p.parseInfixExpression)
	p.registerInfix(token.EQUAL, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LEFT_PAREN, p.parseCallExpression)    // special one
	p.registerInfix(token.LEFT_BRACKET, p.parseIndexExpression) // special one

	// grouped
	// we only need to parse the left pren !!!
	p.registerPrefix(token.LEFT_PAREN, p.parseGroupedExpression)

	// other keywords if, else... etc
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNC, p.parseFunctionLiteral)

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

func (p *Parser) getPrecedence(tokenType token.TokenType) int {
	if p, ok := precedences[tokenType]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("No prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// check if we have a parsing function associated to the current token, if yes call it (parse it according to its type)
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.currToken.Type)
		return nil
	}
	// parse the prefix
	leftExp := prefix()

	for !p.TokenIs(p.peekToken, token.SEMICOLON) && precedence < p.getPrecedence(p.peekToken.Type) {
		infix := p.infixParseFns[p.peekToken.Type]

		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	id := &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}

	if p.TokenIs(p.peekToken, token.ASSIGN) {
		return p.parseAssignIdentifier(id)
	}
	return id
}

func (p *Parser) parseAssignIdentifier(id *ast.Identifier) ast.Expression {

	// x = x + 30
	ae := &ast.AssignExpression{
		Token: p.currToken,
		Name:  id,
	}

	if !p.expectNextToken(token.ASSIGN) {
		return nil
	}

	// advance after the assign =
	p.nextToken()

	ae.Value = p.parseExpression(LOWEST)

	if p.TokenIs(p.peekToken, token.SEMICOLON) {
		p.nextToken()
	}

	return ae
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

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	for p.TokenIs(p.peekToken, token.SEMICOLON) {
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

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.TokenIs(p.peekToken, token.SEMICOLON) {
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

// parse inside a list of expression : take an end : ])
func (p *Parser) parseListExpression(endToken token.TokenType) []ast.Expression {
	res := []ast.Expression{}

	if p.TokenIs(p.peekToken, endToken) {
		p.nextToken()
		return res
	}

	p.nextToken()
	res = append(res, p.parseExpression(LOWEST))

	for p.TokenIs(p.peekToken, token.COMMA) {
		// move the next value and the comma
		p.nextToken()
		p.nextToken()

		res = append(res, p.parseExpression(LOWEST))
	}

	if !p.expectNextToken(endToken) {
		return nil
	}

	return res
}

func (p *Parser) parseListLiteral() ast.Expression {
	list := &ast.ListLiteral{
		Token: p.currToken,
	}
	list.Values = p.parseListExpression(token.RIGHT_BRACKET)
	return list
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	ind := &ast.IndexExpression{
		Token: p.currToken,
		Left:  left,
	}

	p.nextToken()

	// the index should be the next token parsed
	ind.Index = p.parseExpression(LOWEST)

	if !p.expectNextToken(token.RIGHT_BRACKET) {
		return nil
	}

	if p.TokenIs(p.peekToken, token.ASSIGN) {
		p.nextToken()
		p.nextToken()

		ind.Value = p.parseExpression(LOWEST)
	}

	return ind
}

func (p *Parser) parseStringLiteral() ast.Expression {
	lit := &ast.StringLiteral{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}
	return lit
}
func (p *Parser) parsePrefixExpression() ast.Expression {
	pe := &ast.PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
	}

	p.nextToken()

	pe.Right = p.parseExpression(PREFIX)

	return pe
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	ie := &ast.InfixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
		Left:     left,
	}

	// get the current precedence level
	precedence := p.getPrecedence(p.currToken.Type)

	p.nextToken()

	ie.Right = p.parseExpression(precedence)

	return ie
}

func (p *Parser) parseBooleanExpression() ast.Expression {
	return &ast.Boolean{
		Token: p.currToken,
		Value: p.TokenIs(p.currToken, token.TRUE), // return true or false
	}
}
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectNextToken(token.RIGHT_PAREN) {
		return nil
	}
	return exp
}
func (p *Parser) parseIfExpression() ast.Expression {
	exp := &ast.IfExpression{
		Token: p.currToken,
	}

	if !p.expectNextToken(token.LEFT_PAREN) {
		return nil
	}

	p.nextToken()
	exp.Condition = p.parseExpression(LOWEST)

	if !p.expectNextToken(token.RIGHT_PAREN) {
		return nil
	}

	if !p.expectNextToken(token.LEFT_BRACE) {
		return nil
	}

	exp.Consequence = p.parseBlockStatement()

	if p.TokenIs(p.peekToken, token.ELSE) {
		p.nextToken()

		if !p.expectNextToken(token.LEFT_BRACE) {
			return nil
		}

		exp.Alternative = p.parseBlockStatement()
	}

	return exp
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := ast.BlockStatement{
		Token: p.currToken,
	}

	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.TokenIs(p.currToken, token.RIGHT_BRACE) && !p.TokenIs(p.currToken, token.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return &block
}

func (p *Parser) parseFunctionParams() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	// empty body
	if p.TokenIs(p.peekToken, token.RIGHT_PAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()
	ident := &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}
	// append the first arg
	identifiers = append(identifiers, ident)

	for p.TokenIs(p.peekToken, token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{
			Token: p.currToken,
			Value: p.currToken.Literal,
		}
		identifiers = append(identifiers, ident)
	}
	if !p.expectNextToken(token.RIGHT_PAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{
		Token: p.currToken,
	}

	if !p.expectNextToken(token.LEFT_PAREN) {
		return nil
	}

	// the params could be :
	//		myFunc(x, y, fn(x, y) { return x > y; });
	//		( ) --> empty
	//		(1 + 2, 3 * 8)
	lit.Parameters = p.parseFunctionParams()

	if !p.expectNextToken(token.LEFT_BRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{
		Token:    p.currToken,
		Function: function,
	}

	exp.Arguments = p.parseListExpression(token.RIGHT_PAREN)
	return exp
}
