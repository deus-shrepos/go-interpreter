package main

import (
	"os"

	"github.com/go-interpreter/internal/repl"
)

func main() {
	r := repl.NewRepl()
	if len(os.Args) > 1 {
		programPath := os.Args[1]
		r.LoadProgram(programPath)
	} else {
		r.LoadProgram("examples/program.txt")
	}

}
