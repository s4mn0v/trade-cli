package ui

import "github.com/awesome-gocui/gocui"

func (m *Manager) InitKeybindings(g *gocui.Gui) error {
	// Global Actions
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, m.Quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, m.ToggleFocus); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlL, gocui.ModNone, m.ClearLogs); err != nil {
		return err
	}

	// --- Mode Switching ---
	if err := g.SetKeybinding("", gocui.KeyCtrlS, gocui.ModNone, m.SetModeSpot); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlF, gocui.ModNone, m.SetModeFutures); err != nil {
		return err
	}

	// Order Mode
	if err := g.SetKeybinding("", gocui.KeyCtrlO, gocui.ModNone, m.EnterOrderMode); err != nil {
		return err
	}

	// Actions: B/L for Positive direction, S for Negative direction
	if err := g.SetKeybinding("", 'b', gocui.ModNone, m.HandleAction1); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'l', gocui.ModNone, m.HandleAction1); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 's', gocui.ModNone, m.HandleAction2); err != nil {
		return err
	}

	// History Panel Navigation (Selection)
	if err := g.SetKeybinding("history", gocui.KeyArrowUp, gocui.ModNone, m.HistoryUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("history", gocui.KeyArrowDown, gocui.ModNone, m.HistoryDown); err != nil {
		return err
	}

	// Leverage Panel
	if err := g.SetKeybinding("", 'L', gocui.ModNone, m.ToggleLeverage); err != nil {
		return err
	}
	if err := g.SetKeybinding("leverage_pop", 'k', gocui.ModNone, m.LeverageUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("leverage_pop", 'j', gocui.ModNone, m.LeverageDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("leverage_pop", gocui.KeyEsc, gocui.ModNone, m.CloseLeverage); err != nil {
		return err
	}
	if err := g.SetKeybinding("leverage_pop", gocui.KeyEnter, gocui.ModNone, m.ConfirmLeverage); err != nil {
		return err
	}
	if err := g.SetKeybinding("leverage_pop", 'r', gocui.ModNone, m.ResetLeverage); err != nil {
		return err
	}

	// Logs Panel Navigation (Scrolling)
	if err := g.SetKeybinding("logs", gocui.KeyArrowUp, gocui.ModNone, m.ScrollUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("logs", gocui.KeyArrowDown, gocui.ModNone, m.ScrollDown); err != nil {
		return err
	}

	return nil
}
