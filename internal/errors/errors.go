package errors

import (
	"fmt"
)

// TODO: I don't want to call this InterpreterError. If I give my TW interpreter, I will name it that. Let me get back to it later on.
type ExecutionErrorType string

const (
	RUNTIME_ERROR ExecutionErrorType = "Runtime Error"
	PROGRAM_ERROR ExecutionErrorType = "Program Error"
	PARSER_ERROR  ExecutionErrorType = "Parser Error"
	SCANNER_ERROR ExecutionErrorType = "Scanner Error"
)

func (s ExecutionErrorType) String() string {
	return string(s)
}

type ExecutionError struct {
	Type    ExecutionErrorType
	Line    int
	Where   string
	Message string
}

// Report to user where and why that thing went wrong
func (err ExecutionError) Error() string {
	return fmt.Sprintf("%s [line %d] %s: %s", err.Type, err.Line, err.Where, err.Message)
}
