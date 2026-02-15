package parser

import (
	"fmt"
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

		if !testIdentifier(t, declStmt.Name, tt.expectedIdentifier) {
			return
		}

		stmt_type := declStmt.Type
		if !testTypeExpression(t, stmt_type, tt.exprectedType) {
			return
		}
	}
}

func testIdentifier(t *testing.T, s *ast.Identifier, name string) bool {

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

func TestAssignStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"x := 5;", "x", 5},
		{"foobar := y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]

		if stmt.TokenLiteral() != ":=" {
			t.Errorf("s.TokenLiteral() is not \":\". got=%q.", stmt.TokenLiteral())
		}

		assignStmt, ok := stmt.(*ast.AssignStatement)
		if !ok {
			t.Errorf("s not *ast.AssignStatement. got=%T", stmt)
		}

		if !testIdentifier(t, assignStmt.Name, tt.expectedIdentifier) {
			return
		}

		val := assignStmt.Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		ident, ok := exp.(*ast.Identifier)
		if !ok {
			t.Errorf("exp not *ast.Identifier. got=%T", exp)
			return false
		}
		return testIdentifier(t, ident, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.TokenLiteral())
		return false
	}

	return true
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"-15", "-", 15},
		{"-foo", "-", "foo"},
		{"(-f)", "-", "f"},
		{"(-15)", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}

	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		// {"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 = 5;", 5, "=", 5},
		{"5 <> 5;", 5, "<>", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar = barfoo;", "foobar", "=", "barfoo"},
		{"foobar <> barfoo;", "foobar", "<>", "barfoo"},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			for _, el := range program.Statements {
				fmt.Printf("%s\n", el.String())
			}
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue,
			tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 = 3 < 4",
			"((5 > 4) = (3 < 4))",
		},
		{
			"5 < 4 <> 3 > 4",
			"((5 < 4) <> (3 > 4))",
		},
		{
			"3 + 4 * 5 = 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) = ((3 * 1) + (4 * 5)))",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if x < y then x := y; end`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
			stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.AssignStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Name, "x") {
		return
	}

	val := consequence.Value
	if !testLiteralExpression(t, val, "y") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if x < y then
							x := y;
						else 
							y := x;
						end`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.AssignStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.AssignStatment. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Name, "x") {
		return
	}

	val := consequence.Value
	if !testLiteralExpression(t, val, "y") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements does not contain 1 statements. got=%d\n",
			len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.AssignStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.AssignStatement. got=%T",
			exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Name, "y") {
		return
	}

	alt_val := alternative.Value
	if !testLiteralExpression(t, alt_val, "x") {
		return
	}
}

func TestBeginExpression(t *testing.T) {
	input := `begin 
							x := y;
						end`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.BeginExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.BeginExpression. got=%T", stmt.Expression)
	}

	if len(exp.Block.Statements) != 1 {
		t.Errorf("block is not 1 statements. got=%d\n",
			len(exp.Block.Statements))
	}

	asignStmt, ok := exp.Block.Statements[0].(*ast.AssignStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.AssignStatement. got=%T",
			exp.Block.Statements[0])
	}

	if !testIdentifier(t, asignStmt.Name, "x") {
		return
	}

	alt_val := asignStmt.Value
	if !testLiteralExpression(t, alt_val, "y") {
		return
	}

}

func TestGoto(t *testing.T) {
	input := `goto self`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.GotoStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.GotoStatement. got=%T",
			program.Statements[0])
	}

	if !testIdentifier(t, stmt.Name, "self") {
		return
	}
}

// func TestParsion(t *testing.T) {
// 	input := `begin
//   x: integer;
//   count: integer;

//   begin
//     count := 10;

//     if count > 0 then
//       count := count - 1;
//     else
//       goto done;
//     end;
//   end;
// end;`

// 	l := lexer.New(input)
// 	p := New(l)
// 	program := p.ParseProgram()
// 	checkParserErrors(t, p)

// 	if len(program.Statements) != 1 {
// 		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
// 			1, len(program.Statements))
// 	}
// }

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) bool {

	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}
