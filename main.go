package main

import (
	"Crafting-interpreters/internal/interpreter"
	"Crafting-interpreters/internal/parser"
	"Crafting-interpreters/internal/scanner"
	"fmt"
)

func main() {
	//l := repl.Repl{}
	//l.LoadProgram("./examples/program.txt")
	programText := "(10+(10*True))"
	tokenScanner := scanner.TokenScanner{}
	tokenScanner.Init(programText)
	_ = tokenScanner.ScanTokens()
	p := parser.Parser{Tokens: tokenScanner.Tokens}
	expr, err := p.Parse()
	if err != nil {
		fmt.Println(err)
	}
	interpreter := interpreter.Interpreter{}
	_ = interpreter.Interpret(expr)
}
