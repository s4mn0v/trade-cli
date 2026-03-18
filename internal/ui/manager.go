package ui

import (
	"errors"
	"fmt"
	"sync"

	"github.com/awesome-gocui/gocui"
)

const (
	ModeSpot    = "SPOT"
	ModeFutures = "FUTURES"
)

type Manager struct {
	mu        sync.RWMutex
	OrderMode bool
	Mode      string // Tracks "SPOT" or "FUTURES"
	History   *HistoryTable
	Logger    *UILogger
}

func NewManager() *Manager {
	return &Manager{
		History: NewHistoryTable(),
		Logger:  NewUILogger(),
		Mode:    ModeSpot, // Default to Spot
	}
}

func (m *Manager) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	orderH, logW := 5, 30
	histW := maxX - logW - 1

	// 1. Order Panel
	if v, err := g.SetView("order_panel", 0, 0, maxX-1, orderH, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		g.SetCurrentView("order_panel")
	} else {
		v.Title = fmt.Sprintf(" Place order [%s] ", m.Mode)
		v.Clear()
		if m.Mode == ModeSpot {
			fmt.Fprint(v, "\n  (Ctrl+O, b) = Buy | (Ctrl+O, s) = Sell | (Ctrl+S) Spot | (Ctrl+F) Futures")
		} else {
			fmt.Fprint(v, "\n  (Ctrl+O, l) = Long | (Ctrl+O, s) = Short | (Ctrl+S) Spot | (Ctrl+F) Futures")
		}
	}

	// 2. History Panel
	if v, err := g.SetView("history", 0, orderH+1, histW, maxY-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " History "
	} else {
		m.History.Render(v, histW, m.Mode) // Pass mode to History
	}

	// 3. Logs
	if v, err := g.SetView("logs", histW+1, orderH+1, maxX-1, maxY-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " Logs "
	} else {
		m.Logger.Render(v)
	}

	m.applyDynamicStyles(g)
	return nil
}

func (m *Manager) applyDynamicStyles(g *gocui.Gui) {
	curr := ""
	if v := g.CurrentView(); v != nil {
		curr = v.Name()
	}

	// Default Mode Colors
	modeColor := gocui.ColorGreen // SPOT
	if m.Mode == ModeFutures {
		modeColor = gocui.ColorRed // FUTURES
	}

	views := []string{"order_panel", "history", "logs"}
	for _, name := range views {
		if v, err := g.View(name); err == nil {
			if (name == "order_panel" && m.OrderMode) || curr == name {
				v.FrameColor = gocui.ColorYellow // Active/Warning
			} else {
				v.FrameColor = modeColor // Mode Color
			}
		}
	}
}

func SetTerminalSize(rows, cols int) {
	fmt.Printf("\033[8;%d;%dt", rows, cols)
}
