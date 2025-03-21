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
	Literal any
	Line    int
	Char    int
}

// ToString This convert the token to a representation; something like this:  1 + variable + VAR
func (token *Token) ToString() string {
	tokenType := strconv.Itoa(int(token.Type))
	literal, _ := token.Literal.(string)

	return strings.Join([]string{tokenType, " ", token.Lexeme, " ", literal}, "")
}

// String This will be used by the fmt package to print out to the standard out using this formatted string
func (token Token) String() string {
	return fmt.Sprintf("Token<Type=%v, Lexeme=%v, Literal=%v, Line=%v, Char=%v>",
		token.Type, token.Lexeme, token.Literal, token.Line, token.Char) //nolint:lll
}
