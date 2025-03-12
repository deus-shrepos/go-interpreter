package main

import (
	"Crafting-interpreters/src"
	"fmt"
)

func main() {
	// l := src.Lox{}
	// l.LoadProgram("./examples/program.txt")
	programText := "(10*(4*5))"
	scanner := src.TokenScanner{}
	scanner.Init(programText)
	_ = scanner.ScanTokens()
	p := src.Parser{Tokens: scanner.Tokens, Lox: src.Lox{}}
	expr, _ := p.Parse()
	astPrinter := src.PrintAST{}
	astString := astPrinter.Print(expr)
	fmt.Printf(astString)

	//fmt.Println(scanner.Tokens)
	//v1 := Literal{Value: 10}
	//v2 := Literal{Value: 20}
	//bin1 := Binary{Left: v1, Operator: "+", Right: v2}
	//fmt.Println(bin1)
	// err := src.GenerateAst("src/")
	// if err != nil {
	// 	panic(err)
	// }
	// MinusBinOp := src.Binary{
	// 	Left: src.Literal{Value: 10},
	// 	Operator: src.Token{
	// 		Type:    src.MINUS,
	// 		Lexeme:  "-",
	// 		Literal: "-",
	// 		Line:    1,
	// 	},
	// 	Right: src.Literal{Value: 20},
	// }
	// AddBinOp := src.Binary{
	// 	Left: MinusBinOp,
	// 	Operator: src.Token{
	// 		Type:    src.PLUS,
	// 		Lexeme:  "+",
	// 		Literal: "+",
	// 		Line:    1,
	// 	},
	// 	Right: src.Literal{Value: 20},
	// }
	// printAst := src.PrintAST{}
	// str := AddBinOp.Accept(printAst)
	// fmt.Println(str)

}
