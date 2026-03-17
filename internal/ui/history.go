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

func (h *HistoryTable) Render(v *gocui.View, width int) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	v.Clear()

	fmt.Fprintf(v, "  \033[1m%-8s %-11s %-9s %-10s %-8s %-8s\033[0m\n",
		"PAIR", "DATE", "DIR", "PRICE", "TOTAL", "STATUS")
	fmt.Fprintln(v, "  "+strings.Repeat("─", width-6))

	_, viewY := v.Size()
	headerRows := 2
	visibleHeight := viewY - headerRows

	ox, oy := v.Origin()
	if h.SelectedIdx < oy {
		v.SetOrigin(ox, h.SelectedIdx)
	} else if h.SelectedIdx >= oy+visibleHeight {
		v.SetOrigin(ox, h.SelectedIdx-visibleHeight+1)
	}

	for i, e := range h.Entries {

		dirColor := "\033[32m"

		if e.Direction == "SHORT" {
			dirColor = "\033[31m"
		}

		linePrefix := "  "
		lineSuffix := ""
		resetToNormal := "\033[0m"

		innerReset := "\033[0m"
		if i == h.SelectedIdx {
			linePrefix = "\033[30;47m  "
			innerReset = "\033[30;47m"

			lineSuffix = "\033[0m"

		}

		fmt.Fprintf(v, "%s%-8s %-11s %s%-9s%s %-10s %-8s %-8s %s\n",
			linePrefix,
			e.Pair,
			e.Date,
			dirColor, e.Direction, innerReset,

			e.Price,
			e.Total,
			e.Status,
			lineSuffix,
		)

		fmt.Fprint(v, resetToNormal)
	}
}
