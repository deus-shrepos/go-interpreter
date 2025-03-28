package ast

import "github.com/go-interpreter/internal/token"

type Visitor interface {
	VisitBinary(node Binary) (any, error)
	VisitGrouping(node Grouping) (any, error)
	VisitLiteral(node Literal) (any, error)
	VisitUnary(node Unary) (any, error)
}

type Expr interface {
	Accept(vistior Visitor) (any, error)
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (node Binary) Accept(visitor Visitor) (any, error) {
	return visitor.VisitBinary(node)
}

type Grouping struct {
	Expression Expr
}

func (node Grouping) Accept(visitor Visitor) (any, error) {
	return visitor.VisitGrouping(node)
}

type Literal struct {
	Value any
}

func (node Literal) Accept(visitor Visitor) (any, error) {
	return visitor.VisitLiteral(node)
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

func (node Unary) Accept(visitor Visitor) (any, error) {
	return visitor.VisitUnary(node)
}

type Stmt struct {
	Expression Expr
}

func (node Stmt) Accept(visitor Visitor) (any, error) {
	return visitor.VisitStmt(node)
}

type Print struct {
	Expression Expr
}

func (node Print) Accept(visitor Visitor) (any, error) {
	return visitor.VisitPrint(node)
}
