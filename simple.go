package console

import (
	"fmt"
	"io"
)

// A SimpleConsole represents a Logger implementation for simple shell sessions
type SimpleConsole struct {
	out io.Writer
}

// NewSimpleConsole returns a new SimpleConsole with the given writer
func NewSimpleConsole(w io.Writer) *SimpleConsole {
	return &SimpleConsole{
		out: w,
	}
}

// Log takes something that behaves like a string and logs it
func (l *SimpleConsole) Log(msg interface{}) (n int, err error) {
	return fmt.Fprintln(l.out, msg)
}
