package repl

import (
	"fmt"
	"os"
	"strings"

	// "github.com/go-interpreter/internal/interpreter"
	// "github.com/go-interpreter/internal/interpreter"
	"github.com/chzyer/readline"
	"github.com/go-interpreter/internal/printer"
	"github.com/go-interpreter/internal/scanner"

	parser "github.com/go-interpreter/internal/parser"
)

// TODO(ME): NEED TO MAKE BETTER. COMING SOON.
type Repl struct {
	HadError bool
}

func NewRepl() *Repl {
	return &Repl{}

}

func (repl *Repl) LoadProgram(path string) {
	if repl.HadError {
		panic("Program had errors!")
	}
	err := repl.loadProgram(path)
	if err != nil {
		panic(err)
	}
}

// Runfile We want to scan the tokens in a file
// We want to scan correct tokens defined
// in out hypothetical language
func (repl *Repl) loadProgram(path string) error {

	// We are reading the program text here
	file, err := os.ReadFile(path)
	if err != nil {
		_ = fmt.Errorf("an error occured during the program file read: %s", err)
		return err
	}
	// We store that byte file for scanning
	tokenScanner := scanner.NewTokenScanner(string(file))
	repl.run(&tokenScanner)
	return nil
}

// This your token scanner for the program
func (repl *Repl) run(tokenScanner *scanner.TokenScanner) {
	_ = tokenScanner.ScanTokens()
	p := parser.NewParser(tokenScanner.Tokens)
	parsedStatments := p.Parse()
	printer := printer.PrintAST{}
	printer.PrintAll(parsedStatments)
	// inter := interpreter.NewInterpreter()
	//err := inter.Interpret(parsedStatments)
	//if err != nil {
	//	repl.HadError = true
	//	fmt.Println(err)
	//}
	//if repl.HadError {
	//	return
	//}
	// astPrinter := printer.PrintAST{}
	// astPrinter.Print(expr)
}

func (repl *Repl) RunCli() {
	rl, err := readline.NewEx(
		&readline.Config{
			Prompt:                 ">>> ",
			HistoryFile:            "/tmp/readline-multiline",
			DisableAutoSaveHistory: true,
		})
	if err != nil {
		panic(err)
	}

	defer rl.Close()

	var cmds []string
	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		cmds = append(cmds, line)
		if !strings.HasSuffix(line, "{") {
			rl.SetPrompt("... ")
			continue
		}
		cmd := strings.Join(cmds, " ")
		cmds = cmds[:0]
		rl.SetPrompt("> ")
		rl.SaveHistory(cmd)

	}
}
