package main

import (
	"github.com/go-interpreter/internal/parser"
	"github.com/go-interpreter/internal/scanner"
)

func main() {
	//l := repl.Repl{}
	//l.LoadProgram("./examples/program.txt"
	programText := "var x = 10; x = 20;"
	tokenScanner := scanner.NewTokenScanner(programText)
	_ = tokenScanner.ScanTokens()
	p := parser.Parser{Tokens: tokenScanner.Tokens}
	_ = p.Parse()
	//inter := interpreter.NewInterpreter()
	//inter.Interpret(expr)
	//printer := printer.PrintAST{}
	//fmt.Println("\n", printer.Print(expr))
}
