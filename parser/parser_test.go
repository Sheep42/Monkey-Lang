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
		{"let x = 5", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
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

func TestStringLiteralExpression(t *testing.T) {

	tests := []struct {
		input    string
		expected string
	}{
		{`"hello, world"`, "hello, world"},
		{`'hello, world'`, "hello, world"},
	}
	for _, tt := range tests {

		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		literal, ok := stmt.Expression.(*ast.StringLiteral)

		if !ok {
			t.Fatalf("Statement was incorrect type. Expected=\"*ast.StringLiteral\". Got=\"%T\"", stmt)
		}

		if literal.Value != tt.expected {
			t.Errorf("Literal value was incorrect. Expected=%q. Got=%q", tt.expected, literal.Value)
		}

	}
}

func TestArrayLiteralExpression(t *testing.T) {

	input := "[1, 2 * 2, 3 + 3]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)

	if !ok {
		t.Errorf("Expression was incorrect type. Expected=ast.ArrayLiteral. Got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Errorf("Number of array elements was incorrect. Expected=%d. Got=%d", 3, len(array.Elements))
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)

}

func TestIndexExpression(t *testing.T) {

	input := "myArray[1 + 1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)

	if !ok {
		t.Fatalf("Expression was incorrect type. Expected=%s. Got=%T", "*ast.IndexExpression", stmt.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}

}

func TestParsingHashLiteral(t *testing.T) {

	input := `{"one" : 1, "two" : 2, "three" : 3}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)

	if !ok {
		t.Fatalf("Expression was incorrect type. Exptected=%s. Got=%T", "*ast.HashLiteral", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Fatalf("hash.Pairs was incorrect length. Expected=%d. Got=%d", 3, len(hash.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for k, v := range hash.Pairs {

		literal, ok := k.(*ast.StringLiteral)

		if !ok {
			t.Errorf("Key was not correct type. Expected=%s. Got=%T", "*ast.StringLiteral", k)
		}

		expectedVal := expected[literal.String()]

		testIntegerLiteral(t, v, expectedVal)

	}

}

func TestParsingEmptyHashLiteral(t *testing.T) {

	input := `{}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)

	if !ok {
		t.Fatalf("Expression was incorrect type. Exptected=%s. Got=%T", "*ast.HashLiteral", stmt.Expression)
	}

	if len(hash.Pairs) != 0 {
		t.Fatalf("hash.Pairs was incorrect length. Expected=%d. Got=%d", 0, len(hash.Pairs))
	}

}

func TestParsingHashLiteralExpressions(t *testing.T) {

	input := `{"one" : 0 + 1, "two" : 10 - 8, "three" : 15 / 5}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)

	if !ok {
		t.Fatalf("Expression was incorrect type. Exptected=%s. Got=%T", "*ast.HashLiteral", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Fatalf("hash.Pairs was incorrect length. Expected=%d. Got=%d", 3, len(hash.Pairs))
	}

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

	for k, v := range hash.Pairs {

		literal, ok := k.(*ast.StringLiteral)

		if !ok {
			t.Errorf("Key was not correct type. Expected=%s. Got=%T", "*ast.StringLiteral", k)
		}

		testFunc, ok := tests[literal.String()]

		if !ok {
			t.Errorf("No test function for key %q was found", literal.String())
			continue
		}

		testFunc(v)

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
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
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

func TestIfExpression(t *testing.T) {

	input := "if( x < y ) { x }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {

		t.Fatalf("Program contains incorrect number of statements. Expected=%d. Got=%d", 1, len(program.Statements))

	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {

		t.Fatalf("Statement is not an ast.ExpressionStatement. Got=%T", program.Statements[0])

	}

	exp, ok := stmt.Expression.(*ast.IfExpression)

	if !ok {

		t.Fatalf("Expression is not an ast.IfExpression. Got=%T", stmt.Expression)

	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {

		t.Errorf(
			"Consequence contains incorrect number of statements. Expected=%d. Got=%d",
			1,
			len(exp.Consequence.Statements),
		)

	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {

		t.Fatalf(
			"Consequence.Statements[0] is not correct type. Expected=%s. Got=%T",
			"*ast.ExpressionStatement",
			exp.Consequence.Statements[0],
		)

	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {

		t.Errorf("exp.Alternative was not nil. Got=%+v", exp.Alternative)

	}

}

func TestIfElseExpression(t *testing.T) {

	input := "if( x < y ) { x } else { y }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {

		t.Fatalf("Program contains incorrect number of statements. Expected=%d. Got=%d", 1, len(program.Statements))

	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {

		t.Fatalf("Statement is not an ast.ExpressionStatement. Got=%T", program.Statements[0])

	}

	exp, ok := stmt.Expression.(*ast.IfExpression)

	if !ok {

		t.Fatalf("Expression is not an ast.IfExpression. Got=%T", stmt.Expression)

	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {

		t.Errorf(
			"Consequence contains incorrect number of statements. Expected=%d. Got=%d",
			1,
			len(exp.Consequence.Statements),
		)

	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {

		t.Fatalf(
			"Consequence.Statements[0] is not correct type. Expected=%s. Got=%T",
			"*ast.ExpressionStatement",
			exp.Consequence.Statements[0],
		)

	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {

		t.Errorf(
			"Consequence contains incorrect number of statements. Expected=%d. Got=%d",
			1,
			len(exp.Consequence.Statements),
		)

	}

	alt, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)

	if !ok {

		t.Fatalf(
			"Alternative.Statements[0] is not correct type. Expected=%s. Got=%T",
			"*ast.ExpressionStatement",
			exp.Alternative.Statements[0],
		)

	}

	if !testIdentifier(t, alt.Expression, "y") {
		return
	}

}

func TestFunctionLiteralParsing(t *testing.T) {

	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {

		t.Fatalf("Program has incorrect number of statements. Expected=%d. Got=%d", 1, len(program.Statements))

	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {

		t.Fatalf("Statement is of incorrect type. Expected=%s. Got=%T", "*ast.ExpressionStatement", program.Statements[0])

	}

	fn, ok := stmt.Expression.(*ast.FunctionLiteral)

	if !ok {

		t.Fatalf("Expression is of incorrect type. Expected=%s. Got=%T", "*ast.FunctionLiteral", stmt.Expression)

	}

	if len(fn.Parameters) != 2 {

		t.Fatalf("Function has incorrect number of params. Expected=%d. Got=%d", 2, len(fn.Parameters))

	}

	testLiteralExpression(t, fn.Parameters[0], "x")
	testLiteralExpression(t, fn.Parameters[1], "y")

	if len(fn.Body.Statements) != 1 {

		t.Fatalf("Function body has incorrect number of statements. Expected=%d. Got=%d", 1, len(fn.Body.Statements))

	}

	body, ok := fn.Body.Statements[0].(*ast.ExpressionStatement)

	if !ok {

		t.Fatalf("Body statement is of incorrect type. Expected=%s. Got=%T", "*ast.ExpressionStatement", fn.Body.Statements[0])

	}

	testInfixExpression(t, body.Expression, "x", "+", "y")

}

func TestFunctionParamParsing(t *testing.T) {

	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn(){}", expectedParams: []string{}},
		{input: "fn(x){}", expectedParams: []string{"x"}},
		{input: "fn(x, y) {}", expectedParams: []string{"x", "y"}},
		{input: "fn(x, y, z) {}", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {

		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		fn := stmt.Expression.(*ast.FunctionLiteral)

		if len(fn.Parameters) != len(tt.expectedParams) {

			t.Errorf("length of params is incorrect. Expected=%d. Got=%d", len(tt.expectedParams), len(fn.Parameters))

		}

		for i, ident := range tt.expectedParams {

			testLiteralExpression(t, fn.Parameters[i], ident)

		}
	}
}

func TestCallExpressionParsing(t *testing.T) {

	input := `add(1, 2 * 3, 4 + 5);`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {

		t.Fatalf("program.Statements contains the wrong number of statements. Expected=%d. Got=%d", 1, len(program.Statements))

	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {

		t.Fatalf("stmt is of incorrect type. Expected=%s. Got=%T", "*ast.ExpressionStatement", program.Statements[0])

	}

	exp, ok := stmt.Expression.(*ast.CallExpression)

	if !ok {

		t.Fatalf("stmt.Expression is of incorrect type. Expected=%s. Got=%T", "*ast.CallExpression", stmt.Expression)

	}

	if !testIdentifier(t, exp.Function, "add") {

		return

	}

	if len(exp.Arguments) != 3 {

		t.Fatalf("Wrong number of arguments. Expected=%d. Got=%d", 3, len(exp.Arguments))

	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)

}

func TestCallExpressionArgs(t *testing.T) {

	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "add()", expectedParams: []string{}},
		{input: "add(x)", expectedParams: []string{"x"}},
		{input: "add(x, y)", expectedParams: []string{"x", "y"}},
		{input: "add(x, y, z)", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {

		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		callFn := stmt.Expression.(*ast.CallExpression)

		if len(callFn.Arguments) != len(tt.expectedParams) {

			t.Errorf("number of args is incorrect. Expected=%d. Got=%d", len(tt.expectedParams), len(callFn.Arguments))

		}

		for i, ident := range tt.expectedParams {

			testLiteralExpression(t, callFn.Arguments[i], ident)

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
