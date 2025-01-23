package src

import (
	"fmt"
	"os"
)

type Lox struct {
	HadError bool
}

func (lox *Lox) LoadProgram(path string) {
	if lox.HadError {
		panic("Program had errors!")
	}
	err := lox.runfile(path)
	if err != nil {
		panic(err)
	}
}

// Runfile We want to scan the tokens in a file
// We want to scan correct tokens defined
// in out hypothetical langauge
func (lox *Lox) runfile(path string) error {

	// We are reading the program text here
	file, err := os.ReadFile(path)
	if err != nil {
		_ = fmt.Errorf("an error occured during the program file read: %s", err)
		return err
	}
	// We store that byte file for scanning
	tokenScanner := TokenScanner{Source: string(file)}
	lox.run(tokenScanner)
	return nil
}

// This your token scanner for the program
func (lox *Lox) run(tokenScanner TokenScanner) {
	//tokens := tokenScanner.ScanTokens()
	//parser := Parser{Tokens: tokens}
	//expr, err := parser.Parse()
	//if src.HadError {
	//	return
	//}
	//astPrinter := tool.PrintAST{}
	fmt.Println()
}

// ProgramError Signal to user that something went wrong
func (lox *Lox) ProgramError(line int, message string) {
	lox.report(line, "", message)
}

// Report to user where and why that thing went wrong
func (lox *Lox) report(line int, where string, message string) {
	fmt.Printf("[line %d] Error %s: %s", line, where, message)
	lox.HadError = true
}
