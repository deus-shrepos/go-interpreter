package main

import (
	"crafting-interpreters/lox"
	"crafting-interpreters/tool"
	"fmt"
)

var type_ = map[string]lox.TokenType{
	"and": lox.AND,
	"or":  lox.OR,
}

func main() {
	//lox := lox{}
	//lox.loadprogram("program.txt")
	//programText := `//This is a comment`
	//scanner := TokenScanner{}
	//scanner.Init(programText)
	//_ = scanner.ScanTokens()
	//fmt.Println(scanner.Tokens)
	//v1 := Literal{Value: 10}
	//v2 := Literal{Value: 20}
	//bin1 := Binary{Left: v1, Operator: "+", Right: v2}
	//fmt.Println(bin1)
	err := tool.GenerateAst("lox/")
	if err != nil {
		panic(err)
	}
	MinusBinOp := lox.Binary{
		Left: lox.Literal{Value: 10},
		Operator: lox.Token{
			Type:    lox.MINUS,
			Lexeme:  "-",
			Literal: "-",
			Line:    1,
		},
		Right: lox.Literal{Value: 20},
	}
	AddBinOp := lox.Binary{
		Left: MinusBinOp,
		Operator: lox.Token{
			Type:    lox.PLUS,
			Lexeme:  "+",
			Literal: "+",
			Line:    1,
		},
		Right: lox.Literal{Value: 20},
	}
	printAst := tool.PrintAST{}
	str := AddBinOp.Accept(printAst)
	fmt.Println(str)
}
