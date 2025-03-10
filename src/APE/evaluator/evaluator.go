package evaluator

import (
	"fmt"
	"APE/ast"
	"APE/object"
)


// boolean values  

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// switches between availale statments in order to intepret 
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// statements
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.StringLiteral: 
		return &object.String{Value: node.Value}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return val 
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env) 
	case *ast.FunctionLiteral: 
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}
	// Expressions
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
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
	case *ast.CallExpression: 
		function := Eval(node.Function, env)
		if isError(function) { 
			return function 
		}
		args := evalExpressions(node.Arguments, env) 
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
// For the MethodCallExpression case, use this ultra-minimal implementation:

case *ast.MethodCallExpression:
	o := Eval(node.Object, env)
	if isError(o) {
		return o
	}
	a := evalExpressions(node.Arguments, env)
	if len(a) > 0 && isError(a[0]) {
		return a[0]
	}

	// Here we ONLY check type using a type switch
	switch arr := o.(type) {
	case *object.Array:
		// It's an array
		if node.Method == "map" {
			// Handle map
			if len(a) != 1 {
				return newError("wrong number of arguments for map")
			}
			
			fn, ok := a[0].(*object.Function)
			if !ok {
				return newError("argument to map must be a function")
			}
			
			result := []object.Object{}
			for _, e := range arr.Elements {
				val := applyFunction(fn, []object.Object{e})
				if isError(val) {
					return val
				}
				result = append(result, val)
			}
			
			return &object.Array{Elements: result}
		} else if node.Method == "filter" {
			// Handle filter
			if len(a) != 1 {
				return newError("wrong number of arguments for filter")
			}
			
			fn, ok := a[0].(*object.Function)
			if !ok {
				return newError("argument to filter must be a function")
			}
			
			result := []object.Object{}
			for _, e := range arr.Elements {
				condition := applyFunction(fn, []object.Object{e})
				if isError(condition) {
					return condition
				}
				
				if isTruthy(condition) {
					result = append(result, e)
				}
			}
			
			return &object.Array{Elements: result}
		} else if node.Method == "reduce" {
			// Handle reduce
			if len(a) != 2 {
				return newError("wrong number of arguments for reduce")
			}
			
			fn, ok := a[0].(*object.Function)
			if !ok {
				return newError("first argument to reduce must be a function")
			}
			
			accum := a[1]
			for _, e := range arr.Elements {
				accum = applyFunction(fn, []object.Object{accum, e})
				if isError(accum) {
					return accum
				}
			}
			
			return accum
		} else {
			return newError("unknown method %s", node.Method)
		}
	default:
		// Not an array
		return newError("no methods for this type")
	}
	case *ast.ArrayLiteral: 
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
		}
	}
	return nil
}


func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}
	return &object.Hash{Pairs: pairs}
}

func evalIndexExpression(left, index object.Object) object.Object {
    
    switch {
    case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
        arrayObject := left.(*object.Array)
        idx := index.(*object.Integer).Value
        max := int64(len(arrayObject.Elements) - 1)
        
        
        if idx < 0 || idx > max {
            return NULL
        }
        element := arrayObject.Elements[idx]
        return element
		case left.Type() == object.HASH_OBJ: 
			return evalHashIndexExpression(left, index)
    default:
        return newError("index operator not supported: %s", left.Type())
    }
}


func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}


func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
		case *object.Function:
			extendedEnv := extendFunctionEnv(fn, args)
			evaluated := Eval(fn.Body, extendedEnv)
			return unwrapReturnValue(evaluated)			
		case *object.Builtin: 
			return fn.Fn(args...)
		default:
			return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}


func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj 
}

func evalExpression(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env) 
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}


func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)

		if isError(evaluated){
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {

	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())

	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("operator mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

}

func evalStringInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return &object.String{Value: leftVal + rightVal}

}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unkown operator: %s%s", operator, right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}

}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalProgram(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement, env)

		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}

	return result
}

func evalBlockStatement(statements *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range statements.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
