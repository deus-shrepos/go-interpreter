package interpreter

type Callable interface {
	call(interpreter Interpreter, arguments any)
}
