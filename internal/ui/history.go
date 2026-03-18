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
	mu          sync.RWMutex
	Entries     []HistoryEntry
	SelectedIdx int
}

func NewHistoryTable() *HistoryTable {
	return &HistoryTable{Entries: []HistoryEntry{}}
}

func (h *HistoryTable) Add(e HistoryEntry) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Entries = append([]HistoryEntry{e}, h.Entries...)
}

func (h *HistoryTable) Reset() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Entries = []HistoryEntry{} // Clear the slice
	h.SelectedIdx = 0            // Reset selection
}

func (h *HistoryTable) Render(v *gocui.View, width int, mode string) {
	v.Clear()

	// Dynamic Column Header
	col5 := "TOTAL"
	if mode == ModeFutures {
		col5 = "PNL"
	}

	fmt.Fprintf(v, "  \033[1m%-8s %-11s %-9s %-10s %-8s %-8s\033[0m\n",
		"PAIR", "DATE", "DIR", "PRICE", col5, "STATUS")
	fmt.Fprintln(v, "  "+strings.Repeat("─", width-6))

	h.mu.RLock()
	defer h.mu.RUnlock()
	for i, e := range h.Entries {
		dirColor := "\033[32m" // Default Green
		if e.Direction == "SHORT" || e.Direction == "SELL" {
			dirColor = "\033[31m"
		}

		line := fmt.Sprintf(" %-8s %-11s %s%-9s\033[39m %-10s %-8s %-8s ",
			e.Pair, e.Date, dirColor, e.Direction, e.Price, e.Total, e.Status)

		if i == h.SelectedIdx {
			fmt.Fprintf(v, "\033[30;47m%s\033[0m\n", line)
		} else {
			fmt.Fprintf(v, " %s\n", line)
		}
	}
}
