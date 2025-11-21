package interpreter

type ControlSignalType int

const (
	BREAK ControlSignalType = iota
	CONTINUE
	RETURN
)

type ControlSignal struct {
	Type  ControlSignalType
	Value any
}
