package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 838383;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]

		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
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

func TestReturnStatements(t *testing.T) {

	input := `
		return 5;
		return 10;
		return 9281875;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {

		returnStmt, ok := stmt.(*ast.ReturnStatement)

		if !ok {
			t.Fatalf("stmt not *ast.ReturnStatement. got=%T", stmt)
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral() not 'return'. got=%q", returnStmt.TokenLiteral())
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

func TestParsingPrefixExpressions(t *testing.T) {

	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
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

		if !testIntegerLiteral(t, expr.Right, tt.integerValue) {

			return

		}

	}

}

func TestParsingInfixExpressions(t *testing.T) {

	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
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

		if !testIntegerLiteral(t, expr.Left, tt.leftValue) {

			return

		}

		if expr.Operator != tt.operator {

			t.Fatalf("expression operator is not %s. got=%s", tt.operator, expr.Operator)

		}

		if !testIntegerLiteral(t, expr.Right, tt.rightValue) {

			return

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
