// Package audit provides a structured audit log for port scan events.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of an audit event.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

// Entry is a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     Level     `json:"level"`
	Event     string    `json:"event"`
	Detail    string    `json:"detail,omitempty"`
}

// Logger writes audit entries to a destination.
type Logger struct {
	w   io.Writer
	json bool
}

// New creates an audit Logger. If w is nil, os.Stdout is used.
func New(w io.Writer, jsonFormat bool) *Logger {
	if w == nil {
		w = os.Stdout
	}
	return &Logger{w: w, json: jsonFormat}
}

// Log writes an audit entry at the given level.
func (l *Logger) Log(level Level, event, detail string) error {
	e := Entry{
		Timestamp: time.Now().UTC(),
		Level:     level,
		Event:     event,
		Detail:    detail,
	}
	if l.json {
		b, err := json.Marshal(e)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(l.w, "%s\n", b)
		return err
	}
	_, err := fmt.Fprintf(l.w, "[%s] %s %s %s\n",
		e.Timestamp.Format(time.RFC3339), e.Level, e.Event, e.Detail)
	return err
}

// Info is a convenience wrapper for LevelInfo.
func (l *Logger) Info(event, detail string) error {
	return l.Log(LevelInfo, event, detail)
}

// Warn is a convenience wrapper for LevelWarn.
func (l *Logger) Warn(event, detail string) error {
	return l.Log(LevelWarn, event, detail)
}

// Error is a convenience wrapper for LevelError.
func (l *Logger) Error(event, detail string) error {
	return l.Log(LevelError, event, detail)
}
