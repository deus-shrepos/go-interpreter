package interpreter

import (
	_ "errors"
	"fmt"
	"strings"

	"github.com/go-interpreter/internal/ast"
	"github.com/go-interpreter/internal/errors"
	"github.com/go-interpreter/internal/token"
)

/*
The job of the interpreter (runtime evaluation) is to take a tree to its source and evaluate it literals
*/

type Interpreter struct{}

func (i *Interpreter) Interpret(expr ast.Expr, printExpression bool) {
	if expr == nil {
		return
	}
	value, err := i.eval(expr)
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
	}
	if printExpression {
		fmt.Print(stringify(value))
	}
}

func (i *Interpreter) VisitLiteral(expr ast.Literal) (any, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitUnary(expr ast.Unary) (any, error) {
	right, _ := i.eval(expr.Right) // POST ORDER EVALUATION
	switch expr.Operator.Type {
	case token.MINUS:
		return -right.(float64), nil
	case token.BANG:
		return !IsTruthy(right), nil
	default:
		return nil, errors.ExecutionError{Type: errors.RUNTIME_ERROR,
			Line:    expr.Operator.Line,
			Where:   expr.Operator.Char,
			Message: fmt.Sprintf("%s is not a valid operator", expr.Operator.Lexeme)}
	}
}

func (i *Interpreter) VisitGrouping(expr ast.Grouping) (any, error) {
	return i.eval(expr.Expression)
}

func (i *Interpreter) eval(expr ast.Expr) (any, error) {
	return expr.Accept(i)
}

func (i *Interpreter) VisitBinary(expr ast.Binary) (any, error) {
	left, _ := i.eval(expr.Left)
	right, _ := i.eval(expr.Right)
	switch expr.Operator.Type {
	case token.MINUS:
		return right.(float64) - left.(float64), nil
	case token.PLUS:
		// Check if the operands are strings
		if leftValue, ok := left.(string); ok {
			if rightValue, ok := right.(string); ok {
				return rightValue + leftValue, nil
			}
		}
		if rightValue, ok := right.(string); ok {
			// We know that the right is a string, so we need to check if the left is a number
			err := checkIfNumber(left, expr.Operator)
			if err != nil {
				return nil, err
			}
			return rightValue + fmt.Sprintf("%v", left.(float64)), nil
		}
		err := checkIfNumbers(left, right, expr.Operator)
		if err != nil {
			return nil, err
		}
		return right.(float64) + left.(float64), nil

	case token.SLASH:
		err := checkIfNumbers(left, right, expr.Operator)
		if err != nil {
			return nil, err
		}
		return right.(float64) / left.(float64), nil
	case token.STAR:
		err := checkIfNumbers(left, right, expr.Operator)
		if err != nil {
			return nil, err
		}
		return right.(float64) * left.(float64), nil
	case token.GREATER:
		err := checkIfNumbers(left, right, expr.Operator)
		if err != nil {
			return nil, err
		}
		return right.(float64) > left.(float64), nil
	case token.GREATER_EQUAL:
		err := checkIfNumbers(left, right, expr.Operator)
		if err != nil {
			return nil, err
		}
		return right.(float64) >= left.(float64), nil
	case token.LESS:
		err := checkIfNumbers(left, right, expr.Operator)
		if err != nil {
			return nil, err
		}
		return right.(float64) < left.(float64), nil
	case token.LESS_EQUAL:
		err := checkIfNumbers(left, right, expr.Operator)
		if err != nil {
			return nil, err
		}
		return right.(float64) <= left.(float64), nil
	case token.BANG_EQUAL:
		return !isEqual(left, right), nil
	case token.EQUAL_EQUAL:
		return isEqual(left, right), nil
	case token.AND:
		err := checkIfBooleans(left, right, expr.Operator)
		if err != nil {
			return nil, err
		}
		return left.(bool) && right.(bool), nil
	case token.OR:
		err := checkIfBooleans(left, right, expr.Operator)
		if err != nil {
			return nil, err
		}
		return left.(bool) || right.(bool), nil
	default:
		return nil, errors.ExecutionError{Type: errors.RUNTIME_ERROR,
			Line:    expr.Operator.Line,
			Where:   expr.Operator.Char,
			Message: fmt.Sprintf("%s is not a valid operator", expr.Operator.Lexeme)}
	}
}

func isEqual(left, right any) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil {
		return false
	}
	return left == right
}
func IsTruthy(object any) bool {
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
		valueString = strings.TrimSuffix(valueString, ".0")
		return valueString
	}
	return fmt.Sprintf("%v", object)
}

func checkIfNumber(object any, operator token.Token) error {
	if _, ok := object.(float64); !ok {
		return errors.ExecutionError{Type: errors.RUNTIME_ERROR,
			Line:    operator.Line,
			Where:   operator.Char,
			Message: fmt.Sprintf("'%v' Operand must be a number", object)}
	}
	return nil
}

func checkIfBoolean(object any, operator token.Token) error {
	if _, ok := object.(bool); !ok {
		return errors.ExecutionError{Type: errors.RUNTIME_ERROR,
			Line:    operator.Line,
			Where:   operator.Char,
			Message: fmt.Sprintf("'%v' Operand must be a boolean", object)}
	}
	return nil
}

func checkIfBooleans(left, right any, operator token.Token) error {
	err := checkIfBoolean(left, operator)
	if err != nil {
		return err
	}
	err = checkIfBoolean(right, operator)
	if err != nil {
		return err
	}
	return nil
}

func checkIfNumbers(left, right any, operator token.Token) error {
	err := checkIfNumber(left, operator)
	if err != nil {
		return err
	}
	err = checkIfNumber(right, operator)
	if err != nil {
		return err
	}
	return nil
}
