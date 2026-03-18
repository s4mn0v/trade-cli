package ui

import (
	"time"

	"github.com/awesome-gocui/gocui"
)

func (m *Manager) Quit(g *gocui.Gui, v *gocui.View) error { return gocui.ErrQuit }

func (m *Manager) ToggleFocus(g *gocui.Gui, v *gocui.View) error {
	order, history, logs := "order_panel", "history", "logs"
	target := history
	if v != nil {
		switch v.Name() {
		case order:
			target = history
		case history:
			target = logs
		case logs:
			target = order
		}
	}
	_, err := g.SetCurrentView(target)
	return err
}

func (m *Manager) EnterOrderMode(g *gocui.Gui, v *gocui.View) error {
	m.mu.Lock()
	m.OrderMode = true
	m.mu.Unlock()
	return nil
}

func (m *Manager) HandleShort(g *gocui.Gui, v *gocui.View) error {
	if m.OrderMode {
		m.History.Add(HistoryEntry{
			Pair: "BTCUSDT", Date: time.Now().Format("01-02 15:04"),
			Direction: "SHORT", Price: "65100.20", Total: "0.01", Status: "FILLED",
		})
		m.Logger.Info("SHORT executed")
		m.OrderMode = false
	}
	return nil
}

func (m *Manager) HandleLong(g *gocui.Gui, v *gocui.View) error {
	if m.OrderMode {
		m.History.Add(HistoryEntry{
			Pair: "BTCUSDT", Date: time.Now().Format("01-02 15:04"),
			Direction: "LONG", Price: "65250.40", Total: "0.01", Status: "FILLED",
		})
		m.Logger.Info("LONG executed")
		m.OrderMode = false
	}
	return nil
}

func (m *Manager) HistoryUp(g *gocui.Gui, v *gocui.View) error {
	m.History.mu.Lock()
	if m.History.SelectedIdx > 0 {
		m.History.SelectedIdx--
	}
	m.History.mu.Unlock()
	return nil
}

func (m *Manager) HistoryDown(g *gocui.Gui, v *gocui.View) error {
	m.History.mu.Lock()
	if m.History.SelectedIdx < len(m.History.Entries)-1 {
		m.History.SelectedIdx++
	}
	m.History.mu.Unlock()
	return nil
}

func (m *Manager) SetModeSpot(g *gocui.Gui, v *gocui.View) error {
	m.mu.Lock()
	m.Mode = ModeSpot
	m.mu.Unlock()

	m.History.Reset() // Wipes the table data
	m.Logger.Info("Switched to SPOT. Table cleared (Fetching Spot API...)")
	return nil
}

func (m *Manager) SetModeFutures(g *gocui.Gui, v *gocui.View) error {
	m.mu.Lock()
	m.Mode = ModeFutures
	m.mu.Unlock()

	m.History.Reset() // Wipes the table data
	m.Logger.Info("Switched to FUTURES. Table cleared (Fetching Futures API...)")
	return nil
}

func (m *Manager) HandleAction1(g *gocui.Gui, v *gocui.View) error { // BUY or LONG
	if !m.OrderMode {
		return nil
	}

	direction := "LONG"
	if m.Mode == ModeSpot {
		direction = "BUY"
	}

	m.History.Add(HistoryEntry{
		Pair: "BTCUSDT", Date: time.Now().Format("01-02 15:04"),
		Direction: direction, Price: "65000.00", Total: "0.01", Status: "FILLED",
	})
	m.Logger.Info(direction + " order executed")
	m.OrderMode = false
	return nil
}

func (m *Manager) HandleAction2(g *gocui.Gui, v *gocui.View) error { // SELL or SHORT
	if !m.OrderMode {
		return nil
	}

	direction := "SHORT"
	if m.Mode == ModeSpot {
		direction = "SELL"
	}

	m.History.Add(HistoryEntry{
		Pair: "BTCUSDT", Date: time.Now().Format("01-02 15:04"),
		Direction: direction, Price: "65000.00", Total: "0.01", Status: "FILLED",
	})
	m.Logger.Info(direction + " order executed")
	m.OrderMode = false
	return nil
}

// Scroll methods for the Logs panel
func (m *Manager) ScrollUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		if oy > 0 {
			v.SetOrigin(ox, oy-1)
		}
	}
	return nil
}

func (m *Manager) ScrollDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		v.SetOrigin(ox, oy+1)
	}
	return nil
}

// ClearLogs resets the log panel
func (m *Manager) ClearLogs(g *gocui.Gui, v *gocui.View) error {
	m.Logger.Clear()
	return nil
}
