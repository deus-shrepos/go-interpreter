package ast

import "crafting-interpreters/lox"

type Expr interface{}
type Visitor interface{
	VisitBinary(node Binary) interface{}
	VisitGrouping(node Grouping) interface{}
	VisitLiteral(node Literal) interface{}
	VisitUnary(node Unary) interface{}
	
}
type Binary struct {
	Left Expr
	Operator lox.Token
	Right Expr
}
func (node Binary) Accept(visitor Visitor) interface{} {
	return visitor.VisitBinary(node)
}

type Grouping struct {
	Expression Expr
}
func (node Grouping) Accept(visitor Visitor) interface{} {
	return visitor.VisitGrouping(node)
}

type Literal struct {
	Value interface{}
}
func (node Literal) Accept(visitor Visitor) interface{} {
	return visitor.VisitLiteral(node)
}

type Unary struct {
	Operator lox.Token
	Right Expr
}
func (node Unary) Accept(visitor Visitor) interface{} {
	return visitor.VisitUnary(node)
}

