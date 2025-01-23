package src

import (
	"strconv"
)

/*
	Lexemes (literal) -> tokens (partially semantic???)
	Regex Machine -> Grammar that match a pattern of (sub)strings
	We have a scanner "class"(go doesn't do classes)
		=> Scan the text
		=> identify the tokens and wrap them in a token abstraction
		=> add the token to the token stack
		=> repeat
		=> ???
		=> profit??
*/

// TokenScanner Basically a TokenScanner that keeps the information
// of the scanned tokens in a stack. It also tracks the state of the
// scanning; such as the start and current(or it could be the end if current = EOF)
type TokenScanner struct {
	Source  string
	Tokens  []Token
	Start   int
	Current int
	Line    int
	Lox     Lox
}

// Init This initializes the source code
func (scanner *TokenScanner) Init(source string) {
	scanner.Source = source
	scanner.Start = 0
	scanner.Current = 0
	scanner.Line = 0

	// In case we want to use some Lox functionality - some static classes and stuff
	scanner.Lox = Lox{}
}

// ScanTokens scans the source code and produces a list of tokens based on the language grammar.
func (scanner *TokenScanner) ScanTokens() []Token {
	for !scanner.isAtEnd() {
		scanner.Start = scanner.Current
		scanner.ScanToken()
	}
	scanner.Tokens = append(scanner.Tokens, Token{EOF, "", nil, scanner.Line})
	return scanner.Tokens
}

// ScanToken reads the next character in the source and determines the appropriate token type to add to the token list.
func (scanner *TokenScanner) ScanToken() {
	c := scanner.advance()

	switch c {
	case "(":
		scanner.AddToken(LEFT_PAREN)
		break
	case ")":
		scanner.AddToken(RIGHT_PAREN)
		scanner.AddToken(SEMICOLON)
		break
	case "{":
		scanner.AddToken(LEFT_BRACE)
		break
	case "}":
		scanner.AddToken(RIGHT_BRACE)
		scanner.AddToken(SEMICOLON)
		break
	case ",":
		scanner.AddToken(COMMA)
		break
	case ".":
		scanner.AddToken(DOT)
		break
	case "-":
		scanner.AddToken(MINUS)
		break
	case "+":
		scanner.AddToken(PLUS)
		break
	case ";":
		scanner.AddToken(SEMICOLON)
		break
	case "*":
		scanner.AddToken(STAR)
		break
	case "!":
		bang := BANG
		if scanner.match("=") {
			bang = BANG_EQUAL
		}
		scanner.AddToken(bang)
	case "=":
		equal := EQUAL
		if scanner.match("=") {
			equal = EQUAL_EQUAL
		}
		scanner.AddToken(equal)
		break
	case "<":
		lessEqual := LESS
		if scanner.match("=") {
			lessEqual = LESS_EQUAL
		}
		scanner.AddToken(lessEqual)
		break
	case ">":
		greaterEqual := GREATER
		if scanner.match("=") {
			greaterEqual = GREATER_EQUAL
		}
		scanner.AddToken(greaterEqual)
		break
	case "/":
		if scanner.match("/") {
			// This is a comment - keep going until you reach at the end of the line
			for scanner.peek() != "\n" && !scanner.isAtEnd() {
				scanner.advance()
			}
		} else {
			scanner.AddToken(SLASH)
		}
		break
	case " ":
	case "\r":
	case "\t":
		break
	case "\n":
		scanner.Line++
		break
	case "\"":
		scanner.string()
		break
	default:
		if isDigit(c) {
			scanner.number()
		} else if isAlpha(c) {
			scanner.identifier()
		} else {
			scanner.Lox.ProgramError(scanner.Line, "Unexpected Error")
		}
		break
	}
}

// string scans a string literal, handles line tracking for multi-line strings, and adds the token to the token list.
func (scanner *TokenScanner) string() {
	for scanner.peek() != "\"" && !scanner.isAtEnd() {
		if scanner.peek() == "\n" {
			scanner.Line++
			scanner.advance()
		}
		if scanner.isAtEnd() {
			scanner.Lox.ProgramError(scanner.Line, "Unterminated string")
			return
		}
		scanner.advance()
		value := scanner.Source[scanner.Start+1 : scanner.Current-1]
		scanner.addToken(STRING, value)
	}
}

// number scans a numeric literal, supports fractional parts, and adds the token with its parsed value to the token list.
func (scanner *TokenScanner) number() {
	for isDigit(scanner.peek()) {
		scanner.advance()
	}
	// Look for a fractional part
	if scanner.peek() == "." && isDigit(scanner.peekNext()) {
		// Consumer the "."
		scanner.advance()
		for isDigit(scanner.peek()) {
			scanner.advance()
		}
	}
	doubleNumber, _ := strconv.ParseFloat(scanner.Source[scanner.Start:scanner.Current], 64)
	scanner.addToken(NUMBER, doubleNumber)
}

// identifier scans for identifiers and keywords doing something called "maximal munch"
func (scanner *TokenScanner) identifier() {
	for isAlphaNumeric(scanner.peek()) {
		scanner.advance()
	}
	text := scanner.Source[scanner.Start:scanner.Current]
	tokenType := keywords[text]
	if tokenType == 0 {
		tokenType = IDENTIFIER
	}
	scanner.AddToken(tokenType)
}

// blockComment - This is for scanning the block comment
func (scanner *TokenScanner) blockComment() {}

// peekNext returns the character after the current position in the source string without advancing the scanner.
// If the next position is out of bounds, it returns the null character ("\x00").
func (scanner *TokenScanner) peekNext() string {
	if scanner.Current+1 >= len(scanner.Source) {
		return "\x00"
	}
	return string(scanner.Source[scanner.Current+1])
}

// if there is no match, we return false. Otherwise, we keep scanning and return true
func (scanner *TokenScanner) match(expected string) bool {
	if scanner.isAtEnd() {
		return false
	}
	if string(scanner.Source[scanner.Current]) != expected {
		return false
	}
	scanner.Current++
	return true
}

// peek Peeks at the end of the string and returns it
// if it is at the end of the line, we return "\0" line end delimiter
// otherwise we return the current character
func (scanner *TokenScanner) peek() string {
	if scanner.isAtEnd() {
		return "\x00"
	}
	return string(scanner.Source[scanner.Current])
}

// AddToken Wrapper around addToken() in case we just need the TokenType added
func (scanner *TokenScanner) AddToken(tokenType TokenType) {
	scanner.addToken(tokenType, nil)
}

// addToken Scans the source and appends the tokens to the token array
func (scanner *TokenScanner) addToken(tokenType TokenType, literal interface{}) {
	text := scanner.Source[scanner.Start:scanner.Current]
	scanner.Tokens = append(scanner.Tokens, Token{
		Type:    tokenType,
		Lexeme:  text,
		Literal: literal,
		Line:    scanner.Line,
	})

}

// advance() - Scans and advances
func (scanner *TokenScanner) advance() string {
	char := scanner.Source[scanner.Current]
	scanner.Current += 1
	return string(char)
}
func (scanner *TokenScanner) isAtEnd() bool {
	return scanner.Current >= len(scanner.Source)
}

// Helper functions and stuff over here
func isDigit(c string) bool {
	return c >= "0" && c <= "9"
}
func isAlpha(c string) bool {
	return (c >= "a" && c <= "z") || (c >= "A" && c <= "Z") || (c == "_")
}
func isAlphaNumeric(c string) bool {
	return isAlpha(c) || isDigit(c)
}
