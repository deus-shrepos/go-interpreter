package interpreter

import "time"

type ClockFunction struct{}

func (f ClockFunction) call(i *Interpreter, args []any) (any, error) {
	return float64(time.Now().Local().UnixMilli()), nil

}

func (f ClockFunction) arity() int {
	return 0
}
