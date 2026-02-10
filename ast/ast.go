package ast

import (
	"bytes"
	"interp/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

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

type MarkerStatement struct {
	Token  token.Token // the token.LEX_COLON token
	Marker *Identifier
}

func (ms *MarkerStatement) statementNode()       {}
func (ms *MarkerStatement) TokenLiteral() string { return ms.Token.Literal }
func (ms *MarkerStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ms.Marker.String())
	out.WriteString(ms.TokenLiteral())

	return out.String()
}

type DeclStatment struct {
	Token token.Token // the token.LEX_COLON token
	Name  *Identifier
	Type  *Type
	Value Expression
}

func (ds *DeclStatment) statementNode() {}

func (ds *DeclStatment) TokenLiteral() string { return ds.Token.Literal }

func (ds *DeclStatment) String() string {
	var out bytes.Buffer

	out.WriteString(ds.Name.String())
	out.WriteString(ds.TokenLiteral())
	out.WriteString(" ")
	if ds.Type != nil {
		out.WriteString(ds.Type.String())
	}
	if ds.Value != nil {
		out.WriteString(" = ")
		out.WriteString(ds.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

type AssignStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (as *AssignStatement) statementNode()       {}
func (as *AssignStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AssignStatement) String() string {
	var out bytes.Buffer

	out.WriteString(as.Name.String())
	out.WriteString(" := ")
	out.WriteString(as.Value.String())

	return out.String()
}

type Identifier struct {
	Token token.Token // the token.LEX_IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type Type struct {
	Token token.Token
	Value string
}

func (t *Type) expressionNode()      {}
func (t *Type) TokenLiteral() string { return t.Token.Literal }
func (t *Type) String() string       { return t.Value }
