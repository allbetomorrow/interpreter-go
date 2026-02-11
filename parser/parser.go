package parser

import (
	"fmt"
	"interp/ast"
	"interp/lexer"
	"interp/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.LEX_EQ:   EQUALS,
	token.LEX_NE:   EQUALS,
	token.LEX_LT:   LESSGREATER,
	token.LEX_GT:   LESSGREATER,
	token.LEX_PLUS: SUM,
	token.LEX_MIN:  SUM,
	token.LEX_MULT: PRODUCT,
	token.LEX_DIV:  PRODUCT,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.LEX_INT, p.parseIntegerLiteral)
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekTokenIsType() bool {
	return p.peekToken.Type == token.KW_INTEGER || p.peekToken.Type == token.KW_REAL
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.LEX_EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)

		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LEX_IDENT:
		if p.peekTokenIs(token.LEX_ASSIGN) {
			return p.parseAssignStatement()
		}

		curTok := p.curToken

		if !p.expectPeek(token.LEX_COLON) {
			return nil
		}

		if p.peekTokenIsType() {
			return p.parseDeclStatement(curTok)
		} else {
			return p.parseMarkerStatement(curTok)
		}
	}
	return nil
}

func (p *Parser) parseDeclStatement(t token.Token) *ast.DeclStatment {
	stmt := &ast.DeclStatment{
		Token: p.curToken,
		Name: &ast.Identifier{
			Token: t,
			Value: t.Literal,
		},
	}

	p.nextToken()

	stmt.Type = &ast.Type{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if p.peekTokenIs(token.LEX_SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseAssignStatement() *ast.AssignStatement {
	name := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken()

	stmt := &ast.AssignStatement{Token: p.curToken, Name: name}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.LEX_SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseMarkerStatement(t token.Token) *ast.MarkerStatement {
	return nil
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()
	for !p.peekTokenIs(token.LEX_SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	return nil
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
