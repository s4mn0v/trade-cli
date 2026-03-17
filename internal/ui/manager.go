package ui

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/awesome-gocui/gocui"
)

type HistoryEntry struct {
	Time, Side, Price, Status string
}

type Manager struct {
	mu        sync.RWMutex
	OrderMode bool
	History   []HistoryEntry
}

func NewManager() *Manager {
	return &Manager{
		History: []HistoryEntry{}, // Start empty
	}
}

func SetTerminalSize(rows, cols int) {
	fmt.Printf("\033[8;%d;%dt", rows, cols)
}

func (m *Manager) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	orderHeight := 5

	// Panel 1: Order
	if v, err := g.SetView("order_panel", 0, 0, maxX-1, orderHeight, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " Place order "
		fmt.Fprint(v, "\n  (Ctrl + O, s) = Short | (Ctrl + O, l) = Long")
		g.SetCurrentView("order_panel")
	}

	// Panel 2: History
	if v, err := g.SetView("history", 0, orderHeight+1, maxX-1, maxY-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " History "
	} else {
		v.Clear()
		m.renderTable(v, maxX)
	}

	// Apply dynamic styling
	if v, err := g.View("order_panel"); err == nil {
		if m.OrderMode {
			v.FrameColor = gocui.ColorYellow
			v.Title = " Place order [ORDER MODE ACTIVE] "
		} else {
			v.FrameColor = gocui.ColorDefault
			v.Title = " Place order "
		}
	}

	return nil
}

func (m *Manager) renderTable(v *gocui.View, width int) {
	// Use Fprintf instead of Fprintln to process the %-8s specifiers
	// Added \n at the end of the string
	fmt.Fprintf(v, "  \033[1m%-8s %-10s %-12s %-10s\033[0m\n", "TIME", "SIDE", "PRICE", "STATUS")

	// Use Fprintln here because it's just a simple string
	fmt.Fprintln(v, "  "+strings.Repeat("─", width-6))

	for _, e := range m.History {
		color := "\033[32m" // Green
		if strings.Contains(e.Side, "SHORT") {
			color = "\033[31m" // Red
		}
		fmt.Fprintf(v, "  %-8s %s%-10s\033[0m %-12s %-10s\n", e.Time, color, e.Side, e.Price, e.Status)
	}
}
