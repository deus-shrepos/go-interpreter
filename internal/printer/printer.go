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

// VisitBinary generates a string representation of a binary expression.
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

// VisitGrouping creates a string representation of a grouping expression.
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

// VisitLiteral generates a string representation of a literal expression.
func (printer *PrintAST) VisitLiteral(node ast.Literal) (interface{}, error) {
	return fmt.Sprintf("%sLiteral(%v)",
		strings.Repeat("  ", printer.indentation),
		node.Value,
	), nil
}

// VisitUnary generates a string representation of a unary expression.
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

// VisitVariable prints a variable expression.
func (printer *PrintAST) VisitVariable(node ast.Variable) (interface{}, error) {
	return fmt.Sprintf("%sVariable(%s)",
		strings.Repeat("  ", printer.indentation),
		node.Name.Lexeme,
	), nil
}

// VisitAssign prints an assignment expression.
func (printer *PrintAST) VisitAssign(node ast.Assign) (interface{}, error) {
	printer.indentation++
	value, _ := node.Value.Accept(printer)
	printer.indentation--
	return fmt.Sprintf("%sAssign(\n%s%s\n%s\n%s)",
		strings.Repeat("  ", printer.indentation),
		strings.Repeat("  ", printer.indentation+1),
		node.Name.Lexeme,
		value.(string),
		strings.Repeat("  ", printer.indentation),
	), nil
}

// VisitLogical prints a logical expression.
func (printer *PrintAST) VisitLogical(node ast.Logical) (interface{}, error) {
	printer.indentation++
	left, _ := node.Left.Accept(printer)
	right, _ := node.Right.Accept(printer)
	printer.indentation--
	return fmt.Sprintf("%sLogical(\n%s\n%s%s\n%s\n%s)",
		strings.Repeat("  ", printer.indentation),
		left.(string),
		strings.Repeat("  ", printer.indentation+1),
		node.Operator.Lexeme,
		right.(string),
		strings.Repeat("  ", printer.indentation),
	), nil
}

// VisitExpressionStmt prints an expression statement.
func (printer *PrintAST) VisitExpressionStmt(node ast.ExpressionStmt) (interface{}, error) {
	printer.indentation++
	expr, _ := node.Expression.Accept(printer)
	printer.indentation--
	return fmt.Sprintf("%sExpressionStmt(\n%s\n%s)",
		strings.Repeat("  ", printer.indentation),
		expr.(string),
		strings.Repeat("  ", printer.indentation),
	), nil
}

// VisitPrintStmt prints a print statement.
func (printer *PrintAST) VisitPrintStmt(node ast.PrintStmt) (interface{}, error) {
	printer.indentation++
	expr, _ := node.Expression.Accept(printer)
	printer.indentation--
	return fmt.Sprintf("%sPrintStmt(\n%s\n%s)",
		strings.Repeat("  ", printer.indentation),
		expr.(string),
		strings.Repeat("  ", printer.indentation),
	), nil
}

// VisitVarStmt prints a variable declaration statement.
func (printer *PrintAST) VisitVarStmt(node ast.VarStmt) (interface{}, error) {
	printer.indentation++
	var initStr string
	if node.Initializer != nil {
		init, _ := node.Initializer.Accept(printer)
		initStr = init.(string)
	} else {
		initStr = "nil"
	}
	printer.indentation--
	return fmt.Sprintf("%sVarStmt(%s,\n%s\n%s)",
		strings.Repeat("  ", printer.indentation),
		node.Name.Lexeme,
		initStr,
		strings.Repeat("  ", printer.indentation),
	), nil
}

// VisitBlockStmt prints a block of statements.
func (printer *PrintAST) VisitBlockStmt(node ast.Block) (interface{}, error) {
	printer.indentation++
	var stmts []string
	for _, stmt := range node.Statements {
		s, _ := stmt.Accept(printer)
		stmts = append(stmts, s.(string))
	}
	printer.indentation--
	return fmt.Sprintf("%sBlock(\n%s\n%s)",
		strings.Repeat("  ", printer.indentation),
		strings.Join(stmts, "\n"),
		strings.Repeat("  ", printer.indentation),
	), nil
}

// VisitReturnStmt prints a return statement.
func (printer *PrintAST) VisitReturnStmt(node ast.ReturnStmt) (interface{}, error) {
	printer.indentation++
	returnExpr, _ := node.Value.Accept(printer)
	return fmt.Sprintf("%sReturnStmt(\n%s%s\n%s%s",
		strings.Repeat(" ", printer.indentation),
		node.Keyword,
		strings.Repeat(" ", printer.indentation),
		returnExpr,
		strings.Repeat(" ", printer.indentation),
	), nil

}

// VisitIfStmt prints an if statement.
func (printer *PrintAST) VisitIfStmt(node ast.IfStmt) (interface{}, error) {
	printer.indentation++
	cond, _ := node.Condition.Accept(printer)
	thenBranch, _ := node.ThenBranch.Accept(printer)
	var elseBranchStr string
	if node.ElseBranch != nil {
		elseBranch, _ := node.ElseBranch.Accept(printer)
		elseBranchStr = fmt.Sprintf("\n%sElse:\n%s", strings.Repeat("  ", printer.indentation+1), elseBranch.(string))
	}
	printer.indentation--
	return fmt.Sprintf("%sIfStmt(\n%s\n%sThen:\n%s%s\n%s)",
		strings.Repeat("  ", printer.indentation),
		cond.(string),
		strings.Repeat("  ", printer.indentation+1),
		thenBranch.(string),
		elseBranchStr,
		strings.Repeat("  ", printer.indentation),
	), nil
}

// VisitFunctionStmt prints a function declaration.
func (printer *PrintAST) VisitFunctionStmt(node ast.Function) (interface{}, error) {
	printer.indentation++
	var params []string
	for _, param := range node.Parameters {
		params = append(params, param.Lexeme)
	}
	var bodyStmts []string
	for _, stmt := range node.Body {
		s, _ := stmt.Accept(printer)
		bodyStmts = append(bodyStmts, s.(string))
	}
	printer.indentation--
	return fmt.Sprintf("%sFunction(%s,\n%s\n%s)",
		strings.Repeat("  ", printer.indentation),
		strings.Join(params, ", "),
		strings.Join(bodyStmts, "\n"),
		strings.Repeat("  ", printer.indentation),
	), nil
}

// VisitFunctionCall prints a function call expression.
func (printer *PrintAST) VisitFunctionCall(node ast.Call) (interface{}, error) {
	printer.indentation++
	callee, _ := node.Callee.Accept(printer)
	var args []string
	for _, arg := range node.Args {
		a, _ := arg.Accept(printer)
		args = append(args, a.(string))
	}
	printer.indentation--
	return fmt.Sprintf("%sCall(\n%s\n%sArguments:\n%s\n%s)",
		strings.Repeat("  ", printer.indentation),
		callee.(string),
		strings.Repeat("  ", printer.indentation+1),
		strings.Join(args, "\n"),
		strings.Repeat("  ", printer.indentation),
	), nil
}

// TODO: FIX IT
func (printer *PrintAST) VisitFunctionExpression(node ast.FunctionExpr) (interface{}, error) {
	return nil, nil
}

// VisitWhileStmt prints a while statement.
func (printer *PrintAST) VisitWhileStmt(node ast.WhileStmt) (interface{}, error) {
	printer.indentation++
	cond, _ := node.Condition.Accept(printer)
	body, _ := node.Body.Accept(printer)
	printer.indentation--
	return fmt.Sprintf("%sWhileStmt(\n%s\n%sBody:\n%s\n%s)",
		strings.Repeat("  ", printer.indentation),
		cond.(string),
		strings.Repeat("  ", printer.indentation+1),
		body.(string),
		strings.Repeat("  ", printer.indentation),
	), nil
}

// VisitBreakStmt prints a break statement.
func (printer *PrintAST) VisitBreakStmt() (interface{}, error) {
	return fmt.Sprintf("%sBreakStmt", strings.Repeat("  ", printer.indentation)), nil
}

// VisitContinueStmt prints a continue statement.
func (printer *PrintAST) VisitContinueStmt() (interface{}, error) {
	return fmt.Sprintf("%sContinueStmt", strings.Repeat("  ", printer.indentation)), nil
}

// VisitNodes visits both expressions and statements
func (printer *PrintAST) VisitNodes(node interface{}) string {
	var result interface{}
	switch n := node.(type) {
	case ast.Expr:
		result, _ = n.Accept(printer)
	case ast.Stmt:
		result, _ = n.Accept(printer)
	default:
		fmt.Println(n)
		return "Unknown node type"
	}
	return result.(string)
}

// PrintAll prints all the AST nodes
func (printer *PrintAST) PrintAll(stmts []ast.Stmt) {
	astStringBuilder := strings.Builder{}
	for _, s := range stmts {
		astStringBuilder.WriteString(printer.VisitNodes(s))
		astStringBuilder.WriteString("\n")
	}
	fmt.Println(astStringBuilder.String())
}
