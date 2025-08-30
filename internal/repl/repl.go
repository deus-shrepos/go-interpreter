package repl

import (
	"fmt"
	"os"

	"github.com/go-interpreter/internal/interpreter"
	"github.com/go-interpreter/internal/scanner"

	parser "github.com/go-interpreter/internal/parser"
)

// TODO(ME): NEED TO MAKE BETTER. COMING SOON.
type Repl struct {
	HadError bool
}

func NewRepl() *Repl {
	return &Repl{}

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
	tokenScanner := scanner.NewTokenScanner(string(file))
	repl.run(&tokenScanner)
	return nil
}

// This your token scanner for the program
func (repl *Repl) run(tokenScanner *scanner.TokenScanner) {
	_ = tokenScanner.ScanTokens()
	p := parser.NewParser(tokenScanner.Tokens)
	inter := interpreter.NewInterpreter()
	err := inter.Interpret(p.Parse())
	if err != nil {
		repl.HadError = true
		fmt.Println(err)
	}
	if repl.HadError {
		return
	}
	// astPrinter := printer.PrintAST{}
	// astPrinter.Print(expr)
}
