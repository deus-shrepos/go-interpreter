package errors

import (
	"Crafting-interpreters/internal/token"
	"fmt"
)

// TODO: I don't want to call this InterpreterError. If I give my TW interpreter, I will name it that. Let me get back to it later on.
type ExecutionErrorType string

const (
	RUNTIME_ERROR ExecutionErrorType = "Runtime Error"
	PROGRAM_ERROR ExecutionErrorType = "Program Error"
	PARSER_ERROR  ExecutionErrorType = "Parse Error"
)

func (s ExecutionErrorType) String() string {
	return string(s)
}

type ExecutionError struct {
	Type    ExecutionErrorType
	Op      token.Token
	Message string
}

func (err ExecutionError) Error() string {
	return report(err)
}

// Report to user where and why that thing went wrong
func report(error ExecutionError) string {
	return fmt.Sprintf("%s [line %d] %s: %s", error.Type, error.Op.Line, "", error.Message)
}
