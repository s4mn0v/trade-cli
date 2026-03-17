package ui

import "github.com/awesome-gocui/gocui"

// InitKeybindings centralizes all input mapping
func (m *Manager) InitKeybindings(g *gocui.Gui) error {
	// --- GLOBAL KEYS ---
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, m.Quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, m.ToggleFocus); err != nil {
		return err
	}

	// --- ORDER PANEL KEYS ---
	// Mode Entry
	if err := g.SetKeybinding("order_panel", gocui.KeyCtrlO, gocui.ModNone, m.EnterOrderMode); err != nil {
		return err
	}
	// Action Keys (S and L)
	if err := g.SetKeybinding("order_panel", 's', gocui.ModNone, m.HandleShort); err != nil {
		return err
	}
	if err := g.SetKeybinding("order_panel", 'l', gocui.ModNone, m.HandleLong); err != nil {
		return err
	}

	// --- HISTORY PANEL KEYS ---
	if err := g.SetKeybinding("history", gocui.KeyArrowUp, gocui.ModNone, m.ScrollUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("history", gocui.KeyArrowDown, gocui.ModNone, m.ScrollDown); err != nil {
		return err
	}

	return nil
}
