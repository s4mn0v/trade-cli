package ui

import "github.com/awesome-gocui/gocui"

func (m *Manager) InitKeybindings(g *gocui.Gui) error {
	// Global
	g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, m.Quit)
	g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, m.ToggleFocus)

	// Order Actions
	g.SetKeybinding("", gocui.KeyCtrlO, gocui.ModNone, m.EnterOrderMode)
	g.SetKeybinding("", 's', gocui.ModNone, m.HandleShort)
	g.SetKeybinding("", 'l', gocui.ModNone, m.HandleLong)

	// History Row Selection (Must use the names defined in handlers.go)
	g.SetKeybinding("history", gocui.KeyArrowUp, gocui.ModNone, m.HistoryUp)
	g.SetKeybinding("history", gocui.KeyArrowDown, gocui.ModNone, m.HistoryDown)

	return nil
}
