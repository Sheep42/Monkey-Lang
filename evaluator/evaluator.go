package evaluator

import (
	"fmt"

	"github.com/Sheep42/Monkey-Lang/ast"
	"github.com/Sheep42/Monkey-Lang/object"
)

// Literal Null / True / False
var (
	Null  = &object.Null{}
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
)

// Eval evaluates the AST
func Eval(node ast.Node, env *object.Environment) object.Object {

	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObj(node.Value)

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)

		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}

		return &object.Array{Elements: elements}

	case *ast.IndexExpression:
		left := Eval(node.Left, env)

		if isError(left) {
			return left
		}

		index := Eval(node.Index, env)

		if isError(index) {
			return index
		}

		return evalIndexExpression(left, index)

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	case *ast.PrefixExpression:

		right := Eval(node.Right, env)

		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:

		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.IfExpression:
		return evalIfElseExpression(node, env)

	case *ast.ReturnStatement:

		val := Eval(node.ReturnValue, env)

		if isError(val) {
			return val
		}

		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:

		val := Eval(node.Value, env)

		if isError(val) {
			return val
		}

		env.Set(node.Name.Value, val)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:

		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}

	case *ast.CallExpression:

		fn := Eval(node.Function, env)

		if isError(fn) {
			return fn
		}

		args := evalExpressions(node.Arguments, env)

		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFn(fn, args)

	}

	return nil

}

func isError(obj object.Object) bool {

	if obj != nil {

		return obj.Type() == object.ErrorObj

	}

	return false

}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {

	var res object.Object

	for _, stmt := range program.Statements {

		res = Eval(stmt, env)

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

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {

	var res object.Object

	for _, stmt := range block.Statements {

		res = Eval(stmt, env)

		if res != nil {

			t := res.Type()

			// bail out early if we hit a return or error
			if t == object.ReturnValueObj || t == object.ErrorObj {
				return res
			}

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

	case left.Type() == object.StringObj && right.Type() == object.StringObj:
		return evalInfixStringExpression(operator, left, right)

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

func evalInfixStringExpression(operator string, left, right object.Object) object.Object {

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {

	case "+":
		return &object.String{Value: leftVal + rightVal}

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

func evalIfElseExpression(ie *ast.IfExpression, env *object.Environment) object.Object {

	condition := Eval(ie.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return Null
	}

}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {

	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: %s", node.Value)

}

func evalIndexExpression(left, index object.Object) object.Object {

	switch {

	case left.Type() == object.ArrayObj && index.Type() == object.IntegerObj:
		return evalArrayIndexExpression(left, index)

	case left.Type() == object.HashObj:
		return evalHashIndexExpression(left, index)

	default:
		return newError("Index operator not supported: %s[%s]", left.Type(), index.Type())

	}

}

func evalArrayIndexExpression(array, index object.Object) object.Object {

	arr := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arr.Elements) - 1)

	if idx < 0 || idx > max {
		return Null
	}

	return arr.Elements[idx]

}

func evalHashIndexExpression(hash, index object.Object) object.Object {

	hashObject := hash.(*object.Hash)
	key, ok := index.(object.Hashable)

	if !ok {
		return newError("Invalid HashKey: %q. Type %q is unsupported.", index.Inspect(), index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]

	if !ok {
		return Null
	}

	return pair.Value

}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {

	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valNode := range node.Pairs {

		key := Eval(keyNode, env)

		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)

		if !ok {
			return newError("Invalid HashKey: %q. Type %q is unsupported.", key.Inspect(), key.Type())
		}

		val := Eval(valNode, env)

		if isError(val) {
			return val
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: val}

	}

	return &object.Hash{Pairs: pairs}

}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {

	var res []object.Object

	for _, e := range exps {

		evaluated := Eval(e, env)

		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		res = append(res, evaluated)

	}

	return res

}

func applyFn(fn object.Object, args []object.Object) object.Object {

	switch fn := fn.(type) {

	case *object.Function:

		extendedEnv := extendFnEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnVal(evaluated)

	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type())

	}

}

func extendFnEnv(fn *object.Function, args []object.Object) *object.Environment {

	env := object.NewEnclosedEnvironment(fn.Env)

	for i, param := range fn.Parameters {

		env.Set(param.Value, args[i])

	}

	return env

}

func unwrapReturnVal(obj object.Object) object.Object {

	if returnVal, ok := obj.(*object.ReturnValue); ok {
		return returnVal.Value
	}

	return obj

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
