package main

import (
	"os"

	"github.com/go-interpreter/internal/repl"
)

func main() {
	r := repl.NewRepl()
	if len(os.Args) > 1 {

		if os.Args[1] == "interactive" {
			r.RunCli()
		}
		programPath := os.Args[1]
		r.LoadProgram(programPath)
	} else {
		r.LoadProgram("examples/program.txt")
	}
}

// go run main.go path/program.txt
//
