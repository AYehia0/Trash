/*
The AST we're going to build consists of nodes which are connected to each other (a tree ofc)
Some of these nodes implement the Statement and some the Expression interface

These interfaces only contain dummy methods called statementNode and
expressionNode respectively. They are not strictly necessary but help by guiding the Go
compiler and possibly causing it to throw errors when we use a Statement where an Expression
should’ve been used, and vice versa

*/

package ast

import "trash/token"

// TokenLiteral() will be used only for debugging and testing
type Node interface {
	TokenLiteral() string
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

//  Identifers in other parts do produce values, e.g.: let x = valueProducingIdentifier;
// we’ll use Identifier here to represent the name in a variable binding and later reuse it, to represent an identifer as part of or as a complete expression.
type Identifier struct {
	Token token.Token
	Value string
}

type LetStatement struct {
	Token token.Token // token.IDENT
	Name  *Identifier // left side
	Value Expression  // right side
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
