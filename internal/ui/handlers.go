package ui

import (
	"time"

	"github.com/awesome-gocui/gocui"
)

// Quit closes the application
func (m *Manager) Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// ToggleFocus cycles between the Order Panel and the History Panel
func (m *Manager) ToggleFocus(g *gocui.Gui, v *gocui.View) error {
	target := "history"
	if v != nil && v.Name() == "history" {
		target = "order_panel"
	}

	if _, err := g.SetCurrentView(target); err != nil {
		return err
	}
	return nil
}

// EnterOrderMode activates the yellow frame state
func (m *Manager) EnterOrderMode(g *gocui.Gui, v *gocui.View) error {
	m.mu.Lock()
	m.OrderMode = true
	m.mu.Unlock()
	// Ensure we are focused on the order panel to see the color change
	g.SetCurrentView("order_panel")
	return nil
}

// HandleShort executes a short and adds it to the history component
func (m *Manager) HandleShort(g *gocui.Gui, v *gocui.View) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.OrderMode {
		m.History.Add(HistoryEntry{
			Pair:      "BTCUSDT",
			Date:      time.Now().Format("01-02 15:04"),
			Direction: "SHORT",
			Price:     "65100.20",
			Total:     "0.01",
			Status:    "FILLED",
		})
		m.OrderMode = false
	}
	return nil
}

// HandleLong executes a long and adds it to the history component
func (m *Manager) HandleLong(g *gocui.Gui, v *gocui.View) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.OrderMode {
		m.History.Add(HistoryEntry{
			Pair:      "BTCUSDT",
			Date:      time.Now().Format("01-02 15:04"),
			Direction: "LONG",
			Price:     "65250.40",
			Total:     "0.01",
			Status:    "FILLED",
		})
		m.OrderMode = false
	}
	return nil
}

// HistoryUp moves the row selection up
func (m *Manager) HistoryUp(g *gocui.Gui, v *gocui.View) error {
	m.History.mu.Lock()
	defer m.History.mu.Unlock()
	if m.History.SelectedIdx > 0 {
		m.History.SelectedIdx--
	}
	return nil
}

// HistoryDown moves the row selection down
func (m *Manager) HistoryDown(g *gocui.Gui, v *gocui.View) error {
	m.History.mu.Lock()
	defer m.History.mu.Unlock()
	if m.History.SelectedIdx < len(m.History.Entries)-1 {
		m.History.SelectedIdx++
	}
	return nil
}
