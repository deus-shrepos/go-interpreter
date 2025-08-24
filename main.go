package main

import (
	"github.com/go-interpreter/internal/interpreter"
	"github.com/go-interpreter/internal/parser"
	"github.com/go-interpreter/internal/scanner"
)

func main() {
	// THE MAIN FILE WILL CHANGE
	programText := `
	var globalScope = 10;
	  {
	 	 var innerScope1 = 20;
	     {
			var innerScope2 = 30;
	  		print innerScope2;
			print "\n";
	     }
	     print innerScope1;
		 print "\n";
	  }
	print globalScope;
	`
	tokenScanner := scanner.NewTokenScanner(programText)
	_ = tokenScanner.ScanTokens()
	p := parser.Parser{Tokens: tokenScanner.Tokens}
	expr := p.Parse()
	inter := interpreter.NewInterpreter()
	err := inter.Interpret(expr)
	if err != nil {
		panic(err)
	}
	//printer := printer.PrintAST{}
	//fmt.Println("\n", printer.Print(expr))
}
