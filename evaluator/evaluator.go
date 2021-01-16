package evaluator

import (
	"fmt"
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
		return evalProgram(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObj(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)

		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}

		right := Eval(node.Right)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.IfExpression:
		return evalIfElseExpression(node)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)

		if isError(val) {
			return val
		}

		return &object.ReturnValue{Value: val}

	}

	return nil

}

func isError(obj object.Object) bool {

	if obj != nil {

		return obj.Type() == object.ErrorObj

	}

	return false

}

func evalProgram(stmts []ast.Statement) object.Object {

	var res object.Object

	for _, stmt := range stmts {

		res = Eval(stmt)

		// bail out early if we hit a return or error
		switch res := res.(type) {
		case *object.ReturnValue:
			return res.Value
		case *object.Error:
			return res
		}

	}

	return res

}

func evalBlockStatement(block *ast.BlockStatement) object.Object {

	var res object.Object

	for _, stmt := range block.Statements {

		res = Eval(stmt)
		t := res.Type()

		// bail out early if we hit a return or error
		if res != nil && t == object.ReturnValueObj || t == object.ErrorObj {

			return res

		}

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
		return newError("unknown operator: %s%s", operator, right.Type())
	}

}

func evalInfixExpression(operator string, left, right object.Object) object.Object {

	switch {
	case left.Type() == object.IntegerObj && right.Type() == object.IntegerObj:
		return evalInfixIntegerExpression(operator, left, right)

	case operator == "==":
		return nativeBoolToBooleanObj(left == right)

	case operator == "!=":
		return nativeBoolToBooleanObj(left != right)

	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

}

func evalInfixIntegerExpression(operator string, left, right object.Object) object.Object {

	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {

	case "*":
		return &object.Integer{Value: leftVal * rightVal}

	case "/":
		return &object.Integer{Value: leftVal / rightVal}

	case "+":
		return &object.Integer{Value: leftVal + rightVal}

	case "-":
		return &object.Integer{Value: leftVal - rightVal}

	case "<":
		return nativeBoolToBooleanObj(leftVal < rightVal)

	case ">":
		return nativeBoolToBooleanObj(leftVal > rightVal)

	case "==":
		return nativeBoolToBooleanObj(leftVal == rightVal)

	case "!=":
		return nativeBoolToBooleanObj(leftVal != rightVal)

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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

		return newError("unknown operator: -%s", right.Type())

	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}

}

func evalIfElseExpression(ie *ast.IfExpression) object.Object {

	condition := Eval(ie.Condition)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return Null
	}

}

func isTruthy(obj object.Object) bool {

	switch obj {

	case Null:
		return false
	case False:
		return false
	case True:
		return true
	default:
		return true

	}

}

func newError(format string, a ...interface{}) *object.Error {

	return &object.Error{Message: fmt.Sprintf(format, a...)}

}
