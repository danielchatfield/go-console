package console

import (
	"fmt"
	"io"
)

// InteractiveLogEntry is the interface that wraps the methods required for an
// interactive log entry (e.g. as used in an interactive shell)
//
// InteractiveString returns the interactive string e.g. with spinner on front
// Subscribe takes the channel where change signals are sent and returns a
// channel to send a done signal
// Unsubscribe indicates that events are no longer being listened for
type InteractiveLogEntry interface {
	InteractiveString() string
	Subscribe(refresh chan struct{}) chan struct{}
	Unsubscribe()
}

// An InteractiveConsole represents a Logger implementation for interactive
// shell sessions.
type InteractiveConsole struct {
	out  io.Writer
	done <-chan struct{}
}

// NewInteractiveConsole returns an InteractiveConsole with the given writer
func NewInteractiveConsole(w io.Writer) *InteractiveConsole {
	return &InteractiveConsole{
		out: w,
	}
}

// Log takes something that behaves like a string and logs it
func (l *InteractiveConsole) Log(msg interface{}) (n int, err error) {
	// An InteractiveLogEntry works by printing the message then when an update
	// occurs (e.g. spinner rotates or message changes) then it wipes the
	// previous output from the terminal and reprints. If another log entry
	// gets printed whilst it is active then it will delete the wrong line so
	// we should make sure the previous log entry is finished first.
	if l.done != nil {
		select {
		case <-l.done:
			// it is done
		default: // TODO: handle this properly
			panic("Cannot use console.Log() when a log entry is active")
		}
	}

	// lets try an interface upgrade
	if imsg, ok := msg.(InteractiveLogEntry); ok {
		// subscribe to channel
		refresh := make(chan struct{})
		done := imsg.Subscribe(refresh)
		l.done = done

		// setup goroutine to "listen" for changes
		go func() {
			for {
				select {
				case <-refresh:
					fmt.Fprint(l.out, "\r")
					fmt.Fprint(l.out, imsg.InteractiveString())
				case <-done:
					fmt.Fprint(l.out, "\r")
					fmt.Fprintln(l.out, imsg.InteractiveString())
					return
				}
			}
		}()
		return

	}
	return fmt.Fprintln(l.out, msg)
}

func (l *InteractiveConsole) closeChannel() {

}
