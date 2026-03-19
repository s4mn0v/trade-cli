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
	h.Entries = []HistoryEntry{}
	h.SelectedIdx = 0
}

func (h *HistoryTable) Render(v *gocui.View, width int, mode string) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	v.Clear()

	// 1. Header Logic
	col5 := "TOTAL"
	if mode == "FUTURES" {
		col5 = "PNL"
	}

	_, _ = fmt.Fprintf(v, "  \033[1m%-8s %-11s %-9s %-10s %-8s %-8s\033[0m\n", "PAIR", "DATE", "DIR", "PRICE", col5, "STATUS")
	_, _ = fmt.Fprintln(v, "  "+strings.Repeat("─", width-6))

	// 2. Scrolling Logic
	_, viewY := v.Size()
	headerRows := 2
	visibleHeight := viewY - headerRows
	ox, oy := v.Origin()
	if h.SelectedIdx < oy {
		_ = v.SetOrigin(ox, h.SelectedIdx)
	} else if h.SelectedIdx >= oy+visibleHeight {
		_ = v.SetOrigin(ox, h.SelectedIdx-visibleHeight+1)
	}

	// 3. Row Rendering
	for i, e := range h.Entries {
		if i == h.SelectedIdx {
			rowText := fmt.Sprintf("  %-8s %-11s %-9s %-10s %-8s %-8s ",
				e.Pair, e.Date, e.Direction, e.Price, e.Total, e.Status)
			_, _ = fmt.Fprintf(v, "\033[7m%s\033[0m\n", rowText)
		} else {
			// NORMAL: Apply Green/Red to the DIR column
			dirColor := "\033[32m" // Green
			if e.Direction == "SHORT" || e.Direction == "SELL" {
				dirColor = "\033[31m" // Red
			}
			_, _ = fmt.Fprintf(v, "  %-8s %-11s %s%-9s\033[0m %-10s %-8s %-8s\n",
				e.Pair, e.Date, dirColor, e.Direction, e.Price, e.Total, e.Status)
		}
	}
}
