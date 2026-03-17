package ui

import (
	"fmt"
	"strings"
	"sync"

	"github.com/awesome-gocui/gocui"
)

type HistoryEntry struct {
	Pair, Date, Direction, Price, Total, Status string
}

type HistoryTable struct {
	mu      sync.RWMutex
	Entries []HistoryEntry
}

func NewHistoryTable() *HistoryTable {
	return &HistoryTable{Entries: []HistoryEntry{}}
}

func (h *HistoryTable) Add(e HistoryEntry) {
	h.mu.Lock()
	defer h.mu.Unlock()
	// Prepend: Newest entries at the top
	h.Entries = append([]HistoryEntry{e}, h.Entries...)
}

func (h *HistoryTable) Render(v *gocui.View, width int) {
	v.Clear()
	// Header: PAIR | DATE | DIR | PRICE | TOTAL | STATUS
	fmt.Fprintf(v, "  \033[1m%-8s %-11s %-9s %-10s %-8s %-8s\033[0m\n",
		"PAIR", "DATE", "DIR", "PRICE", "TOTAL", "STATUS")
	fmt.Fprintln(v, "  "+strings.Repeat("─", width-6))

	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, e := range h.Entries {
		dirColor := "\033[32m" // Green for LONG
		if e.Direction == "SHORT" {
			dirColor = "\033[31m"
		}

		fmt.Fprintf(v, "  %-8s %-11s %s%-9s\033[0m %-10s %-8s %-8s\n",
			e.Pair, e.Date, dirColor, e.Direction, e.Price, e.Total, e.Status)
	}
}
