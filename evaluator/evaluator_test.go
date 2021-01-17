package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {

	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"-5 - 8", -13},
		{"-(3 * 3 * 3 + 10)", -37},
		{"-3 + -3", -6},
		{"3 - 9", -6},
	}

	for _, tt := range tests {

		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)

	}

}

func TestEvalBooleanExpression(t *testing.T) {

	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 == 2", false},
		{"1 != 2", true},
		{"1 != 1", false},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"false == true", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"(1 > 2) != false", false},
	}

	for _, tt := range tests {

		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)

	}

}

func TestBangOperator(t *testing.T) {

	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {

		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)

	}
}

func TestIfElseExpressions(t *testing.T) {

	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if( true ) { 10 }", 10},
		{"if( false ) { 10 }", nil},
		{"if( 1 ) { 10 }", 10},
		{"if( 1 < 2 ) { 10 }", 10},
		{"if( 1 > 2 ) { 10 }", nil},
		{"if( 1 < 2 ) { 10 } else { 20 }", 10},
		{"if( 1 > 2 ) { 10 } else { 20 }", 20},
	}

	for _, tt := range tests {

		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObj(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {

	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 10; 9;", 10},
		{`
			if(10 > 1) {
				if(10 > 1) {
					return 10;
				}

				return 1;
			}
		`, 10},
	}

	for _, tt := range tests {

		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)

	}

}

func TestErrorHandling(t *testing.T) {

	tests := []struct {
		input       string
		expectedMsg string
	}{
		{
			"5 + true",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if(10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
				if(10 > 1) {

					if(10 > 1) {

						return true + false;

					}

				}

				return 1;
			`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
	}

	for _, tt := range tests {

		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)

		if !ok {

			t.Errorf("Evaluated result was not error message. Expected=%s. Got=%T(%+v)", tt.expectedMsg, evaluated, evaluated)
			continue

		}

		if errObj.Message != tt.expectedMsg {

			t.Errorf("Wrong error message. Expected=%s. Got=%s", tt.expectedMsg, errObj.Message)
			continue

		}

	}
}

func TestLetStatements(t *testing.T) {

	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {

		testIntegerObject(t, testEval(tt.input), tt.expected)

	}

}

func TestFunctionObject(t *testing.T) {

	input := "fn(x) { x + 2; }"

	evaluated := testEval(input)

	fn, ok := evaluated.(*object.Function)

	if !ok {
		t.Fatalf("object is not Function. Got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("Function has incorrect number of params. Expected=1. Got=%d", len(fn.Parameters))
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("Function parameter name is incorrect. Expected='x'. Got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("Function body is incorrect. Expected=%q. Got=%q", expectedBody, fn.Body.String())
	}

}

func TestCallStatements(t *testing.T) {

	tests := []struct {
		input    string
		expected int64
	}{
		{"let i = fn(x) { x; }; i(5);", 5},
		{"let i = fn(x) { return x; }; i(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"let i = fn(x) { x; }; i(5);", 5},
		{"fn(x) { x; }(5);", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}

}

func TestClosures(t *testing.T) {

	input := `
		let adder = fn(x) {
			fn(y) { x + y; }
		};

		let addTwo = adder(2);

		addTwo(2);
	`

	testIntegerObject(t, testEval(input), 4)

}

func testNullObj(t *testing.T, obj object.Object) bool {

	if obj != Null {

		t.Errorf("object is not Null. got=%T", obj)
		return false

	}

	return true

}

func testEval(input string) object.Object {

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)

}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {

	res, ok := obj.(*object.Integer)

	if !ok {

		t.Errorf("object is incorrect type. Expected=Integer. Got=%T", obj)
		return false

	}

	if res.Value != expected {

		t.Errorf("object has incorrect value. Expected=%d. Got=%d", expected, res.Value)
		return false

	}

	return true

}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {

	res, ok := obj.(*object.Boolean)

	if !ok {

		t.Errorf("object is incorrect type. Expected=Boolean. Got=%T", obj)
		return false

	}

	if res.Value != expected {

		t.Errorf("object has incorrect value. Expected=%t. Got=%t", expected, res.Value)
		return false

	}

	return true

}
