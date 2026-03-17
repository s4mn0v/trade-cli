package ui

import (
	"errors"
	"fmt"
	"sync"

	"github.com/awesome-gocui/gocui"
)

type Manager struct {
	mu        sync.RWMutex
	OrderMode bool
	History   *HistoryTable
}

func NewManager() *Manager {
	return &Manager{History: NewHistoryTable()}
}

func (m *Manager) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	orderH := 5

	// 1. Order Panel
	if v, err := g.SetView("order_panel", 0, 0, maxX-1, orderH, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " Place order "
		fmt.Fprint(v, "\n  (Ctrl + O, s) = Short | (Ctrl + O, l) = Long")
		g.SetCurrentView("order_panel")
	}

	// 2. History Panel
	if v, err := g.SetView("history", 0, orderH+1, maxX-1, maxY-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " History "
		v.Autoscroll = false // Must be false
	} else {
		m.History.Render(v, maxX)
	}

	// --- FOCUS HIGHLIGHTING ---
	curr := ""
	if v := g.CurrentView(); v != nil {
		curr = v.Name()
	}

	if v, err := g.View("order_panel"); err == nil {
		if m.OrderMode {
			v.FrameColor, v.Title = gocui.ColorYellow, " Place order [ORDER MODE ACTIVE] "
		} else if curr == "order_panel" {
			v.FrameColor = gocui.ColorCyan
		} else {
			v.FrameColor = gocui.ColorDefault
			v.Title = " Place order "
		}
	}

	if v, err := g.View("history"); err == nil {
		if curr == "history" {
			v.FrameColor = gocui.ColorCyan
		} else {
			v.FrameColor = gocui.ColorDefault
		}
	}

	return nil
}

func SetTerminalSize(rows, cols int) {
	fmt.Printf("\033[8;%d;%dt", rows, cols)
}
