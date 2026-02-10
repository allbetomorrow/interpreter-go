package ast

import (
	"interp/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&DeclStatment{
				Token: token.Token{Type: token.LEX_COLON, Literal: ":"},
				Name: &Identifier{
					Token: token.Token{Type: token.LEX_IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Type: &Type{
					Token: token.Token{Type: token.KW_INTEGER, Literal: "integer"},
					Value: "integer",
				},
			},
		},
	}

	if program.String() != "myVar: integer;" {
		t.Errorf("program.String wrong. got=%q", program.String())
	}
}
