package logentry

import (
	"fmt"
	"time"

	"github.com/danielchatfield/go-chalk"
	"github.com/danielchatfield/go-indicator"
)

// Status represents the status of the command
type Status int

// The supported statuses
const (
	RUNNING Status = iota
	SUCCESS
	FAILURE
)

// Command is a log entry that represents the running of a command
type Command struct {
	i       *indicator.Indicator
	cmd     string
	msg     string
	status  Status
	done    chan struct{} // indicates that the command has finished
	refresh chan struct{} // indicates that the command needs to be reprinted
}

//NewCommand returns a new command
func NewCommand(cmd string, msg string) *Command {
	c := &Command{
		i:    indicator.New(),
		cmd:  cmd,
		msg:  msg,
		done: make(chan struct{}),
	}

	go func() {
		for {
			select {
			case <-c.done:
				return
			case <-time.Tick(100 * time.Millisecond):
				c.i.Next()
				c.Refresh()
			}
		}
	}()

	return c
}

// Refresh sends a refresh signal if a subscriber is attached
func (c *Command) Refresh() {
	if c.refresh != nil {
		c.refresh <- struct{}{}
	}
}

// Subscribe sets the refresh event channel so that the console knows when
// the command needs reprinting
func (c *Command) Subscribe(refresh chan struct{}) chan struct{} {
	c.refresh = refresh
	return c.done
}

// Unsubscribe sets the refresh channel to nil to prevent any firther updates
// from being sent
func (c *Command) Unsubscribe() {
	c.refresh = nil
	close(c.done)
}

func (c *Command) String() string {
	return fmt.Sprintf("   %s %s", c.cmd, c.msg)
}

// InteractiveString returns a string with the Indicator
func (c *Command) InteractiveString() string {
	var cmd string

	switch c.status {
	case SUCCESS:
		cmd = chalk.Green(c.cmd)
	case FAILURE:
		cmd = chalk.Red(c.cmd)
	default:
		cmd = chalk.Cyan(c.cmd)
	}
	return fmt.Sprintf(" %s %s  %s", c.i, cmd, c.msg)
}

func (c *Command) terminate() {
	c.Refresh()
	c.done <- struct{}{}
	close(c.done)
	<-time.Tick(50 * time.Millisecond)
}

// Success changes the indicator to the success indicator then
func (c *Command) Success() {
	c.status = SUCCESS
	c.i.Success()
	c.terminate()
}

// Failure changes the indicator to the success indicator then
func (c *Command) Failure() {
	c.status = FAILURE
	c.i.Failure()
	c.terminate()
}
