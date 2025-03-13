package errors

import "fmt"

type Error struct{}

// ProgramError Signal to user that something went wrong
func (err *Error) ProgramError(line int, message string) {
	err.report(line, "", message)
}

// Report to user where and why that thing went wrong
func (err *Error) report(line int, where string, message string) {
	fmt.Printf("[line %d] Error %s: %s", line, where, message)
}
