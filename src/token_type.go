package src

// TokenType represents the category or type of a token in a lexical analysis process, such as operators, keywords, or literals.
type TokenType int

const (

	// LEFT_PAREN SINGLE TOKENS
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// BANG ONE OR TWO TOKENS (BINARY OPERATORS MORE LIKE???)
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// IDENTIFIER LITERALS (WHATEVER IT IS, IT IS)
	IDENTIFIER
	STRING
	NUMBER

	// LANGUAGE KEYWORDS
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR

	// WHILE represents the 'while' keyword in the language, used for defining looping constructs within the code.
	WHILE

	// MISC
	EOF
)

// keywords is a map linking string representations of language keywords to their corresponding TokenType values.
var keywords = map[string]TokenType{
	"and":    AND,
	"or":     OR,
	"class":  CLASS,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
	"this":   THIS,
	"else":   ELSE,
}
