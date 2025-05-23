package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

type Logger struct {
	levelPriority int
	levels        map[string]int
	stdout        *log.Logger
	stderr        *log.Logger
}

type Event struct {
	logger *Logger
	level  string
	fields map[string]interface{}
}

func New(level string) *Logger {
	levelOrder := []string{LevelDebug, LevelInfo, LevelWarn, LevelError}
	levels := make(map[string]int, len(levelOrder))
	for i, lvl := range levelOrder {
		levels[lvl] = i
	}

	normalizedLevel := strings.ToLower(strings.TrimSpace(level))
	priority, exists := levels[normalizedLevel]
	if !exists {
		priority = levels[LevelInfo]
	}

	return &Logger{
		levelPriority: priority,
		levels:        levels,
		stdout:        log.New(os.Stdout, "", log.LstdFlags),
		stderr:        log.New(os.Stderr, "", log.LstdFlags),
	}
}

func (l *Logger) isLevelEnabled(level string) bool {
	return l.levelPriority <= l.levels[level]
}

func (l *Logger) Debug() *Event {
	return l.newEvent(LevelDebug)
}

func (l *Logger) Info() *Event {
	return l.newEvent(LevelInfo)
}

func (l *Logger) Warn() *Event {
	return l.newEvent(LevelWarn)
}

func (l *Logger) Error() *Event {
	return l.newEvent(LevelError)
}

func (l *Logger) newEvent(level string) *Event {
	if !l.isLevelEnabled(level) {
		return nil
	}
	return &Event{
		logger: l,
		level:  level,
		fields: make(map[string]interface{}),
	}
}

func (e *Event) Str(key, val string) *Event {
	if e != nil {
		e.fields[key] = val
	}
	return e
}

func (e *Event) Int(key string, val int) *Event {
	if e != nil {
		e.fields[key] = val
	}
	return e
}

func (e *Event) Err(err error) *Event {
	if e != nil && err != nil {
		e.fields[LevelError] = err.Error()
	}
	return e
}

func (e *Event) Msg(msg string) {
	if e == nil {
		return
	}

	out := fmt.Sprintf("[%s] %s", strings.ToUpper(e.level), msg)
	for k, v := range e.fields {
		out += fmt.Sprintf(" | %s=%v", k, v)
	}

	switch e.level {
	case LevelError:
		e.logger.stderr.Println(out)
	default:
		e.logger.stdout.Println(out)
	}
}

func (e *Event) Msgf(format string, args ...interface{}) {
	if e != nil {
		e.Msg(fmt.Sprintf(format, args...))
	}
}
