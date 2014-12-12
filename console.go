package console

import (
	"os"

	"github.com/danielchatfield/go-console/logentry"
)

// Status represents the status of the task
type Status int

// The Status constants
const (
	SUCCESS Status = iota
	WARNING
	FAILURE
)

var (
	console = NewConsole()
)

// Console is the interface that wraps the basic logging functions.
type Console interface {
	Log(arg interface{}) (n int, err error)
}

// NewConsole returns a new InteractiveConsole if the current terminal supports
// interactivity otherwise returns a SimpleConsole
func NewConsole() Console {
	return NewInteractiveConsole(os.Stdout)
}

// Log logs using the default console
func Log(msg interface{}) (n int, err error) {
	return console.Log(msg)
}

// LogCommand creates a new CommandLogEntry, logs it and returns it
func LogCommand(cmd string) *logentry.Command {
	c := logentry.NewCommand(cmd, "")
	console.Log(c)
	return c
}

// LogSection logs a section header
func LogSection(str string) {
	c := logentry.NewSection(str)
	console.Log(c)
}
