package ui

import (
	"time"

	"github.com/awesome-gocui/gocui"
)

func (m *Manager) Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (m *Manager) ToggleFocus(g *gocui.Gui, v *gocui.View) error {
	target := "order_panel"
	if v != nil && v.Name() == "order_panel" {
		target = "history"
	}
	_, err := g.SetCurrentView(target)
	return err
}

func (m *Manager) EnterOrderMode(g *gocui.Gui, v *gocui.View) error {
	m.OrderMode = true
	return nil
}

func (m *Manager) HandleShort(g *gocui.Gui, v *gocui.View) error {
	if m.OrderMode {
		m.History.Add(HistoryEntry{
			Pair: "BTCUSDT", Date: time.Now().Format("01-02 15:04"),
			Direction: "SHORT", Price: "65100", Total: "0.01", Status: "FILLED",
		})
		m.OrderMode = false // Reset mode after order
	}
	return nil
}

func (m *Manager) HandleLong(g *gocui.Gui, v *gocui.View) error {
	if m.OrderMode {
		m.History.Add(HistoryEntry{
			Pair: "BTCUSDT", Date: time.Now().Format("01-02 15:04"),
			Direction: "LONG", Price: "65200", Total: "0.01", Status: "FILLED",
		})
		m.OrderMode = false // Reset mode after order
	}
	return nil
}

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
