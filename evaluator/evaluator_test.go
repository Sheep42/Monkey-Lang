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

func testEval(input string) object.Object {

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)

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
