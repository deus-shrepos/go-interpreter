package main

import (
	// "fmt"

	// "fmt"

	"github.com/go-interpreter/internal/repl"
)

func main() {
	r := repl.NewRepl()
	r.LoadProgram("examples/while_stmts.txt")
}
