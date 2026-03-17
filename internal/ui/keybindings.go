package ui

import "github.com/awesome-gocui/gocui"

func (m *Manager) InitKeybindings(g *gocui.Gui) error {
	// Global Actions
	g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, m.Quit)
	g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, m.ToggleFocus)
	g.SetKeybinding("", gocui.KeyCtrlL, gocui.ModNone, m.ClearLogs)

	// Order Actions
	g.SetKeybinding("", gocui.KeyCtrlO, gocui.ModNone, m.EnterOrderMode)
	g.SetKeybinding("", 's', gocui.ModNone, m.HandleShort)
	g.SetKeybinding("", 'l', gocui.ModNone, m.HandleLong)

	// History Panel Navigation (Selection)
	g.SetKeybinding("history", gocui.KeyArrowUp, gocui.ModNone, m.HistoryUp)
	g.SetKeybinding("history", gocui.KeyArrowDown, gocui.ModNone, m.HistoryDown)

	// Logs Panel Navigation (Scrolling)
	g.SetKeybinding("logs", gocui.KeyArrowUp, gocui.ModNone, m.ScrollUp)
	g.SetKeybinding("logs", gocui.KeyArrowDown, gocui.ModNone, m.ScrollDown)

	return nil
}
