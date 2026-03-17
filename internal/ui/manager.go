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
	Logger    *UILogger
}

func NewManager() *Manager {
	return &Manager{
		History: NewHistoryTable(),
		Logger:  NewUILogger(),
	}
}

func (m *Manager) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	orderH := 5
	logW := 30
	histW := maxX - logW - 1

	if v, err := g.SetView("order_panel", 0, 0, maxX-1, orderH, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " Place order "
		fmt.Fprint(v, "\n  (Ctrl+O, s) = Short | (Ctrl+O, l) = Long | (Ctrl+L) = Clear Logs")
		g.SetCurrentView("order_panel")
	}

	if v, err := g.SetView("history", 0, orderH+1, histW, maxY-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " History "
	} else {
		m.History.Render(v, histW)
	}

	if v, err := g.SetView("logs", histW+1, orderH+1, maxX-1, maxY-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " Logs "
		v.Autoscroll = true
	} else {
		m.Logger.Render(v)
	}

	// Dynamic Focus Styling
	curr := ""
	if v := g.CurrentView(); v != nil {
		curr = v.Name()
	}
	for _, name := range []string{"order_panel", "history", "logs"} {
		if v, err := g.View(name); err == nil {
			if name == "order_panel" && m.OrderMode {
				v.FrameColor = gocui.ColorYellow
			} else if curr == name {
				v.FrameColor = gocui.ColorCyan
			} else {
				v.FrameColor = gocui.ColorDefault
			}
		}
	}
	return nil
}

func SetTerminalSize(rows, cols int) {
	fmt.Printf("\033[8;%d;%dt", rows, cols)
}
