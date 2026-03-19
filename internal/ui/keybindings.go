package ui

import "github.com/awesome-gocui/gocui"

func (m *Manager) InitKeybindings(g *gocui.Gui) error {
	// Global Actions
	_ = g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, m.Quit)
	_ = g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, m.ToggleFocus)
	_ = g.SetKeybinding("", gocui.KeyCtrlL, gocui.ModNone, m.ClearLogs)

	// --- Mode Switching ---
	_ = g.SetKeybinding("", gocui.KeyCtrlS, gocui.ModNone, m.SetModeSpot)
	_ = g.SetKeybinding("", gocui.KeyCtrlF, gocui.ModNone, m.SetModeFutures)

	// Order Mode
	_ = g.SetKeybinding("", gocui.KeyCtrlO, gocui.ModNone, m.EnterOrderMode)

	// Actions: B/L for Positive direction, S for Negative direction
	_ = g.SetKeybinding("", 'b', gocui.ModNone, m.HandleAction1)
	_ = g.SetKeybinding("", 'l', gocui.ModNone, m.HandleAction1)
	_ = g.SetKeybinding("", 's', gocui.ModNone, m.HandleAction2)

	// History Panel Navigation (Selection)
	_ = g.SetKeybinding("history", gocui.KeyArrowUp, gocui.ModNone, m.HistoryUp)
	_ = g.SetKeybinding("history", gocui.KeyArrowDown, gocui.ModNone, m.HistoryDown)

	// Leverage Popup
	_ = g.SetKeybinding("", 'L', gocui.ModNone, m.ToggleLeverage)
	_ = g.SetKeybinding("leverage_pop", 'k', gocui.ModNone, m.LeverageUp)
	_ = g.SetKeybinding("leverage_pop", 'j', gocui.ModNone, m.LeverageDown)
	_ = g.SetKeybinding("leverage_pop", gocui.KeyEsc, gocui.ModNone, m.CloseLeverage)
	_ = g.SetKeybinding("leverage_pop", gocui.KeyEnter, gocui.ModNone, m.ConfirmLeverage)
	_ = g.SetKeybinding("leverage_pop", 'r', gocui.ModNone, m.ResetLeverage)

	// Coin Popup
	_ = g.SetKeybinding("", 'p', gocui.ModNone, m.ToggleCoinPopup)
	_ = g.SetKeybinding("coin_pop", gocui.KeyEnter, gocui.ModNone, m.ConfirmCoin)
	_ = g.SetKeybinding("coin_pop", gocui.KeyEsc, gocui.ModNone, m.ToggleCoinPopup)

	// Quantity Panel
	_ = g.SetKeybinding("", 'Q', gocui.ModNone, m.ToggleQuantity)
	_ = g.SetKeybinding("quantity_pop", 'k', gocui.ModNone, m.QuantityUp)
	_ = g.SetKeybinding("quantity_pop", 'j', gocui.ModNone, m.QuantityDown)
	_ = g.SetKeybinding("quantity_pop", gocui.KeyEsc, gocui.ModNone, m.CloseQuantity)
	_ = g.SetKeybinding("quantity_pop", gocui.KeyEnter, gocui.ModNone, m.ConfirmQuantity)
	_ = g.SetKeybinding("quantity_pop", 'r', gocui.ModNone, m.ResetQuantity)

	// Logs Panel Navigation (Scrolling)
	_ = g.SetKeybinding("logs", gocui.KeyArrowUp, gocui.ModNone, m.ScrollUp)
	_ = g.SetKeybinding("logs", gocui.KeyArrowDown, gocui.ModNone, m.ScrollDown)

	return nil
}
