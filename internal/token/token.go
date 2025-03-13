package token

import (
	"fmt"
	"strconv"
	"strings"
)

// Token represents a single unit of lexical information in the source code.
// It includes its type, lexeme, literal value, and the line on which it appears.
type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

// Init initializes a Token with the given type, lexeme, literal value, and line number where it appears.
func (token *Token) Init(tokenType TokenType, Lexeme string, literal interface{}, line int) {
	token.Type = tokenType
	token.Lexeme = Lexeme
	token.Literal = literal
	token.Line = line
}

// ToString This convert the token to a representation; something like this:  1 + variable + VAR
func (token *Token) ToString() string {
	tokenType := strconv.Itoa(int(token.Type))
	literal, _ := token.Literal.(string)
	return strings.Join([]string{tokenType, " ", token.Lexeme, " ", literal}, "")
}

// String This will be used by the fmt package to print out to the standard out using this formatted string
func (token Token) String() string {
	return fmt.Sprintf("Token<Type=%v, Lexeme=%v, Literal=%v, Line=%v>", token.Type, token.Lexeme, token.Literal, token.Line)
}
