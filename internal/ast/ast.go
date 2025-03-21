package ast

import (
	"Crafting-interpreters/internal/token"
)

type Visitor interface {
	VisitBinary(node Binary) (any, error)
	VisitGrouping(node Grouping) (any, error)
	VisitLiteral(node Literal) (any, error)
	VisitUnary(node Unary) (any, error)
}

type Expr interface {
	Accept(visitor Visitor) (any, error)
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (node Binary) Accept(visitor Visitor) (any, error) {
	return visitor.VisitBinary(node) //nolint:wrapcheck
}

type Grouping struct {
	Expression Expr
}

func (node Grouping) Accept(visitor Visitor) (any, error) {
	return visitor.VisitGrouping(node) //nolint:wrapcheck
}

type Literal struct {
	Value any
}

func (node Literal) Accept(visitor Visitor) (any, error) {
	return visitor.VisitLiteral(node) //nolint:wrapcheck
}

type Unary struct {
	Right    Expr
	Operator token.Token
}

func (node Unary) Accept(visitor Visitor) (any, error) {
	return visitor.VisitUnary(node) //nolint:wrapcheck
}
