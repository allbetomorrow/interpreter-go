package parser

import (
	"interp/ast"
	"interp/lexer"
	"testing"
)

func TestDeclStatments(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		exprectedType      string
	}{
		{"x: integer;", "x", "integer"},
		{"foll: integer;", "foll", "integer"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]

		if stmt.TokenLiteral() != ":" {
			t.Errorf("s.TokenLiteral() is not \":\". got=%q.", stmt.TokenLiteral())
		}

		declStmt, ok := stmt.(*ast.DeclStatment)
		if !ok {
			t.Errorf("s not *ast.DeclStatment. got=%T", stmt)
		}

		if !testIdentifier(t, *declStmt.Name, tt.expectedIdentifier) {
			return
		}

		stmt_type := stmt.(*ast.DeclStatment).Type
		if !testTypeExpression(t, stmt_type, tt.exprectedType) {
			return
		}
	}
}

func testIdentifier(t *testing.T, s ast.Identifier, name string) bool {

	if s.Value != name {
		t.Errorf("declStmt.Name.Value not '%s'. got=%s", name, s.Value)
		return false
	}

	if s.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, s.TokenLiteral())
		return false
	}

	return true
}

func testTypeExpression(t *testing.T, got_type ast.Expression, exp_type string) bool {
	if got_type.TokenLiteral() != exp_type {
		t.Errorf("got_type.TokenLiteral() is not \"%s\". got=%s.", exp_type, got_type.TokenLiteral())
		return false
	}
	typeExpression := got_type.(*ast.Type)
	if typeExpression.Value != exp_type {
		t.Errorf("got_type.Value is not \"%s\". got=%s.", exp_type, typeExpression.Value)
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
