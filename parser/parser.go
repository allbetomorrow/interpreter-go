package parser

import (
	"interp/ast"
	"interp/lexer"
	"interp/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statement = []ast.Statement{}

	for p.curToken.Type != token.LEX_EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statement = append(program.Statement, stmt)

		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {

	}
}
