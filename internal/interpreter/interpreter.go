package interpreter

import (
	_ "errors"
	"fmt"
	"strings"

	"github.com/go-interpreter/internal/ast"
	"github.com/go-interpreter/internal/errors"
	"github.com/go-interpreter/internal/token"
)

// Interpreter represents the core structure for the interpreter.
// It is responsible for executing and evaluating code based on the
// implemented logic and rules of the interpreter.
type Interpreter struct {
	environment Environment
}

func NewInterpreter() Interpreter {
	return Interpreter{}
}

// Interpret executes a series of statements provided as input.
// It iterates over each statement, executing them one by one using the exec method.
// If an error occurs during the execution of a statement, it logs the error to the console.
func (i *Interpreter) Interpret(stmts []ast.Stmt) {
	if len(stmts) == 0 {
		return
	}
	for _, statememnt := range stmts {
		_, err := i.exec(statememnt) // WE DO NOT EVAL STATMENTS, WE EXECUTE THEM
		if err != nil {
			fmt.Println(fmt.Errorf("error: %v", err))
		}
	}
}

// VisitVarStmt handles the execution of a variable declaration statement.
// It evaluates the initializer expression if present, defines the variable
// in the current environment with its name and value, and returns any error
// encountered during evaluation. If no initializer is provided, the variable
// is defined with a nil value.
func (i *Interpreter) VisitVarStmt(stmt ast.VarStmt) error {
	var value any = nil
	if stmt.Initializer != nil {
		var err error = nil
		value, err = i.eval(stmt.Initializer)
		if err != nil {
			return err
		}
	}
	i.environment.Define(stmt.Name.Lexeme, value)
	return nil
}

// VisitVarExpr evaluates a variable expression by retrieving its value from the current environment.
// It takes an ast.Variable as input, attempts to get the value associated with the variable's name,
// and returns the value along with any error encountered during the lookup.
func (i *Interpreter) VisitVarExpr(expr ast.Variable) (any, error) {
	value, err := i.environment.Get(expr.Name)
	if err != nil {
		return nil, err
	}
	return value, nil
}

// VisitLiteral evaluates a literal expression and returns its value.
// It takes an ast.Literal as input and returns the value of the literal
// along with any potential error. Literal expressions represent constant
// values such as numbers, strings, or booleans in the abstract syntax tree.
func (i *Interpreter) VisitLiteral(expr ast.Literal) (any, error) {
	return expr.Value, nil
}

// VisitUnary evaluates a unary expression in the abstract syntax tree (AST).
// It performs a post-order evaluation of the operand and applies the unary operator.
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

// VisitExpressionStmt evaluates an expression statement by visiting its expression node.
// It takes an ExpressionStmt from the AST as input and returns the result of evaluating
// the expression along with any potential error encountered during evaluation.
func (i *Interpreter) VisitExpressionStmt(stmt ast.ExpressionStmt) (any, error) {
	return i.eval(stmt.Expression)
}

// VisitPrintStmt evaluates a PrintStmt node in the abstract syntax tree (AST)
// and prints the result of the evaluated expression to the standard output.
// It takes a PrintStmt as input, evaluates its Expression field, and formats
// the result using the stringify function before printing it.
// Returns nil for both the result and error as this function is primarily
// used for side effects (printing).
func (i *Interpreter) VisitPrintStmt(stmt ast.PrintStmt) (any, error) {
	value, _ := i.eval(stmt.Expression)
	fmt.Printf("%s", stringify(value))
	return nil, nil
}

// VisitGrouping evaluates a grouping expression by recursively evaluating
// the inner expression contained within the grouping. It returns the result
// of the evaluation or an error if the evaluation fails.
func (i *Interpreter) VisitGrouping(expr ast.Grouping) (any, error) {
	return i.eval(expr.Expression)
}

// eval evaluates the given AST expression by delegating the evaluation
// to the expression's Accept method. It returns the result of the evaluation
// along with any error encountered during the process.
func (i *Interpreter) eval(expr ast.Expr) (any, error) {
	return expr.Accept(i)
}

// exec executes the given statement by invoking its Accept method,
// passing the current Interpreter instance. It returns the result
// of the statement execution along with any potential error.
func (i *Interpreter) exec(stmt ast.Stmt) (any, error) {
	return stmt.Accept(i)
}

// VisitBinary evaluates a binary expression by visiting its left and right operands
// and applying the operator specified in the expression. It supports various operators
// such as arithmetic, comparison, logical, and string concatenation.
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
