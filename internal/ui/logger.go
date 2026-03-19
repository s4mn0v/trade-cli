package ui

import (
	"fmt"
	"sync"
	"time"

	"github.com/awesome-gocui/gocui"
)

type LogMessage struct {
	Timestamp, Text string
	Level           string // "INFO", "WARN", "ERR"
}

type UILogger struct {
	mu       sync.RWMutex
	Messages []LogMessage
}

func NewUILogger() *UILogger { return &UILogger{Messages: []LogMessage{}} }

func (l *UILogger) Info(msg string)    { l.add(msg, "INFO") }
func (l *UILogger) Warning(msg string) { l.add(msg, "WARN") }
func (l *UILogger) Error(msg string)   { l.add(msg, "ERR") }

func (l *UILogger) add(text string, level string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Messages = append(l.Messages, LogMessage{time.Now().Format("15:04:05"), text, level})
}

func (l *UILogger) Clear() {
	l.mu.Lock()
	l.Messages = []LogMessage{}
	l.mu.Unlock()
}

func (l *UILogger) Render(v *gocui.View) {
	v.Clear()
	l.mu.RLock()
	defer l.mu.RUnlock()
	for _, m := range l.Messages {
		color := "\033[32m" // Green (Info/Success)
		if m.Level == "ERR" {
			color = "\033[31m"
		} // Red (Error)
		if m.Level == "WARN" {
			color = "\033[33m"
		} // Yellow (Warning)

		_, _ = fmt.Fprintf(v, "[%s] %s%s\033[0m\n", m.Timestamp, color, m.Text)
	}
}
