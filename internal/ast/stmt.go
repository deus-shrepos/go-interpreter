package ast

import "github.com/go-interpreter/internal/token"

type StmtVisitor interface {
	VisitExpressionStmt(node ExpressionStmt) (any, error)
	VisitPrintStmt(node PrintStmt) (any, error)
	VisitVarStmt(node VarStmt) (any, error)
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

type VarStmt struct {
	Name        token.Token
	Initializer Expr
}

func (node VarStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitVarStmt(node)
}
