package parser

import (
	"fmt"
	"interp/ast"
	"interp/lexer"
	"interp/token"
	"strconv"
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
	p.registerPrefix(token.LEX_IDENT, p.parseIdentifier)
	p.registerPrefix(token.LEX_INT, p.parseIntegerLiteral)
	p.registerPrefix(token.LEX_MIN, p.parsePrefixExpression)
	p.registerPrefix(token.LEX_LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.KW_IF, p.parseIfExpression)
	p.registerPrefix(token.KW_BEGIN, p.parseBeginExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.LEX_PLUS, p.parseInfixExpression)
	p.registerInfix(token.LEX_MIN, p.parseInfixExpression)
	p.registerInfix(token.LEX_MULT, p.parseInfixExpression)
	p.registerInfix(token.LEX_DIV, p.parseInfixExpression)
	p.registerInfix(token.LEX_GT, p.parseInfixExpression)
	p.registerInfix(token.LEX_LT, p.parseInfixExpression)
	p.registerInfix(token.LEX_EQ, p.parseInfixExpression)
	p.registerInfix(token.LEX_NE, p.parseInfixExpression)
	p.registerInfix(token.LEX_LE, p.parseInfixExpression)
	p.registerInfix(token.LEX_GE, p.parseInfixExpression)

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

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.KW_THEN) {
		return nil
	}

	p.nextToken()
	expression.Consequence = p.parseBlockStatement()

	if p.curTokenIs(token.KW_ELSE) {
		p.nextToken()
		expression.Alternative = p.parseBlockStatement()
	}

	if !p.curTokenIs(token.KW_END) {
		return nil
	}

	return expression
}

func (p *Parser) parseBeginExpression() ast.Expression {
	expression := &ast.BeginExpression{Token: p.curToken}

	p.nextToken()

	expression.Block = p.parseBlockStatement()

	if !p.curTokenIs(token.KW_END) {
		return nil
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token: token.Token{
			Type:    "",
			Literal: "",
		},
	}
	block.Statements = []ast.Statement{}

	for !p.curTokenIs(token.KW_ELSE) && !p.curTokenIs(token.KW_END) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LEX_IDENT:
		if p.peekTokenIs(token.LEX_ASSIGN) {
			return p.parseAssignStatement()
		}

		if !p.peekTokenIs(token.LEX_COLON) {
			return p.parseExpressionStatement()
		}

		curTok := p.curToken
		p.nextToken()
		if p.peekTokenIsType() {
			return p.parseDeclStatement(curTok)
		} else {
			return p.parseMarkerStatement(curTok)
		}
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.LEX_SEMICOLON) {
		p.nextToken()
	}
	return stmt
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

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.LEX_RPAREN) {
		return nil
	}

	return exp
}
