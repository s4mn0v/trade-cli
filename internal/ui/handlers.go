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
		m.addHistoryEntry("SHORT", "65100.20")
		m.OrderMode = false
	}
	return nil
}

func (m *Manager) HandleLong(g *gocui.Gui, v *gocui.View) error {
	if m.OrderMode {
		m.addHistoryEntry("LONG ", "65250.40")
		m.OrderMode = false
	}
	return nil
}

// Internal helper to update the table data
func (m *Manager) addHistoryEntry(side, price string) {
	newEntry := HistoryEntry{
		Time:   time.Now().Format("15:04"),
		Side:   side,
		Price:  price,
		Status: "FILLED",
	}
	// Prepend to show newest at the top
	m.History = append([]HistoryEntry{newEntry}, m.History...)
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
