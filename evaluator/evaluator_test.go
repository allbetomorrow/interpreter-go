package evaluator

import (
	"interp/lexer"
	"interp/object"
	"interp/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5;", 5},
		{"10;", 10},
		{"-5;", -5},
		{"-10;", -10},
		{"5 + 5 + 5 + 5 - 10;", 10},
		{"2 * 2 * 2 * 2 * 2;", 32},
		{"-50 + 100 + -50;", 0},
		{"5 * 2 + 10;", 20},
		{"5 + 2 * 10;", 25},
		{"20 + 2 * -10;", 0},
		{"50 / 2 * 2 + 10;", 60},
		{"2 * (5 + 10);", 30},
		{"3 * 3 * 3 + 10;", 37},
		{"3 * (3 * 3) + 10;", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10;", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"1 < 2;", 1},
		{"1 > 2;", 0},
		{"1 < 1;", 0},
		{"1 > 1;", 0},
		{"1 = 1;", 1},
		{"1 <> 1;", 0},
		{"1 = 2;", 0},
		{"1 <> 2;", 1},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testCompObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if 1 then 10; end;", 10},
		{"if 1 < 2 then 10; end;", 10},
		{"if 1 > 2 then 10; end;", nil},
		{"if 1 > 2 then 10; else 20; end;", 20},
		{"if 1 < 2 then 10; else 20; end;", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestBeginExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"begin 5; end;", 5},
		{"begin a: integer; a := 10; a; end;", 10},
		{"begin end;", nil},
		{"begin a: integer; end;", nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"a: integer; a := 5; a;", 5},
		{"a: integer; a := 5 * 5; a;", 25},
		{"a: integer; a := 5; b: integer; b := a; b;", 5},
		{"a: integer; a := 5; b: integer; b := a; c: integer; c := a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestGoto(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]int64
	}{
		{
			input: `a: integer;
		begin
		 a := 5;
		 goto fs;
		 a := 10;
		end;
		fs:`,
			expected: map[string]int64{
				"a": 5,
			},
		},
		{
			input: `a: integer;
			begin
				a := 3;
				begin
					if a = 3 then
						a := a * 3;
						goto ttt;
					end;
					a := 10; 
				end;
				a := 15;
			end;
			a := 22;
			ttt:`,
			expected: map[string]int64{
				"a": 9,
			},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		Eval(program, env)

		for key, val := range tt.expected {
			env_var, ok := env.Get(key)
			if !ok {
				t.Fatalf("variable %s not exist", key)
			}

			integer_obj, ok := env_var.(*object.Integer)
			if !ok {
				t.Fatalf("env_var is not integer go %T", env_var)
			}

			if integer_obj.Value != val {
				t.Fatalf("value of %s is not %d, got %d", key, val, integer_obj.Value)
			}
		}

	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func testCompObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	// fmt.Printf("type of %q is %T\n", input, program.Statements[0])
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}

	return true
}
