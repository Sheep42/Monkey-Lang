package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {

	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		// {"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {

		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain 1 statement. got=%d",
				len(program.Statements),
			)
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}

	}

}

/*

func TestReturnStatements(t *testing.T) {

	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {

		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {

			t.Fatalf(
				"program.Statements does not contain 1 statement. got=%d",
				len(program.Statements),
			)

		}

		stmt := program.Statements[0]
		returnStmt, ok := stmt.(*ast.ReturnStatement)

		if !ok {
			t.Fatalf("stmt not *ast.returnStatement. got=%T", stmt)
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}

		if testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue) {
			return
		}

	}

}

*/

func TestIdentifierExpression(t *testing.T) {

	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {

		t.Fatalf("program has wrong number of Statements. got=%d", len(program.Statements))

	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {

		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])

	}

	ident, ok := stmt.Expression.(*ast.Identifier)

	if !ok {

		t.Fatalf("Expression not identifier. got=%T", stmt.Expression)

	}

	if ident.Value != "foobar" {

		t.Fatalf("ident.Value is not %s. got=%s", "foobar", ident.Value)

	}

	if ident.TokenLiteral() != "foobar" {

		t.Fatalf("ident.Value is not %s. got=%s", "foobar", ident.TokenLiteral())

	}

}

func TestIntegerLiteralExpression(t *testing.T) {

	input := "5;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {

		t.Fatalf("Wrong number of statements in program.Statements. expected=%d, got=%d", 1, len(program.Statements))

	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {

		t.Fatalf("program.Statements[0] is not an ast.ExpressionStatement. got=%T", program.Statements[0])

	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)

	if !ok {

		t.Fatalf("Expression is not *ast.IntegerLiteral. got=%T", stmt.Expression)

	}

	if literal.Value != 5 {

		t.Fatalf("literal.Value != %d. got=%d", 5, literal.Value)

	}

	if literal.TokenLiteral() != "5" {

		t.Fatalf("literal.TokenLiteral() != \"%s\". got=%s", "5", literal.TokenLiteral())

	}
}

func TestParsingPrefixExpressions(t *testing.T) {

	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		// {"!foobar;", "!", "foobar"},
		// {"-foobar;", "-", "foobar"},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {

		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {

			t.Fatalf("program.Statements contains wrong number of statements. expected=%d. got=%d", 1, len(program.Statements))

		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {

			t.Fatalf("statement is not *ast.ExpressionStatement. got=%T", program.Statements[0])

		}

		expr, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {

			t.Fatalf("expression is not *ast.PrefixExpression. got=%d", stmt.Expression)

		}

		if expr.Operator != tt.operator {

			t.Fatalf("expression operator is not %s. got=%s", tt.operator, expr.Operator)

		}

		if !testLiteralExpression(t, expr.Right, tt.value) {

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
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		// {"foobar + barfoo;", "foobar", "+", "barfoo"},
		// {"foobar - barfoo;", "foobar", "-", "barfoo"},
		// {"foobar * barfoo;", "foobar", "*", "barfoo"},
		// {"foobar / barfoo;", "foobar", "/", "barfoo"},
		// {"foobar > barfoo;", "foobar", ">", "barfoo"},
		// {"foobar < barfoo;", "foobar", "<", "barfoo"},
		// {"foobar == barfoo;", "foobar", "==", "barfoo"},
		// {"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {

		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {

			t.Fatalf("program.Statements contains wrong number of statements. expected=%d. got=%d", 1, len(program.Statements))

		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {

			t.Fatalf("statement is not *ast.ExpressionStatement. got=%T", program.Statements[0])

		}

		expr, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {

			t.Fatalf("expression is not *ast.InfixExpression. got=%d", stmt.Expression)

		}

		if !testInfixExpression(t, expr, tt.leftValue, tt.operator, tt.rightValue) {

			return

		}

	}

}

func TestParsingOperatorPrecedence(t *testing.T) {

	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
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
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		// {
		// 	"1 + (2 + 3) + 4",
		// 	"((1 + (2 + 3)) + 4)",
		// },
		// {
		// 	"(5 + 5) * 2",
		// 	"((5 + 5) * 2)",
		// },
		// {
		// 	"2 / (5 + 5)",
		// 	"(2 / (5 + 5))",
		// },
		// {
		// 	"(5 + 5) * 2 * (5 + 5)",
		// 	"(((5 + 5) * 2) * (5 + 5))",
		// },
		// {
		// 	"-(5 + 5)",
		// 	"(-(5 + 5))",
		// },
		// {
		// 	"!(true == true)",
		// 	"(!(true == true))",
		// },
		// {
		// 	"a + add(b * c) + d",
		// 	"((a + add((b * c))) + d)",
		// },
		// {
		// 	"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
		// 	"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		// },
		// {
		// 	"add(a + b + c * d / f + g)",
		// 	"add((((a + b) + ((c * d) / f)) + g))",
		// },
	}

	for _, tt := range tests {

		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()

		if actual != tt.expected {

			t.Errorf("expected=%q. got=%q", tt.expected, actual)

		}

	}

}

func TestBooleanExpression(t *testing.T) {

	tests := []struct {
		input         string
		expectedValue bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {

		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {

			t.Fatalf(
				"program.Statements does not contain 1 statement. got=%d",
				len(program.Statements),
			)

		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {

			t.Fatalf(
				"Statement is incorrect type. exptected=%s. got=%T",
				"ast.ExpressionStatement",
				program.Statements[0],
			)

		}

		boolean, ok := stmt.Expression.(*ast.Boolean)
		if !ok {

			t.Fatalf(
				"Exression is incorrect type. expected=%s. got=%T",
				"ast.Boolean",
				stmt.Expression,
			)

		}

		if boolean.Value != tt.expectedValue {

			t.Errorf(
				"boolean.Value is incorrect. expected=%t. got=%t",
				tt.expectedValue,
				boolean.Value,
			)

		}

	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("Parser Error: %q", msg)
	}

	t.FailNow()
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)

	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name)
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {

	integer, ok := il.(*ast.IntegerLiteral)

	if !ok {

		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false

	}

	if integer.Value != value {

		t.Errorf("integer value not correct. expected=%d. got=%d", value, integer.Value)
		return false

	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {

		t.Errorf("integer TokenLiteral was not correct. expected=%d. got=%s", value, integer.TokenLiteral())
		return false

	}

	return true

}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {

	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("Expression not ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("Identifier value is incorrect. expected=%s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("Identifier TokenLiteral is incorrect. expected=%s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true

}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {

	switch v := expected.(type) {

	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)

	}

	t.Errorf("Type of Expression not handled. got=%T", exp)

	return false

}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {

	bo, ok := exp.(*ast.Boolean)

	if !ok {

		t.Errorf(
			"Expression was not a Boolean Literal. got=%T",
			exp,
		)
		return false

	}

	if bo.Value != value {

		t.Errorf(
			"Boolean value was incorrect. expected=%t. got=%t",
			value,
			bo.Value,
		)
		return false

	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {

		t.Errorf(
			"TokenLiteral value was incorrect. expected=%t. got=%s",
			value,
			bo.TokenLiteral(),
		)
		return false

	}

	return true

}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {

	opExp, ok := exp.(*ast.InfixExpression)

	if !ok {

		t.Errorf("Expression was not an ast.InfixExpression. got=%T", exp)
		return false

	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {

		t.Errorf("Expression operator is incorrect. expected=%s. got=%s", operator, opExp.Operator)
		return false

	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true

}
