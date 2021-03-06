package evaluator

import (
	"testing"

	"github.com/Sheep42/Monkey-Lang/lexer"
	"github.com/Sheep42/Monkey-Lang/object"
	"github.com/Sheep42/Monkey-Lang/parser"
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

func TestEvalStringLiteral(t *testing.T) {

	tests := []struct {
		input    string
		expected string
	}{
		{`"Hello World"`, "Hello World"},
		{`'Hello World'`, "Hello World"},
	}

	for _, tt := range tests {

		eval := testEval(tt.input)
		str, ok := eval.(*object.String)

		if !ok {
			t.Fatalf("Evaluated obj was incorrect type. Expected=\"*object.String\". Got=\"%T\"", eval)
		}

		if str.Value != tt.expected {
			t.Errorf("String has incorrect value. Expected=%q. Got=%q", tt.expected, str.Value)
		}

	}

}

func TestStringConcatenation(t *testing.T) {

	tests := []struct {
		input    string
		expected string
	}{
		{`"Hello" + " " + "World"`, "Hello World"},
		{`'Hello' + ' ' + 'World'`, "Hello World"},
		{`"Hello" + ' ' + 'World'`, "Hello World"},
	}

	for _, tt := range tests {

		eval := testEval(tt.input)

		str, ok := eval.(*object.String)

		if !ok {
			t.Fatalf("Object is incorrect type. Expected=\"String\". Got=\"%T\"", eval)
		}

		if str.Value != tt.expected {
			t.Errorf("String has incorrect value. Expected=%q. Got=%q", tt.expected, str.Value)
		}

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
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
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
		{
			`{"name": "monkey"}[fn(x) { x }]`,
			`Invalid HashKey: "fn(x) {\nx\n}". Type "FUNCTION" is unsupported.`,
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

func TestBuiltinFns(t *testing.T) {

	tests := []struct {
		input    string
		expected interface{}
	}{

		{`len("")`, 0},
		{`len("four")`, 4},
		{`len('four')`, 4},
		{`len("hello world")`, 11},
		{`len([1, 2])`, 2},
		{`first([1, 2, 3])`, 1},
		{`last([1, 2, 3])`, 3},
		{`rest([1, 2, 3])[0]`, 2},
		{`push([], 1)[0]`, 1},
		{`len(1)`, "len: Unsupported argument. expected=STRING. got=INTEGER"},
		{`len("one", "two")`, "len: wrong number of args. expected=1. got=2"},
	}

	for _, tt := range tests {

		eval := testEval(tt.input)

		switch expected := tt.expected.(type) {

		case int:
			testIntegerObject(t, eval, int64(expected))
		case string:
			errObj, ok := eval.(*object.Error)

			if !ok {
				t.Errorf("object is not Error. got=%T", eval)
				continue
			}

			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q. got=%q", expected, errObj.Message)
			}

		}

	}

}

func TestArrayLiterals(t *testing.T) {

	input := "[1, 2 * 2, 3 + 4]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)

	if !ok {
		t.Fatalf("Object was incorrect type. Expected=%s. Got=%T", "ast.Array", evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("Array has incorrect number of elements. Expected=%d. Got=%d", 3, len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 7)

}

func TestArrayIndexExpressions(t *testing.T) {
	testCases := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"[1, 2, 3][1 + 1]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"let arr = [1, 2, 3]; arr[0];",
			1,
		},
		{
			"let arr = [1, 2, 3]; arr[0] + arr[1] + arr[2];",
			6,
		},
		{
			"let arr = [1, 2, 3]; let i = arr[0]; arr[i];",
			2,
		},
		{
			"[1, 2, 3][3];",
			nil,
		},
		{
			"[1, 2, 3][-1];",
			nil,
		},
	}

	for _, tt := range testCases {

		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)

		if ok {

			testIntegerObject(t, evaluated, int64(integer))

		} else {

			testNullObj(t, evaluated)

		}

	}

}

func TestHashLiteral(t *testing.T) {
	input := `let two = "two";
		{
			"one": 10 - 9,
			two: 1 + 1,
			"thr" + "ee": 6 / 2,
			4: 4,
			true: 5,
			false: 6
		}
	`

	evaled := testEval(input)
	res, ok := evaled.(*object.Hash)

	if !ok {
		t.Fatalf("Eval returned unexpected type. Expected=%s. Got=%T", "Hash", evaled)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		(True.HashKey()):                           5,
		(False.HashKey()):                          6,
	}

	if len(res.Pairs) != len(expected) {
		t.Fatalf("Hash has incorrect number of pairs. Expected=%d. Got=%d", len(expected), len(res.Pairs))
	}

	for expectedKey, expectedVal := range expected {

		pair, ok := res.Pairs[expectedKey]

		if !ok {
			t.Errorf("No pair found in res.Pairs for expected key")
		}

		testIntegerObject(t, pair.Value, expectedVal)

	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
	}

	for _, tt := range tests {

		evaled := testEval(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, evaled, int64(integer))
		} else {
			testNullObj(t, evaled)
		}

	}
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
