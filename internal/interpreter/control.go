package interpreter

type ControlSig int

const (
	BREAK ControlSig = iota
	CONTINUE
)
