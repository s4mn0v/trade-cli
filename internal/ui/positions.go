package ui

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/awesome-gocui/gocui"
)

type Position struct {
	ID        int
	Pair      string
	Side      string // LONG / SHORT
	Entry     string
	Size      string
	PnL       float64
	CreatedAt time.Time
}

type PositionList struct {
	mu          sync.RWMutex
	Active      []*Position
	SelectedIdx int
}

func NewPositionList() *PositionList {
	return &PositionList{Active: []*Position{}}
}

// Render draws the positions table inside the provided view
func (pl *PositionList) Render(v *gocui.View, width int) {
	pl.mu.RLock()
	defer pl.mu.RUnlock()

	if len(pl.Active) == 0 {
		_, _ = fmt.Fprintln(v, "\n\n  \033[90mNo active positions.\033[0m")
		return
	}

	// 1. Header Rendering
	_, _ = fmt.Fprintf(v, "\n\n  %-10s %-8s %-10s %-10s %-10s\n", "PAIR", "SIDE", "ENTRY", "SIZE", "PNL")
	_, _ = fmt.Fprintln(v, "  "+strings.Repeat("─", width-6))

	// 2. SCROLLING LOGIC
	headerRows := 6

	_, viewY := v.Size()
	visibleHeight := viewY - headerRows
	visibleHeight = max(visibleHeight, 1)

	ox, oy := v.Origin()

	if pl.SelectedIdx < oy {
		_ = v.SetOrigin(ox, pl.SelectedIdx)
	} else if pl.SelectedIdx >= oy+visibleHeight {
		_ = v.SetOrigin(ox, pl.SelectedIdx-visibleHeight+1)
	}

	// 3. ROW RENDERING
	for i, p := range pl.Active {
		prefix := "  "
		suffix := ""

		// Highlight the selected row
		if i == pl.SelectedIdx {
			prefix = "\033[7m >" // Invert color (Selected)
			suffix = "\033[0m"
		}

		sideColor := "\033[32m" // Green for Long
		if p.Side == "SHORT" {
			sideColor = "\033[31m" // Red for Short
		}

		pnlColor := "\033[32m"
		if p.PnL < 0 {
			pnlColor = "\033[31m"
		}

		_, _ = fmt.Fprintf(v, "%s%-10s %s%-8s\033[0m %-10s %-10s %s%+.2f%%\033[0m%s\n", prefix, p.Pair, sideColor, p.Side, p.Entry, p.Size, pnlColor, p.PnL, suffix)
	}
}
