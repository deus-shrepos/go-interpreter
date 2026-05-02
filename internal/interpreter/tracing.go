package interpreter

import (
	"fmt"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

var depth atomic.Int32

type Tracer struct {
	Visited      strings.Builder
	TraceEnabled bool
}

func NewTracer() Tracer {
	return Tracer{}
}

func (t Tracer) Trace() func() {
	if !t.TraceEnabled {
		return func() {}
	}
	// get the current function in the stack
	programCounter, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(programCounter)
	name := getFunctionNameOnly(fn.Name())
	d := depth.Add(1) - 1
	fmt.Printf("%sEntering-> %s\n", strings.Repeat("   ", int(d)), name)
	// t.Visited.WriteString(fmt.Sprintf("Entering -> %s\n", fn.Name()))
	start := time.Now()

	// defer this function where it needs to be traced
	return func() {
		fmt.Printf("%s<- Exiting %s (took %s)\n", strings.Repeat("   ", int(d)), name, time.Since(start))
		// t.Visited.WriteString(fmt.Sprintf("%s <- Exiting (took %s)\n", fn.Name(), time.Since(start)))
		depth.Add(-1)
	}
}

func getFunctionNameOnly(s string) string {
	if i := strings.LastIndex(s, "/"); i >= 0 {
		s = s[i+1:]
	}

	if i := strings.Index(s, "."); i >= 0 {
		s = s[i+1:]
	}

	s = strings.ReplaceAll(s, "(*", "")
	s = strings.ReplaceAll(s, ")", "")
	return s
}
