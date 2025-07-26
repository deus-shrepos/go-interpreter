package ast

import (
	"github.com/go-interpreter/internal/token"
)

type ExprVisitor interface {
	VisitBinary(node Binary) (any, error)
	VisitGrouping(node Grouping) (any, error)
	VisitLiteral(node Literal) (any, error)
	VisitUnary(node Unary) (any, error)
	VisitVariable(node Variable) (any, error)
	VisitAssign(node Assign) (any, error)
}

type Expr interface {
	Accept(vistior ExprVisitor) (any, error)
}

type Assign struct {
	Name  token.Token
	Value Expr
}

func (node Assign) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitAssign(node)
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (node Binary) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitBinary(node)
}

type Grouping struct {
	Expression Expr
}

func (node Grouping) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitGrouping(node)
}

type Literal struct {
	Value any
}

func (node Literal) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitLiteral(node)
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

func (node Unary) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitUnary(node)
}

type Variable struct {
	Name token.Token
}

func (node Variable) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitVariable(node)
}
