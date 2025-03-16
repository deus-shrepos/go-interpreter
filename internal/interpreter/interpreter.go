package interpreter

import (
	"Crafting-interpreters/internal/ast"
	"Crafting-interpreters/internal/errors"
	"Crafting-interpreters/internal/token"
	_ "errors"
	"fmt"
	"strings"
)

/*
The job of the interpreter (runtime evaluation) is to take a tree to its source and evaluate it literals
*/

type Interpreter struct {
	Error errors.Error
}

func (i *Interpreter) Interpret(expr ast.Expr) error {
	value, err := i.eval(expr)
	if err != nil {
		return err
	}
	fmt.Printf(stringify(value))
	return nil
}

func (i *Interpreter) VisitLiteral(expr ast.Literal) (interface{}, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitUnary(expr ast.Unary) (interface{}, error) {
	right, _ := i.eval(expr.Right) // POST ORDER EVALUATION
	switch expr.Operator.Type {
	case token.MINUS:
		return -right.(float64), nil
	case token.BANG:
		return !IsTruthy(right), nil
	default:
		return nil, fmt.Errorf("%s is not a valid operator", expr.Operator.Lexeme)
	}
}

func (i *Interpreter) VisitGrouping(expr ast.Grouping) (interface{}, error) {
	return i.eval(expr.Expression)
}

func (i *Interpreter) eval(expr ast.Expr) (interface{}, error) {
	return expr.Accept(i)
}

func (i *Interpreter) VisitBinary(expr ast.Binary) (interface{}, error) {
	right, _ := i.eval(expr.Right)
	left, _ := i.eval(expr.Left)

	switch expr.Operator.Type {
	case token.MINUS:
		return right.(float64) - left.(float64), nil
	case token.PLUS:
		switch rightValue := right.(type) {
		case float64:
			leftValue, ok := left.(float64)
			if !ok {
				return nil, fmt.Errorf("'%v' operand must be number", leftValue)
			}
			return rightValue + left.(float64), nil
		case string:
			return rightValue + left.(string), nil
		default:
			return nil, fmt.Errorf("'%v' operand must be a number or a string", rightValue)
		}
	case token.SLASH:
		switch rightValue := right.(type) {
		case float64:
			return rightValue / left.(float64), nil
		default:
			return nil, fmt.Errorf("operand must be a number", expr.Operator.Lexeme)
		}
	case token.STAR:
		switch rightValue := right.(type) {
		case float64:
			return rightValue * left.(float64), nil
		default:
			return nil, fmt.Errorf("%s is not a valid operator", expr.Operator.Lexeme)
		}
	case token.GREATER:
		switch rightValue := right.(type) {
		case float64:
			return rightValue > left.(float64), nil
		case string:
			return rightValue > left.(string), nil
		default:
			return nil, fmt.Errorf("%s is not a valid operator", expr.Operator.Lexeme)
		}
	case token.GREATER_EQUAL:
		switch rightValue := right.(type) {
		case float64:
			return rightValue >= left.(float64), nil
		case string:
			return rightValue >= left.(string), nil
		}
	case token.LESS:
		switch rightValue := right.(type) {
		case float64:
			return rightValue < left.(float64), nil
		case string:
			return rightValue < left.(string), nil
		}
	case token.LESS_EQUAL:
		switch rightValue := right.(type) {
		case float64:
			return rightValue <= left.(float64), nil
		case string:
			return rightValue <= left.(string), nil
		}
	case token.BANG_EQUAL:
		return !isEqual(left, right), nil
	case token.EQUAL_EQUAL:
		return isEqual(left, right), nil
	case token.AND:
		switch rightValue := right.(type) {
		case bool:
			return rightValue && left.(bool), nil
		default:
			return nil, fmt.Errorf("%s is not a valid operator", expr.Operator.Lexeme)
		}
	default:
		return nil, fmt.Errorf("%s is not a valid operator", expr.Operator.Lexeme)
	}
	return nil, fmt.Errorf("%s is not a valid operator", expr.Operator.Lexeme)
}

func isEqual(left, right interface{}) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil {
		return false
	}
	return left == right
}
func IsTruthy(object interface{}) bool {
	if object == nil {
		return false
	}
	isBoolean, ok := object.(bool)
	return ok && isBoolean
}

func stringify(object any) string {
	if object == nil {
		return "nil"
	}

	value, isDouble := object.(float64)
	if isDouble {
		valueString := fmt.Sprint(value)
		if strings.HasSuffix(valueString, ".0") {
			valueString = valueString[0 : len(valueString)-2]
		}
		return valueString
	}
	return fmt.Sprintf("%v", object)
}
