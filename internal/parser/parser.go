package parser

import (
	"fmt"

	"github.com/go-interpreter/internal/ast"
	"github.com/go-interpreter/internal/errors"
	"github.com/go-interpreter/internal/token"
)

// Parser is a recursive descent parser for the Lox language.
// It takes a list of tokens and produces an abstract syntax tree.
type Parser struct {
	Tokens  []token.Token
	Current int
}

// Parser initializes a new parser with the given list of tokens.
func (parser *Parser) Parser(tokens []token.Token) {
	parser.Tokens = tokens
	parser.Current = 0
}

func (parser *Parser) Parse() (ast.Expr, error) {
	expr, err := parser.Expression()
	if err != nil {
		fmt.Print(fmt.Errorf("%w", err))
		return nil, err
	}
	return expr, nil
}

// Expression parses an expression from the list of tokens.
// It returns the root node of the abstract syntax tree.
// expression -> equality
func (parser *Parser) Expression() (ast.Expr, error) {
	eq, err := parser.Equality()
	if err != nil {
		return nil, err
	}
	return eq, nil
}

// Equality parses an equality expression from the list of tokens.
// It returns the root node of the abstract syntax tree.
func (parser *Parser) Equality() (ast.Expr, error) {
	expr, err := parser.Comparison()
	if err != nil {
		return nil, err
	}
	for parser.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := parser.previous()
		right, err := parser.Comparison()
		if err != nil {
			return nil, err
		}
		expr = ast.Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

// Comparison It parses a comparison expression from the list of tokens.
// It returns the root node of the abstract syntax tree.
func (parser *Parser) Comparison() (ast.Expr, error) {
	expr, err := parser.Term()
	if err != nil {
		return nil, err
	}
	for parser.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := parser.previous()
		right, err := parser.Term()
		if err != nil {
			return nil, err
		}
		expr = ast.Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

// Term parses a term expression from  list of tokens.
// It returns the root node of the abstract syntax tree.
func (parser *Parser) Term() (ast.Expr, error) {
	expr, err := parser.Factor()
	if err != nil {
		return nil, err
	}
	for parser.match(token.MINUS, token.PLUS) {
		operator := parser.previous()
		right, err := parser.Factor()
		if err != nil {
			return nil, err
		}
		expr = ast.Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

// Factor parses a factor expression from the list of tokens.
// It returns the root node of the abstract syntax tree.
func (parser *Parser) Factor() (ast.Expr, error) {
	expr, err := parser.Unary()
	if err != nil {
		return nil, err
	}
	for parser.match(token.SLASH, token.STAR) {
		operator := parser.previous()
		right, err := parser.Unary()
		if err != nil {
			return nil, err
		}
		expr = ast.Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

// Unary parses a unary expression from the list of tokens.
// It returns the root node of the abstract syntax tree.
func (parser *Parser) Unary() (ast.Expr, error) {
	if parser.match(token.BANG, token.MINUS) {
		operator := parser.previous()
		right, err := parser.Unary()
		if err != nil {
			return nil, err
		}
		return ast.Unary{Operator: operator, Right: right}, nil
	}
	primary, err := parser.Primary()
	if err != nil {
		return nil, err
	}
	return primary, nil
}

// Primary parses a primary expression from the list of tokens.
// It returns the root node of the abstract syntax tree.
func (parser *Parser) Primary() (ast.Expr, error) {
	switch {
	case parser.match(token.FALSE):
		return ast.Literal{Value: false}, nil
	case parser.match(token.TRUE):
		return ast.Literal{Value: true}, nil
	case parser.match(token.NIL):
		return ast.Literal{Value: nil}, nil
	case parser.match(token.NUMBER, token.STRING):
		return ast.Literal{Value: parser.previous().Literal}, nil
	case parser.match(token.LEFT_PAREN):
		expr, e := parser.Expression()
		if e != nil {
			fmt.Println(fmt.Errorf("%v", e))
		}
		_, err := parser.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return ast.Grouping{Expression: expr}, nil
	default:
		// We probaby don't want to panic here because we are syncing the parser
		// We will catch it in parser.match(token.LEFT_PAREN) and report it back to
		// the stdout
		peek := parser.peek()
		return nil, errors.ExecutionError{
			Type:    errors.PARSER_ERROR,
			Line:    peek.Line,
			Where:   peek.Char,
			Message: fmt.Sprintf("Unexpected token '%v'", peek.Lexeme),
		}

	}
}

// Comparison parses a comparison expression from the list of tokens.
// It returns the root node of the abstract syntax tree.
func (parser *Parser) match(types ...token.TokenType) bool {
	for _, tokenType := range types {
		if parser.check(tokenType) {
			parser.advance()
			return true
		}
	}
	return false
}

// This is a helper function to check if the current token is of the given type.
func (parser *Parser) check(type_ token.TokenType) bool {
	if parser.isAtEnd() {
		return false
	}
	return parser.peek().Type == type_
}

// This is a helper function to advance the parser to the next token.
func (parser *Parser) advance() token.Token {
	if !parser.isAtEnd() {
		parser.Current++
	}
	return parser.previous()
}

// This is a helper function to match the type at the end of the list of tokens.
// If it is at the end, we return the null character.
func (parser *Parser) isAtEnd() bool {
	return parser.peek().Type == token.EOF
}

// This is a helper function to peek at the end of the string and return it.
func (parser *Parser) peek() token.Token {
	return parser.Tokens[parser.Current]
}

// This is a helper function to return the previous token.
func (parser *Parser) previous() token.Token {
	return parser.Tokens[parser.Current-1]
}

// This function consumer or otherwise it throws an error
// It consumes the token if it is of the given type.
// If it is not, it throws an error with the given message.
// The caller can handle the error and decide what to do with it.
func (parser *Parser) consume(type_ token.TokenType, message string) (token.Token, error) {
	if parser.check(type_) {
		return parser.advance(), nil
	}
	return token.Token{}, errors.ExecutionError{
		Type:    errors.PARSER_ERROR,
		Line:    parser.peek().Line,
		Where:   parser.peek().Char,
		Message: message,
	}
}
