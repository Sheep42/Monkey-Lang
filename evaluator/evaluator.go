package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

var (
	Null  = &object.Null{}
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {

	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObj(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)

	}

	return nil

}

func evalStatements(stmts []ast.Statement) object.Object {

	var res object.Object

	for _, stmt := range stmts {

		res = Eval(stmt)

	}

	return res

}

func nativeBoolToBooleanObj(input bool) object.Object {

	if input {

		return True

	}

	return False

}

func evalPrefixExpression(operator string, right object.Object) object.Object {

	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalNegationOperatorExpression(right)
	default:
		return Null
	}

}

func evalBangOperatorExpression(right object.Object) object.Object {

	switch right {
	case True:
		return False
	case False:
		return True
	case Null:
		return True
	default:
		return False
	}

}

func evalNegationOperatorExpression(right object.Object) object.Object {

	if right.Type() != object.IntegerObj {

		return Null

	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}

}
