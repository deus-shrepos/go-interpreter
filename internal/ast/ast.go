package ast

import (
	"Crafting-interpreters/internal/token"
)

type Visitor interface {
	VisitBinary(node Binary) (interface{}, error)
	VisitGrouping(node Grouping) (interface{}, error)
	VisitLiteral(node Literal) (interface{}, error)
	VisitUnary(node Unary) (interface{}, error)
}

type Expr interface {
	Accept(visitor Visitor) (interface{}, error)
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (node Binary) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitBinary(node)
}

type Grouping struct {
	Expression Expr
}

func (node Grouping) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitGrouping(node)
}

type Literal struct {
	Value interface{}
}

func (node Literal) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitLiteral(node)
}

type Unary struct {
	Right    Expr
	Operator token.Token
}

func (node Unary) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitUnary(node)
}
