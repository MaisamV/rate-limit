package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// LogEvent represents a log event that can be chained with additional context
type LogEvent interface {
	Str(key, val string) LogEvent
	Int(key string, i int) LogEvent
	Int64(key string, i int64) LogEvent
	Bool(key string, b bool) LogEvent
	Dur(key string, d time.Duration) LogEvent
	Err(err error) LogEvent
	Msg(msg string)
}

// Logger defines the interface for logging operations
type Logger interface {
	Debug() LogEvent
	Info() LogEvent
	Warn() LogEvent
	Error() LogEvent
	Fatal() LogEvent
}

// zerologEvent wraps zerolog.Event to implement LogEvent interface
type zerologEvent struct {
	event *zerolog.Event
}

func (e *zerologEvent) Str(key, val string) LogEvent {
	return &zerologEvent{event: e.event.Str(key, val)}
}

func (e *zerologEvent) Int(key string, i int) LogEvent {
	return &zerologEvent{event: e.event.Int(key, i)}
}

func (e *zerologEvent) Int64(key string, i int64) LogEvent {
	return &zerologEvent{event: e.event.Int64(key, i)}
}

func (e *zerologEvent) Bool(key string, b bool) LogEvent {
	return &zerologEvent{event: e.event.Bool(key, b)}
}

func (e *zerologEvent) Dur(key string, d time.Duration) LogEvent {
	return &zerologEvent{event: e.event.Dur(key, d)}
}

func (e *zerologEvent) Err(err error) LogEvent {
	return &zerologEvent{event: e.event.Err(err)}
}

func (e *zerologEvent) Msg(msg string) {
	e.event.Msg(msg)
}

// zerologLogger wraps zerolog.Logger to implement Logger interface
type zerologLogger struct {
	logger zerolog.Logger
}

func (l *zerologLogger) Debug() LogEvent {
	return &zerologEvent{event: l.logger.Debug()}
}

func (l *zerologLogger) Info() LogEvent {
	return &zerologEvent{event: l.logger.Info()}
}

func (l *zerologLogger) Warn() LogEvent {
	return &zerologEvent{event: l.logger.Warn()}
}

func (l *zerologLogger) Error() LogEvent {
	return &zerologEvent{event: l.logger.Error()}
}

func (l *zerologLogger) Fatal() LogEvent {
	return &zerologEvent{event: l.logger.Fatal()}
}

// New creates a new logger instance
func New() Logger {
	// Configure zerolog
	zerolog.TimeFieldFormat = time.RFC3339

	// Create logger with console writer for development
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()

	return &zerologLogger{logger: logger}
}

// NewWithLevel creates a new logger with specified level
func NewWithLevel(level string) Logger {
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(logLevel)
	return New()
}
