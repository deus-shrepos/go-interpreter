package repl

import (
	"fmt"
	"os"

	"github.com/go-interpreter/internal/printer"
	"github.com/go-interpreter/internal/scanner"

	parser "github.com/go-interpreter/internal/parser"
)

type Repl struct {
	HadError bool
}

func (repl *Repl) LoadProgram(path string) {
	if repl.HadError {
		panic("Program had errors!")
	}
	err := repl.loadProgram(path)
	if err != nil {
		panic(err)
	}
}

// Runfile We want to scan the tokens in a file
// We want to scan correct tokens defined
// in out hypothetical language
func (repl *Repl) loadProgram(path string) error {

	// We are reading the program text here
	file, err := os.ReadFile(path)
	if err != nil {
		_ = fmt.Errorf("an error occured during the program file read: %s", err)
		return err
	}
	// We store that byte file for scanning
	tokenScanner := scanner.TokenScanner{Source: string(file)}
	repl.run(tokenScanner)
	return nil
}

// This your token scanner for the program
func (repl *Repl) run(tokenScanner scanner.TokenScanner) {
	tokens := tokenScanner.ScanTokens()
	parser := parser.Parser{Tokens: tokens}
	expr, _ := parser.Parse()
	if repl.HadError {
		return
	}
	astPrinter := printer.PrintAST{}
	astPrinter.Print(expr)
}
