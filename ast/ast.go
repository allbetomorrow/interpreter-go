package ast

import "interp/token"

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

type Decl interface {
	Node
	declNode()
}

type Program struct {
	Statement []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statement) > 0 {
		return p.Statement[0].TokenLiteral()
	}
	return ""
}

type DeclStatment struct {
	Token token.Token // the token.LEX_COLON token
	Name  *Identifier
	Type  Expression
}

func (ds *DeclStatment) statementNode() {}

func (ds *DeclStatment) TokenLiteral() string { return ds.Token.Literal }

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
