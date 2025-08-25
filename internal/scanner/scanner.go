package scanner

import (
	"fmt"
	"strconv"

	"github.com/go-interpreter/internal/errors"
	"github.com/go-interpreter/internal/token"
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
	Tokens  []token.Token
	Start   int
	Current int
	Line    int
}

// NewTokenScanner Init This initializes the source code
func NewTokenScanner(source string) TokenScanner {
	return TokenScanner{
		Source:  source,
		Start:   0,
		Current: 0,
		Line:    0,
	}
}

// ScanTokens scans the source code and produces a list of tokens based on the language grammar.
func (scanner *TokenScanner) ScanTokens() []token.Token {
	for !scanner.isAtEnd() {
		scanner.Start = scanner.Current
		err := scanner.ScanToken()
		if err != nil {
			fmt.Println(err)
		}
	}
	scanner.Tokens = append(scanner.Tokens, token.Token{Type: token.EOF, Lexeme: "", Literal: nil, Line: scanner.Line})
	return scanner.Tokens
}

// ScanToken reads the next character in the source and determines the appropriate token type to add to the token list.
func (scanner *TokenScanner) ScanToken() error {
	c := scanner.advance()

	switch c {
	case "(":
		scanner.AddToken(token.LEFT_PAREN)
	case ")":
		scanner.AddToken(token.RIGHT_PAREN)
	case "{":
		scanner.AddToken(token.LEFT_BRACE)
	case "}":
		scanner.AddToken(token.RIGHT_BRACE)
	case ",":
		scanner.AddToken(token.COMMA)
	case ".":
		scanner.AddToken(token.DOT)
	case "-":
		scanner.AddToken(token.MINUS)
	case "+":
		scanner.AddToken(token.PLUS)
	case ";":
		scanner.AddToken(token.SEMICOLON)
	case "*":
		scanner.AddToken(token.STAR)
	case "!":
		bang := token.BANG
		if scanner.match("=") {
			bang = token.BANG_EQUAL
		}
		scanner.AddToken(bang)
	case "=":
		equal := token.EQUAL
		if scanner.match("=") {
			equal = token.EQUAL_EQUAL
		}
		scanner.AddToken(equal)
	case "<":
		lessEqual := token.LESS
		if scanner.match("=") {
			lessEqual = token.LESS_EQUAL
		}
		scanner.AddToken(lessEqual)
	case ">":
		greaterEqual := token.GREATER
		if scanner.match("=") {
			greaterEqual = token.GREATER_EQUAL
		}
		scanner.AddToken(greaterEqual)
	case "/":
		if scanner.match("/") {
			// This is a comment - keep going until you reach at the end of the line
			for scanner.peek() != "\n" && !scanner.isAtEnd() {
				scanner.advance()
			}
		} else {
			scanner.AddToken(token.SLASH)
		}
	case " ":
	case "\r":
	case "\t":
		break
	case "\n":
		scanner.Line++
	case "\"":
		err := scanner.string()
		if err != nil {
			return err
		}
	default:
		if isDigit(c) {
			scanner.number()
		} else if isAlpha(c) {
			scanner.identifier()
		} else {
			return errors.ExecutionError{Type: errors.SCANNER_ERROR,
				Line:    scanner.Line,
				Where:   scanner.Current,
				Message: "Unexpected Error"}
		}
	}
	return nil
}

// string scans a string literal, handles line tracking for multi-line strings, and adds the token to the token list.
func (scanner *TokenScanner) string() error {
	// Two cases run here:
	//	1) continue scanning until you find " and you are not at the EOL
	//  2) did not find " but you are the EOL
	for scanner.peek() != "\"" && !scanner.isAtEnd() {
		if scanner.peek() == "\n" {
			scanner.Line++
		}
		scanner.advance()
	}
	// If we are the EOL, we have an unterminated string
	if scanner.isAtEnd() {
		return errors.ExecutionError{Type: errors.SCANNER_ERROR,
			Line:    scanner.Line,
			Where:   scanner.Current,
			Message: "Unterminated string",
		}
	}
	scanner.advance()
	value := scanner.Source[scanner.Start+1 : scanner.Current-1]
	value, err := strconv.Unquote(`"` + value + `"`) // From raw to an actual string
	if err != nil {
		return err
	}
	scanner.addToken(token.STRING, value)
	return nil
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
	scanner.addToken(token.NUMBER, doubleNumber)
}

// identifier scans for identifiers and keywords doing something called "maximal munch"
func (scanner *TokenScanner) identifier() {
	for isAlphaNumeric(scanner.peek()) {
		scanner.advance()
	}
	text := scanner.Source[scanner.Start:scanner.Current]
	tokenType := token.Keywords[text]
	if tokenType == 0 {
		tokenType = token.IDENTIFIER
	}
	scanner.AddToken(tokenType)
}

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
	return scanner.Source[scanner.Current : scanner.Current+1]
}

// AddToken Wrapper around addToken() in case we just need the TokenType added
func (scanner *TokenScanner) AddToken(tokenType token.TokenType) {
	scanner.addToken(tokenType, nil)
}

// addToken Scans the source and appends the tokens to the token array
func (scanner *TokenScanner) addToken(tokenType token.TokenType, literal any) {
	text := scanner.Source[scanner.Start:scanner.Current]
	scanner.Tokens = append(scanner.Tokens, token.Token{
		Type:    tokenType,
		Lexeme:  text,
		Literal: literal,
		Line:    scanner.Line,
		Char:    scanner.Start,
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
