package ast

import "github.com/go-interpreter/internal/token"

type StmtVisitor interface {
	VisitExpressionStmt(node ExpressionStmt) (any, error)
	VisitPrintStmt(node PrintStmt) (any, error)
	VisitVarStmt(node VarStmt) (any, error)
	VisitBlockStmt(node Block) (any, error)
	VisitIfStmt(node IfStmt) (any, error)
	VisitWhileStmt(node WhileStmt) (any, error)
	VisitBreakStmt() (any, error)
	VisitContinueStmt() (any, error)
}

type Stmt interface {
	Accept(visitor StmtVisitor) (any, error)
}
type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

type BreakStmt struct {
	Value string
}

func (node BreakStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitBreakStmt()
}

type ContinueStmt struct{}

func (node ContinueStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitContinueStmt()
}

func (node WhileStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitWhileStmt(node)
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (node IfStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitIfStmt(node)
}

type Block struct {
	Statements []Stmt
}

func (node Block) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitBlockStmt(node)
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
