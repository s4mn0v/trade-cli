package ui

import "github.com/awesome-gocui/gocui"

func (m *Manager) InitKeybindings(g *gocui.Gui) error {
	// Global Quit
	g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, m.Quit)

	// Tab to switch focus
	g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, m.ToggleFocus)

	// Order Mode Initiation (Must be in order_panel or global)
	g.SetKeybinding("", gocui.KeyCtrlO, gocui.ModNone, m.EnterOrderMode)

	// Action Keys (Global so they work as soon as OrderMode is true)
	g.SetKeybinding("", 's', gocui.ModNone, m.HandleShort)
	g.SetKeybinding("", 'l', gocui.ModNone, m.HandleLong)

	// Scrolling (Only when history is focused)
	g.SetKeybinding("history", gocui.KeyArrowUp, gocui.ModNone, m.ScrollUp)
	g.SetKeybinding("history", gocui.KeyArrowDown, gocui.ModNone, m.ScrollDown)

	return nil
}
