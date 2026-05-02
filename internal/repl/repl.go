package repl

import (
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/go-interpreter/internal/interpreter"
	"github.com/go-interpreter/internal/printer"
	"github.com/go-interpreter/internal/scanner"
	"github.com/go-interpreter/internal/utils"

	parser "github.com/go-interpreter/internal/parser"
)

// TODO(ME): LITERALLY CREATE A NICE REPL :)
type Repl struct {
	HadError bool
}

func NewRepl() *Repl {
	return &Repl{}

}

// Runfile We want to scan the tokens in a file
// We want to scan correct tokens defined
// in out hypothetical language
func (repl *Repl) LoadProgramFromPath(path string, astOnly bool, traceEnabled bool) {
	if repl.HadError {
		panic("Program execution stopped!")

	}

	// We are reading the program text here
	file, err := os.ReadFile(path)
	if err != nil {
		_ = fmt.Errorf("an error occured during the program file read: %s", err)
	}
	// We store that byte file for scanning
	tokenScanner := scanner.NewTokenScanner(string(file))
	if astOnly {
		repl.PrintAST(&tokenScanner)
	}
	repl.run(&tokenScanner, traceEnabled)
}

func (repl *Repl) LoadFromString(program string, astOnly bool, traceEnabled bool) {
	tokenScanner := scanner.NewTokenScanner(program)
	if astOnly {
		repl.PrintAST(&tokenScanner)
		return
	}
	repl.run(&tokenScanner, traceEnabled)
}

// This your token scanner for the program
func (repl *Repl) run(tokenScanner *scanner.TokenScanner, traceEnabled bool) {
	_ = tokenScanner.ScanTokens()
	p := parser.NewParser(tokenScanner.Tokens)
	parsedStatments := p.Parse()
	outputStream := utils.NewOutStream(os.Stdout, os.Stdin)
	inter := interpreter.NewInterpreter(outputStream, traceEnabled)
	err := inter.Interpret(parsedStatments)
	if err != nil {
		repl.HadError = true
	}
	if repl.HadError {
		return
	}
}

func (repl *Repl) PrintAST(tokenScanner *scanner.TokenScanner) {
	_ = tokenScanner.ScanTokens()
	p := parser.NewParser(tokenScanner.Tokens)
	parsedStatments := p.Parse()
	print := printer.PrintAST{}
	print.PrintAll(parsedStatments)
}

// TODO: Fix the cli later when we have all the basic langauge constructs
func (repl *Repl) RunCli() {
	rl, err := readline.NewEx(
		&readline.Config{
			Prompt:                 ">>> ",
			HistoryFile:            "/tmp/readline-multiline",
			DisableAutoSaveHistory: true,
			EOFPrompt:              "exit",
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
