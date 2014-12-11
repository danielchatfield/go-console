package console

import (
	"fmt"

	"github.com/danielchatfield/go-indicator"
)

// CommandLogEntry is a log entry that represents the running of a command
type CommandLogEntry struct {
	i       *indicator.Indicator
	cmd     string
	msg     string
	done    chan struct{} // indicates that the command has finished
	refresh chan struct{} // indicates that the command needs to be reprinted
}

//NewCommandLogEntry returns a new commandLogEntry
func NewCommandLogEntry(cmd string, msg string) *CommandLogEntry {
	c := &CommandLogEntry{
		i:    indicator.New(),
		cmd:  cmd,
		msg:  msg,
		done: make(chan struct{}),
	}

	return c
}

// Subscribe sets the refresh event channel so that the console knows when
// the command needs reprinting
func (c *CommandLogEntry) Subscribe(refresh chan struct{}) chan struct{} {
	c.refresh = refresh
	return c.done
}

// Unsubscribe sets the refresh channel to nil to prevent any firther updates
// from being sent
func (c *CommandLogEntry) Unsubscribe() {
	c.refresh = nil
}

func (c *CommandLogEntry) String() string {
	return fmt.Sprintf("   %s  %s", c.cmd, c.msg)
}

// InteractiveString returns a string with the Indicator
func (c *CommandLogEntry) InteractiveString() string {
	return fmt.Sprintf(" %s  %s  %s", c.i, c.cmd, c.msg)
}

// Success changes the indicator to the success indicator then
func (c *CommandLogEntry) Success() {
	c.i.Success()
	close(c.done)
}
