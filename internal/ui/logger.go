package ui

import (
	"fmt"
	"sync"
	"time"

	"github.com/awesome-gocui/gocui"
)

type LogMessage struct {
	Timestamp string
	Text      string
	IsError   bool
}

type UILogger struct {
	mu       sync.RWMutex
	Messages []LogMessage
}

func NewUILogger() *UILogger {
	return &UILogger{Messages: []LogMessage{}}
}

func (l *UILogger) Info(msg string)  { l.add(msg, false) }
func (l *UILogger) Error(msg string) { l.add(msg, true) }

func (l *UILogger) add(text string, isErr bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Messages = append(l.Messages, LogMessage{
		Timestamp: time.Now().Format("15:04:05"),
		Text:      text,
		IsError:   isErr,
	})
}

// Clear removes all messages from the logger
func (l *UILogger) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Messages = []LogMessage{}
}

func (l *UILogger) Render(v *gocui.View) {
	v.Clear()
	l.mu.RLock()
	defer l.mu.RUnlock()
	for _, m := range l.Messages {
		color := "\033[32m" // Green
		if m.IsError {
			color = "\033[31m"
		} // Red
		fmt.Fprintf(v, "[%s] %s%s\033[0m\n", m.Timestamp, color, m.Text)
	}
}
