package main

import (
	"Crafting-interpreters/internal/parser"
	"Crafting-interpreters/internal/printer"
	"Crafting-interpreters/internal/repl"
	"Crafting-interpreters/internal/scanner"
	"fmt"
)

func main() {
	// l := internal.Repl{}
	// l.LoadProgram("./examples/program.txt")
	programText := "(10*(4*5))"
	scanner := scanner.TokenScanner{}
	scanner.Init(programText)
	_ = scanner.ScanTokens()
	p := parser.Parser{Tokens: scanner.Tokens, Lox: repl.Repl{}}
	expr, _ := p.Parse()
	astPrinter := printer.PrintAST{}
	astString := astPrinter.Print(expr)
	fmt.Printf(astString)

	//fmt.Println(scanner.Tokens)
	//v1 := Literal{Value: 10}
	//v2 := Literal{Value: 20}
	//bin1 := Binary{Left: v1, Operator: "+", Right: v2}
	//fmt.Println(bin1)
	// err := internal.GenerateAst("internal/")
	// if err != nil {
	// 	panic(err)
	// }
	// MinusBinOp := internal.Binary{
	// 	Left: internal.Literal{Value: 10},
	// 	Operator: internal.Token{
	// 		Type:    internal.MINUS,
	// 		Lexeme:  "-",
	// 		Literal: "-",
	// 		Line:    1,
	// 	},
	// 	Right: internal.Literal{Value: 20},
	// }
	// AddBinOp := internal.Binary{
	// 	Left: MinusBinOp,
	// 	Operator: internal.Token{
	// 		Type:    internal.PLUS,
	// 		Lexeme:  "+",
	// 		Literal: "+",
	// 		Line:    1,
	// 	},
	// 	Right: internal.Literal{Value: 20},
	// }
	// printAst := internal.PrintAST{}
	// str := AddBinOp.Accept(printAst)
	// fmt.Println(str)

}
