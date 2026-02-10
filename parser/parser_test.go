package parser

import (
	"interp/ast"
	"interp/lexer"
	"testing"
)

func TestDeclStatments(t *testing.T) {
	input := `x: integer;
	foll: integer;`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParserProgram() returned nil")
	}

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statments does not contain 2 statments. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"foll"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testDeclStatment(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testDeclStatment(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != ":" {
		t.Errorf("s.TokenLiteral() is not \":\". got=%q.", s.TokenLiteral())
		return false
	}

	declStmt, ok := s.(*ast.DeclStatment)
	if !ok {
		t.Errorf("s not *ast.DeclStatment. got=%T", s)
		return false
	}

	if declStmt.Name.Value != name {
		t.Errorf("declStmt.Name.Value not '%s'. got=%s", name, declStmt.Name.Value)
		return false
	}

	if declStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, declStmt.Name)
		return false
	}

	return true
}
