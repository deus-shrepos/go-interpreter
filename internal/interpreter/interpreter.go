package interpreter

import (
	_ "errors"
	"fmt"
	"strconv"

	"github.com/go-interpreter/internal/ast"
	"github.com/go-interpreter/internal/errors"
	"github.com/go-interpreter/internal/token"
)

// Interpreter represents the core structure for the interpreter.
// It is responsible for executing and evaluating code based on the
// implemented logic and rules of the interpreter.
type Interpreter struct {
	environment *Environment
}

func NewInterpreter() Interpreter {
	return Interpreter{
		environment: NewEnvironment(nil),
	}
}

// Interpret executes a series of statements provided as input.
// It iterates over each statement, executing them one by one using the exec method.
// If an error occurs during the execution of a statement, it logs the error to the console.
func (i *Interpreter) Interpret(stmts []ast.Stmt) error {
	if len(stmts) == 0 {
		return nil
	}
	for _, statement := range stmts {
		// We panic if the interpreter parses a NIL (because that is parsing error).
		// If we execute a NIL that will make the whole goroutine panic.
		// This would essentially make sure we have an early exit.
		if statement == nil {
			return fmt.Errorf("error: Interpreter panic. Exiting program")
		}
		_, err := i.exec(statement) // WE DO NOT EVAL STATEMENTS, WE EXECUTE THEM
		if err != nil {
			return fmt.Errorf("error: %v", err)
		}
	}
	fmt.Println("") // To get rid of that annoying "%" in the terminal
	return nil
}

// VisitVarStmt handles the execution of a variable declaration statement.
// It evaluates the initializer expression if present, defines the variable
// in the current Environment with its name and value, and returns any error
// encountered during evaluation. If no initializer is provided, the variable
// is defined with a nil value.
func (i *Interpreter) VisitVarStmt(stmt ast.VarStmt) (any, error) {
	var value any = nil
	if stmt.Initializer != nil {
		var err error = nil
		value, err = i.eval(stmt.Initializer)
		if err != nil {
			return nil, err
		}
	}
	i.environment.Define(stmt.Name.Lexeme, value)
	return nil, nil
}

// VisitVariable VisitVarExpr evaluates a variable expression by retrieving its value from the current Environment.
// It takes an ast.Variable as input, attempts to get the value associated with the variable's name,
// and returns the value along with any error encountered during the lookup.
func (i *Interpreter) VisitVariable(expr ast.Variable) (any, error) {
	value, err := i.environment.Get(expr.Name)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (i *Interpreter) VisitIfStmt(stmt ast.IfStmt) (any, error) {
	evaluatedExpr, err := i.eval(stmt.Condition) // Evaluate the if-condition
	if err != nil {
		return nil, err
	}
	if IsTruthy(evaluatedExpr) {
		_, err = i.exec(stmt.ThenBranch)
		if err != nil {
			return nil, err
		}
	} else if stmt.ElseBranch != nil {
		_, err = i.exec(stmt.ElseBranch)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// VisitAssign handles assignment expressions in the AST.
// It evaluates the right-hand side value, then assigns it to the variable in the current Environment.
// Returns the assigned value and any error encountered during evaluation or assignment.
func (i *Interpreter) VisitAssign(expr ast.Assign) (any, error) {
	value, err := i.eval(expr.Value)
	if err != nil {
		return nil, err
	}
	err = i.environment.Assign(expr.Name, value)
	if err != nil {
		return nil, err
	}
	return value, err
}

// VisitBlockStmt VisitBlock executes a block statement by creating a new environment scope.
// It runs each statement in the block within this new environment, ensuring
// that variables declared inside the block do not affect the outer environment.
// Returns nil and any error encountered during execution.
func (i *Interpreter) VisitBlockStmt(blockStmt ast.Block) (any, error) {
	_, err := i.execBlock(blockStmt.Statements, NewEnvironment(i.environment))
	if err != nil {
		return nil, err
	}
	return nil, nil
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
	fmt.Print(stringify(value))
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

// execBlock executes a list of statements within a new environment scope.
// It temporarily replaces the interpreter's current environment with the provided one,
// executes each statement in the block, and restores the previous environment afterward.
// If any statement returns an error, execution stops and the error is returned.
// Returns nil and any error encountered during execution.
func (i *Interpreter) execBlock(stmts []ast.Stmt, environment *Environment) (any, error) {
	previous := i.environment
	i.environment = environment
	for _, stmt := range stmts {
		_, err := i.exec(stmt)
		if err != nil {
			return nil, err
		}
	}
	i.environment = previous
	return nil, nil
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
				return leftValue + rightValue, nil
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

	switch v := object.(type) {
	case nil:
		return false
	case string:
		return len(v) > 0
	case int:
		return v > 0 || v < 0
	case bool:
		return v
	case float64:
		return v > 0.0
	default:
		return false
	}
}

func stringify(object any) string {
	switch value := object.(type) {
	case float64:
		return strconv.FormatFloat(value, 'g', -1, 64)
	case string:
		return value
	case nil:
		return ""
	default:
		return fmt.Sprint(value)
	}
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
