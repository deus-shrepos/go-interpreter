package main

import (
	"Crafting-interpreters/internal/interpreter"
	"Crafting-interpreters/internal/parser"
	"Crafting-interpreters/internal/printer"
	"Crafting-interpreters/internal/scanner"
	"fmt"
)

func main() {
	//l := repl.Repl{}
	//l.LoadProgram("./examples/program.txt"
	programText := "(10+(10 * 2 + (3 * 3)))"
	tokenScanner := scanner.TokenScanner{}
	tokenScanner.Init(programText)
	_ = tokenScanner.ScanTokens()
	p := parser.Parser{Tokens: tokenScanner.Tokens}
	expr, _ := p.Parse()
	interpreter := interpreter.Interpreter{}
	interpreter.Interpret(expr, true)
	printer := printer.PrintAST{}
	fmt.Println("\n", printer.Print(expr))
}
