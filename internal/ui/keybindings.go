package ui

import "github.com/awesome-gocui/gocui"

func (m *Manager) InitKeybindings(g *gocui.Gui) error {
	// Global Actions
	g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, m.Quit)
	g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, m.ToggleFocus)
	g.SetKeybinding("", gocui.KeyCtrlL, gocui.ModNone, m.ClearLogs)

	// --- Mode Switching ---
	g.SetKeybinding("", gocui.KeyCtrlS, gocui.ModNone, m.SetModeSpot)
	g.SetKeybinding("", gocui.KeyCtrlF, gocui.ModNone, m.SetModeFutures)

	// Order Mode
	g.SetKeybinding("", gocui.KeyCtrlO, gocui.ModNone, m.EnterOrderMode)

	// Actions: B/L for Positive direction, S for Negative direction
	g.SetKeybinding("", 'b', gocui.ModNone, m.HandleAction1) // Buy
	g.SetKeybinding("", 'l', gocui.ModNone, m.HandleAction1) // Long
	g.SetKeybinding("", 's', gocui.ModNone, m.HandleAction2) // Sell or Short

	// History Panel Navigation (Selection)
	g.SetKeybinding("history", gocui.KeyArrowUp, gocui.ModNone, m.HistoryUp)
	g.SetKeybinding("history", gocui.KeyArrowDown, gocui.ModNone, m.HistoryDown)

	// Logs Panel Navigation (Scrolling)
	g.SetKeybinding("logs", gocui.KeyArrowUp, gocui.ModNone, m.ScrollUp)
	g.SetKeybinding("logs", gocui.KeyArrowDown, gocui.ModNone, m.ScrollDown)

	return nil
}
