/*
The AST we're going to build consists of nodes which are connected to each other (a tree ofc)
Some of these nodes implement the Statement and some the Expression interface

These interfaces only contain dummy methods called statementNode and
expressionNode respectively. They are not strictly necessary but help by guiding the Go
compiler and possibly causing it to throw errors when we use a Statement where an Expression
should’ve been used, and vice versa

*/

package ast

import (
	"bytes"
	"trash/token"
)

// TokenLiteral() will be used only for debugging and testing
type Node interface {
	TokenLiteral() string
	String() string // for printing the ast, debugging only
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// the program which is the whole parsed code expressed as list of Statements
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

//	Identifers in other parts do produce values, e.g.: let x = valueProducingIdentifier;
//
// we’ll use Identifier here to represent the name in a variable binding and later reuse it, to represent an identifer as part of or as a complete expression.
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string {
	return i.Value
}

type LetStatement struct {
	Token token.Token // token.IDENT
	Name  *Identifier // left side
	Value Expression  // right side
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")

	return out.String()
}

// return x;
// return 69;
// return someFunction(9, someVar);
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

// expression statements, only expression line
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}
