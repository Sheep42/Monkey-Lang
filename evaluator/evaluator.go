package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

var (
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
