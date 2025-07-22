package interpreter

import (
	"fmt"

	"github.com/go-interpreter/internal/errors"
	"github.com/go-interpreter/internal/token"
)

// Environment Represents an Environment in a scope
// Used to store bindings
type Environment struct {
	Values map[string]any
}

// NewEnvirnoment Initiates a new Environment
func NewEnvirnoment() Environment {
	return Environment{
		Values: make(map[string]any),
	}
}

// Define Defines a variable in the envirnoment
// It will set as a mapping bound to a value of any type
func (env *Environment) Define(varName string, value any) {
	env.Values[varName] = value
}

// Get Gets the value of a bound variable in an envirnoment
// If it doesn't find i, it raises an execution Error
func (env *Environment) Get(token token.Token) (any, error) {
	_, exists := env.Values[token.Lexeme]
	if exists {
		return env.Values[token.Lexeme], nil
	}
	return nil, errors.ExecutionError{
		Type:    errors.RUNTIME_ERROR,
		Line:    token.Line,
		Where:   token.Char,
		Message: fmt.Sprintf("%s Undefined variable %s.", token, token.Lexeme),
	}
}
