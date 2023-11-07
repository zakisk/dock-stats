package logger

import (
	"github.com/pterm/pterm"
)

func NewLog() *Log {
	logger := pterm.DefaultLogger.WithKeyStyles(map[string]pterm.Style{
		KeyError:         *pterm.NewStyle(pterm.FgRed),
		KeyProbableCause: *pterm.NewStyle(pterm.FgLightMagenta),
		KeyRemedy:        *pterm.NewStyle(pterm.FgLightYellow),
	})
	return &Log{Logger: logger}
}

func (l *Log) Trace(message string, args []pterm.LoggerArgument) {
	l.Logger.Trace(message, args)
}

func (l *Log) Debug(message string, args []pterm.LoggerArgument) {
	l.Logger.Debug(message, args)
}

func (l *Log) Info(message string, args []pterm.LoggerArgument) {
	l.Logger.Info(message, args)
}

func (l *Log) Warn(message string, args []pterm.LoggerArgument) {
	l.Logger.Warn(message, args)
}

func (l *Log) Error(message string, args []pterm.LoggerArgument) {
	l.Logger.Error(message, args)
}

func (l *Log) Fatal(message string, args []pterm.LoggerArgument) {
	l.Logger.Fatal(message, args)
}

func (l *Log) Print(message string, args []pterm.LoggerArgument) {
	l.Logger.Print(message, args)
}

func (l *Log) Args(args ...any) []pterm.LoggerArgument {
	return l.Logger.Args(args)
}

type Log struct {
	Logger *pterm.Logger
}

type Logger interface {
	Print(message string, args ...pterm.LoggerArgument)
}