package main

import (
	"fmt"

	"github.com/go-interpreter/internal/repl"

	"os"
)

func cliHelp() {
	fmt.Println("Valid Args:")
	fmt.Println("--path: load from a specific path")
	fmt.Println("--string: load from a string")
	fmt.Println("Both can print the AST tree with at the end --ast")

}

func main() {
	r := repl.NewRepl()
	if len(os.Args) > 1 {
		hasASTFlag := false

		if os.Args[len(os.Args)-1] == "--ast" {
			hasASTFlag = true

		}
		switch os.Args[1] {
		case "--interactive":
			r.RunCli()
		case "--path":
			r.LoadProgramFromPath(os.Args[2], hasASTFlag)
		case "--string":
			r.LoadFromString(os.Args[2], hasASTFlag)
		case "--help":
			cliHelp()
		default:
			fmt.Println("Invalid Flags")
			cliHelp()
		}
	}
}
