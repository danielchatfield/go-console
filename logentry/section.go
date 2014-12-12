package logentry

import "github.com/danielchatfield/go-chalk"

// Section represents a secion header
type Section struct {
	msg string
}

// NewSection returns a new Section with the specified message
func NewSection(msg string) *Section {
	return &Section{msg: msg}
}

func (s *Section) String() string {
	return chalk.Format(chalk.INTENSE+chalk.FOREGROUND, "\n\n   %s", s.msg)
}
