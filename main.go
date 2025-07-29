package main

import (
	"github.com/go-interpreter/internal/interpreter"
	"github.com/go-interpreter/internal/parser"
	"github.com/go-interpreter/internal/scanner"
)

func main() {
	//l := repl.Repl{}
	//l.LoadProgram("./examples/program.txt"
	programText := "var x = (1 + 2 + (3+4));print x;"
	tokenScanner := scanner.NewTokenScanner(programText)
	_ = tokenScanner.ScanTokens()
	p := parser.Parser{Tokens: tokenScanner.Tokens}
	expr := p.Parse()
	inter := interpreter.NewInterpreter()
	inter.Interpret(expr)
	//printer := printer.PrintAST{}
	//fmt.Println("\n", printer.Print(expr))
}
