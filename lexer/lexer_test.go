package lexer

import (
	"interp/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `x: integer;
	read x;
	z: integer;
	z := 42;

	if x > y then
    write x;
	else
    write y;
	end;

	loop begin
    x := x - 1;
    goto done;
	end

	done:
	read skip space tab;
	>,=<><=>=<*/
	1343456B, 1343456b, 1343456C, 1343456c, 1343456D, 1343456H, 5ABH of`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LEX_IDENT, "x"},
		{token.LEX_COLON, ":"},
		{token.KW_INTEGER, "integer"},
		{token.LEX_SEMICOLON, ";"},
		{token.KW_READ, "read"},
		{token.LEX_IDENT, "x"},
		{token.LEX_SEMICOLON, ";"},
		{token.LEX_IDENT, "z"},
		{token.LEX_COLON, ":"},
		{token.KW_INTEGER, "integer"},
		{token.LEX_SEMICOLON, ";"},
		{token.LEX_IDENT, "z"},
		{token.LEX_ASSIGN, ":="},
		{token.LEX_INT, "42"},
		{token.LEX_SEMICOLON, ";"},
		{token.KW_IF, "if"},
		{token.LEX_IDENT, "x"},
		{token.LEX_GT, ">"},
		{token.LEX_IDENT, "y"},
		{token.KW_THEN, "then"},
		{token.KW_WRITE, "write"},
		{token.LEX_IDENT, "x"},
		{token.LEX_SEMICOLON, ";"},
		{token.KW_ELSE, "else"},
		{token.KW_WRITE, "write"},
		{token.LEX_IDENT, "y"},
		{token.LEX_SEMICOLON, ";"},
		{token.KW_END, "end"},
		{token.LEX_SEMICOLON, ";"},
		{token.KW_LOOP, "loop"},
		{token.KW_BEGIN, "begin"},
		{token.LEX_IDENT, "x"},
		{token.LEX_ASSIGN, ":="},
		{token.LEX_IDENT, "x"},
		{token.LEX_MIN, "-"},
		{token.LEX_INT, "1"},
		{token.LEX_SEMICOLON, ";"},
		{token.KW_GOTO, "goto"},
		{token.LEX_IDENT, "done"},
		{token.LEX_SEMICOLON, ";"},
		{token.KW_END, "end"},
		{token.LEX_IDENT, "done"},
		{token.LEX_COLON, ":"},
		{token.KW_READ, "read"},
		{token.KW_SKIP, "skip"},
		{token.KW_SPACE, "space"},
		{token.KW_TAB, "tab"},
		{token.LEX_SEMICOLON, ";"},
		{token.LEX_GT, ">"},
		{token.LEX_COMMA, ","},
		{token.LEX_EQ, "="},
		{token.LEX_NE, "<>"},
		{token.LEX_LE, "<="},
		{token.LEX_GE, ">="},
		{token.LEX_LT, "<"},
		{token.LEX_MULT, "*"},
		{token.LEX_DIV, "/"},
		{token.LEX_INT, "1343456B"},
		{token.LEX_COMMA, ","},
		{token.LEX_INT, "1343456b"},
		{token.LEX_COMMA, ","},
		{token.LEX_INT, "1343456C"},
		{token.LEX_COMMA, ","},
		{token.LEX_INT, "1343456c"},
		{token.LEX_COMMA, ","},
		{token.LEX_INT, "1343456D"},
		{token.LEX_COMMA, ","},
		{token.LEX_INT, "1343456H"},
		{token.LEX_COMMA, ","},
		{token.LEX_INT, "5ABH"},
		{token.KW_OF, "of"},
		{token.LEX_EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
