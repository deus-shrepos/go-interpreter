package interpreter

import "github.com/go-interpreter/internal/ast"

// Callable represents any entity that can be invoked with arguments,
// such as functions or methods. It defines the contract for calling
// such entities and retrieving their arity (number of expected arguments).
type Callable interface {
	call(i *Interpreter, args []any) (any, error)
	arity() int
}

type Function struct {
	Declarations ast.Function
	Closure      *Environment
}

// NewFunction creates a new Function instance from the given function declaration.
// It initializes the Function with its declarations and the closed envirnoment.
// This allows the interpreter to manage function calls and their associated data.
func NewFunction(funcDeclaration ast.Function, closure *Environment) Function {
	return Function{
		Declarations: funcDeclaration,
		Closure:      closure,
	}
}

// call executes the function with the provided arguments in a new environment.
// It sets up the function's environment with access to the global scope and
// defines the function's parameters. After executing the function's body,
// it returns the result or an error if one occurs.
func (f Function) call(i *Interpreter, args []any) (any, error) {
	functionEnvirnoment := NewEnvironment(f.Closure) // lexcial scope surronding the function
	for idx, arg := range args {
		functionEnvirnoment.Define(f.Declarations.Parameters[idx].Lexeme, arg)
	}

	value, err := i.execBlock(f.Declarations.Body, functionEnvirnoment)
	if err != nil {
		return nil, err
	}

	if control, isReturn := value.(ControlSignal); isReturn {
		return control.Value, nil
	}
	return nil, nil
}

// arity returns the number of parameters the function expects.
func (f Function) arity() int {
	return len(f.Declarations.Parameters)
}
