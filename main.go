package main

import (
	"crafting-interpreters/src"
	"crafting-interpreters/tool"
	"fmt"
)

var type_ = map[string]src.TokenType{
	"and": src.AND,
	"or":  src.OR,
}

func main() {
	//src := src{}
	//src.loadprogram("program.txt")
	//programText := `//This is a comment`
	//scanner := TokenScanner{}
	//scanner.Init(programText)
	//_ = scanner.ScanTokens()
	//fmt.Println(scanner.Tokens)
	//v1 := Literal{Value: 10}
	//v2 := Literal{Value: 20}
	//bin1 := Binary{Left: v1, Operator: "+", Right: v2}
	//fmt.Println(bin1)
	err := tool.GenerateAst("src/")
	if err != nil {
		panic(err)
	}
	MinusBinOp := src.Binary{
		Left: src.Literal{Value: 10},
		Operator: src.Token{
			Type:    src.MINUS,
			Lexeme:  "-",
			Literal: "-",
			Line:    1,
		},
		Right: src.Literal{Value: 20},
	}
	AddBinOp := src.Binary{
		Left: MinusBinOp,
		Operator: src.Token{
			Type:    src.PLUS,
			Lexeme:  "+",
			Literal: "+",
			Line:    1,
		},
		Right: src.Literal{Value: 20},
	}
	printAst := tool.PrintAST{}
	str := AddBinOp.Accept(printAst)
	fmt.Println(str)
}
