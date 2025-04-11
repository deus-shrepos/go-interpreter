package ast

type StmtVisitor interface {
	VisitExpressionStmt(node ExpressionStmt) (any, error)
	VisitPrintStmt(node PrintStmt) (any, error)
}

type Stmt interface {
	Accept(vistior StmtVisitor) (any, error)
}

type ExpressionStmt struct {
	Expression Expr
}

func (node ExpressionStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitExpressionStmt(node)
}

type PrintStmt struct {
	Expression Expr
}

func (node PrintStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitPrintStmt(node)
}
