package printer

import (
	"Crafting-interpreters/internal/ast"
	"fmt"
)

// PrintAST is a visitor implementation for converting abstract syntax trees into their string representation.
type PrintAST struct{}

// VisitBinary generates a string representation of a binary expression by recursively visiting its left, right, and operator.
func (printer PrintAST) VisitBinary(node ast.Binary) interface{} {
	return fmt.Sprintf("(%s %s %s)",
		node.Left.Accept(printer),
		node.Right.Accept(printer),
	)
}

// VisitGrouping creates a string representation of a grouping expression by recursively visiting its inner expression.
func (printer PrintAST) VisitGrouping(node ast.Grouping) interface{} {
	return fmt.Sprintf("(group %s)", node.Expression.Accept(printer))
}

// VisitLiteral generates a string representation of a literal expression based on its value.
func (printer PrintAST) VisitLiteral(node ast.Literal) interface{} {
	return fmt.Sprintf("%v", node.Value)
}

// VisitUnary generates a string representation of a unary expression by visiting its operator and operand.
func (printer PrintAST) VisitUnary(node ast.Unary) interface{} {
	return fmt.Sprintf("(%s %s)",
		node.Operator,
		node.Right.Accept(printer),
	)
}

// // TODO: Need to print the AST in the REPL
func (printer PrintAST) Print(expression ast.Expr) string {
	return expression.Accept(printer).(string)
}
