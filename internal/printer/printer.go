package printer

import (
	"fmt"
	"strings"

	"github.com/go-interpreter/internal/ast"
)

// PrintAST is a visitor implementation for converting abstract syntax trees into their string representation.
type PrintAST struct {
	indentation int
}

// VisitBinary generates a string representation of a binary
// expression by recursively visiting its left, right, and operator.
func (printer *PrintAST) VisitBinary(node ast.Binary) (interface{}, error) {
	printer.indentation++
	left, _ := node.Left.Accept(printer)
	right, _ := node.Right.Accept(printer)
	printer.indentation--
	return fmt.Sprintf("%sBinary(\n%s\n%s%s %s\n%s\n%s)",
		strings.Repeat("  ", printer.indentation),
		left.(string),
		strings.Repeat("  ", printer.indentation+1),
		node.Operator.Lexeme,
		strings.Repeat("  ", printer.indentation+1),
		right.(string),
		strings.Repeat("  ", printer.indentation),
	), nil
}

// VisitGrouping creates a string representation of a grouping expression by recursively visiting its inner expression.
func (printer *PrintAST) VisitGrouping(node ast.Grouping) (interface{}, error) {
	printer.indentation++
	expr, _ := node.Expression.Accept(printer)
	printer.indentation--
	return fmt.Sprintf("%sGrouping(\n%s\n%s)",
		strings.Repeat("  ", printer.indentation),
		expr.(string),
		strings.Repeat("  ", printer.indentation),
	), nil
}

// VisitLiteral generates a string representation of a literal expression based on its value.
func (printer *PrintAST) VisitLiteral(node ast.Literal) (interface{}, error) {
	return fmt.Sprintf("%sLiteral(%v)",
		strings.Repeat("  ", printer.indentation),
		node.Value,
	), nil
}

// VisitUnary generates a string representation of a unary expression by visiting its operator and operand.
func (printer *PrintAST) VisitUnary(node ast.Unary) (interface{}, error) {
	printer.indentation++
	right, _ := node.Right.Accept(printer)
	printer.indentation--
	return fmt.Sprintf("%sUnary(\n%s%s\n%s)",
		strings.Repeat("  ", printer.indentation),
		strings.Repeat("  ", printer.indentation+1),
		node.Operator.Lexeme,
		right.(string),
	), nil
}

func (printer *PrintAST) VisitVariable(node ast.Variable) (interface{}, error) {
	printer.indentation++
	variableName, _ := node.Accept(printer)
	printer.indentation--
	return fmt.Sprintf("%Variable(\n%s%s\n)",
		strings.Repeat("  ", printer.indentation),
		strings.Repeat("  ", printer.indentation+1),
		variableName,
	), nil
}

func (printer *PrintAST) Print(expression ast.Expr) string {
	expr, _ := expression.Accept(printer)
	return expr.(string)
}
