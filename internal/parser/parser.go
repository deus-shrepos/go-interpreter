package parser

import (
	"Crafting-interpreters/internal/ast"
	"Crafting-interpreters/internal/errors"
	"Crafting-interpreters/internal/token"
	"fmt"
)

/*
	Literals: Numbers, strings, Booleans, and nil
	Unary expressions: A prefix ! to perform a logical not, and - to negate a number
	Binary expressions: the infix arithmetic (+, -, #, /) and logic operators
					(==, !=, <, <=, >, >=)
	Parentheses: A pair of ( and a) wrapped around an expression

	--------------------------------------------
	expression -> literal
				 | unary
				 | binary
				 | grouping;
	literal    ->   NUMBER | STRING | "true" | "false" | "nil";
	grouping   ->  "(" expression ")"
	unary 	   ->   ( "-" | "!" ) expression;
	binary 	   ->   expression operator expression
	operator   -> 	"==" | "!=" | "<" | "<=" | ">" | ">=" |
					| "+" | "-" | "*" | "/";

	What about Syntactic ambiguity? How do we handle that?
	String => 6 / 3 - 1
		generates => T1 and T2
					 is T1 =/= T2? It should not be the case.
					 semantically makes a huge difference.
					 so we ought to take care of it.
		if syntax trees are different, therefore meaning is different.
	Handling is as simple as defining rules for precedence and associativity
	higher precedence -> binds tighter
	lower  precedence -> binds less tightly
	Associativity -> which operator to be evaluated first in the series of the
	same operator.
		left-associative (left-to-right) -> left to right evaluation
			5 - 3 - 1 => (5 - 3) - 1
		right-associative (right-to-left) -> left to right evaluation
		 	a = b = c => a = (b = c)

		----------------------------------------------
		Name				Operators 		Associates
		Equality			== !=			Left
		Comparison			> >= < <=		Left
		Term				- +				Left
		Factor				/ * 			Left
		Unary				! - 			Right

		expression  ->
		equality    ->
		comparison  ->
		term		->
		factor 		-> factor ("/" | "*") unary | unary
		unary 		-> ("!" | "-") unary | primary
		primary 	-> Number | string

	 	example: 1 * 2 / 3
		we have to do left-associative parsing in this case since
		* and / have equal precedence, and we will recurse from left to right
		to avoid any confusion.

		factor => factor "/" unary
		factor => (factor "*" unary) "/" unary
		factor => ( (unary "*" unary) "/" unary)
		...
		unary  => (((primary "*" primary) "/" primary))
		..,
		primary => (Number "*" Number) "/" Number

		Revised:
		expression  -> equality
		equality    -> comparison ( "!=" | "==" ) comparison )*
		comparison  -> term ( ( ">" | ">=" | "<" | "<=") term )*
		term		-> factor ( ( "-" | "+" ) factor )*
		factor 		-> unary ( ( "/" | "*") unary) *
		unary 		-> ("!" | "-") unary | primary
		primary 	-> Number | string | "true" | "false" | "nil" |
						"(" expression ")"

		Recursive Descent Parsing
		-------------------------

		Combinations of L & R: LL(K), LR(1), LALR, or RDP
		It is a top-down parser as it starts from the top or the outermost
		grammar rule ( like expression ) and works its way down into the nested
		subexpressions before finally reaching the leaves of the syntax tree.

		Grammar Notion 				Code Repr
		-------------------------------------------------
		Terminal 					code -> match/consume
		Non-terminal 				call -> rule's func
		|							if/switch
		* or + 						while/for loop
		?							if
*/

// Parser is a recursive descent parser for the Lox language.
// It takes a list of tokens and produces an abstract syntax tree.
type Parser struct {
	Tokens  []token.Token
	Current int
	Error   errors.Error
}

// Parser initializes a new parser with the given list of tokens.
func (parser *Parser) Parser(tokens []token.Token) {
	parser.Tokens = tokens
	parser.Current = 0
}

func (parser *Parser) Parse() (ast.Expr, error) {
	expr, err := parser.Expression()
	if err != nil {
		return nil, err
	}
	return expr, nil
}

// Expression parses an expression from the list of tokens.
// It returns the root node of the abstract syntax tree.
// expression -> equality
func (parser *Parser) Expression() (ast.Expr, error) {
	eql, err := parser.Equality()
	if err != nil {
		return nil, err
	}
	return eql, nil
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
		expr, err := parser.Expression()
		if err != nil {
			return nil, err
		}
		_, err = parser.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return ast.Grouping{Expression: expr}, nil
	default:
		// We probaby don't want to panic here because we are syncing the parser
		peek := parser.peek()
		previous := parser.previous()
		return nil, parser.ParserError(peek, fmt.Sprintf("Unexpected token '%v' found after '%s' \n", peek.Lexeme, previous.Lexeme))

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
func (parser *Parser) consume(type_ token.TokenType, message string) (token.Token, error) {
	if parser.check(type_) {
		return parser.advance(), nil
	}
	err := parser.ParserError(parser.peek(), message)
	return token.Token{}, err
}

// ParserError The function just reports and returns the error
// The caller can handle the error and decide what to do with it
func (parser *Parser) ParserError(token token.Token, message string) error {
	parser.Error.ProgramError(token.Line, message)
	return fmt.Errorf("parser Error Occurred. Exiting")
}
