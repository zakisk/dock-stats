package logger

import (
	"github.com/pterm/pterm"
)

type Message struct {
	Message   string
	KeyValues map[string]any
}

func (m *Message) Args() []pterm.LoggerArgument {
	return pterm.DefaultLogger.ArgsFromMap(m.KeyValues)
}
