package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/awesome-gocui/gocui"
)

func (m *Manager) Quit(g *gocui.Gui, v *gocui.View) error { return gocui.ErrQuit }

func (m *Manager) ToggleFocus(g *gocui.Gui, v *gocui.View) error {
	order, history, logs := "order_panel", "history", "logs"
	target := history
	if v != nil {
		switch v.Name() {
		case order:
			target = history
		case history:
			target = logs
		case logs:
			target = order
		}
	}
	_, err := g.SetCurrentView(target)
	return err
}

func (m *Manager) HandleShort(g *gocui.Gui, v *gocui.View) error {
	if m.OrderMode {
		m.History.Add(HistoryEntry{
			Pair: "BTCUSDT", Date: time.Now().Format("01-02 15:04"),
			Direction: "SHORT", Price: "65100.20", Total: "0.01", Status: "FILLED",
		})
		m.Logger.Info("SHORT executed")
		m.OrderMode = false
	}
	return nil
}

func (m *Manager) HandleLong(g *gocui.Gui, v *gocui.View) error {
	if m.OrderMode {
		m.History.Add(HistoryEntry{
			Pair: "BTCUSDT", Date: time.Now().Format("01-02 15:04"),
			Direction: "LONG", Price: "65250.40", Total: "0.01", Status: "FILLED",
		})
		m.Logger.Info("LONG executed")
		m.OrderMode = false
	}
	return nil
}

func (m *Manager) HistoryUp(g *gocui.Gui, v *gocui.View) error {
	m.History.mu.Lock()
	if m.History.SelectedIdx > 0 {
		m.History.SelectedIdx--
	}
	m.History.mu.Unlock()
	return nil
}

func (m *Manager) HistoryDown(g *gocui.Gui, v *gocui.View) error {
	m.History.mu.Lock()
	if m.History.SelectedIdx < len(m.History.Entries)-1 {
		m.History.SelectedIdx++
	}
	m.History.mu.Unlock()
	return nil
}

func (m *Manager) SetModeSpot(g *gocui.Gui, v *gocui.View) error {
	m.mu.Lock()
	m.Mode = ModeSpot
	m.mu.Unlock()
	m.History.Reset()
	m.Logger.Info("Mode set to SPOT (Success)")
	return nil
}

func (m *Manager) SetModeFutures(g *gocui.Gui, v *gocui.View) error {
	m.mu.Lock()
	m.Mode = ModeFutures
	m.mu.Unlock()
	m.History.Reset()
	m.Logger.Warning("Mode set to FUTURES (Warning: High Risk)")
	return nil
}

func (m *Manager) HandleAction1(g *gocui.Gui, v *gocui.View) error {
	if !m.OrderMode {
		return nil
	}

	m.mu.Lock()
	currentCoin := m.CurrentCoin
	percent := m.PositionPercent
	mode := m.Mode
	m.OrderMode = false // Exit order mode
	m.mu.Unlock()

	if mode == ModeSpot {
		// --- SPOT LOGIC ---
		m.History.Add(HistoryEntry{
			Pair:      currentCoin,
			Date:      time.Now().Format("01-02 15:04"),
			Direction: "BUY",
			Price:     "65000.00",
			Total:     fmt.Sprintf("%.2f", m.SpotBalance*(float64(percent)/100.0)),
			Status:    "FILLED",
		})
		m.Logger.Info(fmt.Sprintf("Spot BUY %s Filled", currentCoin))

	} else {
		// --- FUTURES LOGIC ---
		posID := int(time.Now().UnixNano())
		newPos := &Position{
			ID:    posID,
			Pair:  currentCoin,
			Side:  "LONG",
			Entry: "65000.00",
			Size:  fmt.Sprintf("%d%%", percent),
			PnL:   0.00,
		}

		m.Positions.mu.Lock()
		m.Positions.Active = append(m.Positions.Active, newPos)
		m.Positions.mu.Unlock()
		m.Logger.Info(fmt.Sprintf("Futures LONG %s Opened", currentCoin))

		// TEST TIMER: Auto-close after 10 seconds
		time.AfterFunc(10*time.Second, func() {
			m.removePositionByID(posID, "Auto-Closed (Expired)")
		})
	}
	return nil
}

func (m *Manager) HandleAction2(g *gocui.Gui, v *gocui.View) error {
	if !m.OrderMode {
		return nil
	}

	m.mu.Lock()
	currentCoin := m.CurrentCoin
	percent := m.PositionPercent
	mode := m.Mode
	m.OrderMode = false
	m.mu.Unlock()

	if mode == ModeSpot {
		// --- SPOT LOGIC ---
		m.History.Add(HistoryEntry{
			Pair:      currentCoin,
			Date:      time.Now().Format("01-02 15:04"),
			Direction: "SELL",
			Price:     "65000.00",
			Total:     fmt.Sprintf("%.2f", m.SpotBalance*(float64(percent)/100.0)),
			Status:    "FILLED",
		})
		m.Logger.Error(fmt.Sprintf("Spot SELL %s Filled", currentCoin))

	} else {
		// --- FUTURES LOGIC ---
		posID := int(time.Now().UnixNano())
		newPos := &Position{
			ID:    posID,
			Pair:  currentCoin,
			Side:  "SHORT",
			Entry: "65000.00",
			Size:  fmt.Sprintf("%d%%", percent),
			PnL:   0.00,
		}

		m.Positions.mu.Lock()
		m.Positions.Active = append(m.Positions.Active, newPos)
		m.Positions.mu.Unlock()
		m.Logger.Info(fmt.Sprintf("Futures SHORT %s Opened", currentCoin))

		// TEST TIMER: Auto-close after 10 seconds
		time.AfterFunc(10*time.Second, func() {
			m.removePositionByID(posID, "Auto-Closed (Expired)")
		})
	}
	return nil
}

// Helper function to handle thread-safe removal for both Timer and Manual Close
func (m *Manager) removePositionByID(id int, reason string) {
	m.Positions.mu.Lock()
	defer m.Positions.mu.Unlock()

	for i, p := range m.Positions.Active {
		if p.ID == id {
			// Remove from slice
			m.Positions.Active = append(m.Positions.Active[:i], m.Positions.Active[i+1:]...)

			// Reset selection index if it's now out of bounds
			if m.Positions.SelectedIdx >= len(m.Positions.Active) && len(m.Positions.Active) > 0 {
				m.Positions.SelectedIdx = len(m.Positions.Active) - 1
			}

			m.Logger.Warning(fmt.Sprintf("%s: %s", reason, p.Pair))
			return
		}
	}
}

// Update the manual close handler to use the helper
func (m *Manager) CloseActivePosition(g *gocui.Gui, v *gocui.View) error {
	if m.Mode != ModeFutures {
		return nil
	}

	m.Positions.mu.RLock()
	if len(m.Positions.Active) == 0 {
		m.Positions.mu.RUnlock()
		return nil
	}
	targetID := m.Positions.Active[m.Positions.SelectedIdx].ID
	m.Positions.mu.RUnlock()

	m.removePositionByID(targetID, "Manually Closed")
	return nil
}

// Navigation for positions
func (m *Manager) PositionUp(g *gocui.Gui, v *gocui.View) error {
	m.Positions.mu.Lock()
	if m.Positions.SelectedIdx > 0 {
		m.Positions.SelectedIdx--
	}
	m.Positions.mu.Unlock()
	return nil
}

func (m *Manager) PositionDown(g *gocui.Gui, v *gocui.View) error {
	m.Positions.mu.Lock()
	if m.Positions.SelectedIdx < len(m.Positions.Active)-1 {
		m.Positions.SelectedIdx++
	}
	m.Positions.mu.Unlock()
	return nil
}

func (m *Manager) EnterOrderMode(g *gocui.Gui, v *gocui.View) error {
	m.mu.Lock()
	m.OrderMode = true
	m.mu.Unlock()
	m.Logger.Warning("Entering Order Mode... (Awaiting Action)") // Yellow Log
	return nil
}

// ScrollUp methods for the Logs panel
func (m *Manager) ScrollUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		if oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return nil
			}
		}
	}
	return nil
}

func (m *Manager) ScrollDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+1); err != nil {
			return nil
		}
	}
	return nil
}

// OpenLeverage checks if mode is Futures then opens popup

func (m *Manager) ConfirmLeverage(g *gocui.Gui, v *gocui.View) error {
	m.mu.Lock()
	m.FuturesLeverage = m.LeveragePopup.CurrentVal
	m.ShowLeverage = false
	m.mu.Unlock()

	m.Logger.Info(fmt.Sprintf("F/Leverage: %d", m.FuturesLeverage))

	_, err := g.SetCurrentView("order_panel")
	return err
}

func (m *Manager) ToggleLeverage(g *gocui.Gui, v *gocui.View) error {
	if m.Mode != ModeFutures {
		m.Logger.Error("Leverage only available in FUTURES mode")
		return nil
	}

	m.mu.Lock()
	m.ShowLeverage = !m.ShowLeverage
	if m.ShowLeverage {
		m.LeveragePopup.CurrentVal = m.FuturesLeverage
	}
	m.mu.Unlock()

	if !m.ShowLeverage {
		_, err := g.SetCurrentView("order_panel")
		return err
	}
	return nil
}

func (m *Manager) LeverageUp(g *gocui.Gui, v *gocui.View) error {
	if m.LeveragePopup.CurrentVal < 125 {
		m.LeveragePopup.CurrentVal++
	}
	return nil
}

func (m *Manager) LeverageDown(g *gocui.Gui, v *gocui.View) error {
	if m.LeveragePopup.CurrentVal > 1 {
		m.LeveragePopup.CurrentVal--
	}
	return nil
}

func (m *Manager) ResetLeverage(g *gocui.Gui, v *gocui.View) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.LeveragePopup.CurrentVal = 5
	return nil
}

func (m *Manager) CloseLeverage(g *gocui.Gui, v *gocui.View) error {
	m.ShowLeverage = false
	_, err := g.SetCurrentView("order_panel")
	return err
}

// ToggleQuantity opens the percentage popup
func (m *Manager) ToggleQuantity(g *gocui.Gui, v *gocui.View) error {
	m.mu.Lock()
	m.ShowQuantity = !m.ShowQuantity
	if m.ShowQuantity {
		m.QuantityPopup.CurrentVal = m.PositionPercent
	}
	m.mu.Unlock()

	if !m.ShowQuantity {
		_, err := g.SetCurrentView("order_panel")
		return err
	}
	return nil
}

func (m *Manager) QuantityUp(g *gocui.Gui, v *gocui.View) error {
	if m.QuantityPopup.CurrentVal < 100 {
		m.QuantityPopup.CurrentVal += 5 // Step by 5%
	}
	return nil
}

func (m *Manager) QuantityDown(g *gocui.Gui, v *gocui.View) error {
	if m.QuantityPopup.CurrentVal > 0 {
		m.QuantityPopup.CurrentVal -= 5
	}
	return nil
}

func (m *Manager) ConfirmQuantity(g *gocui.Gui, v *gocui.View) error {
	m.mu.Lock()
	m.PositionPercent = m.QuantityPopup.CurrentVal
	m.ShowQuantity = false
	m.mu.Unlock()

	m.Logger.Info(fmt.Sprintf("Size Set: %d%%", m.PositionPercent))
	_, err := g.SetCurrentView("order_panel")
	return err
}

func (m *Manager) ResetQuantity(g *gocui.Gui, v *gocui.View) error {
	m.QuantityPopup.CurrentVal = 10
	return nil
}

func (m *Manager) CloseQuantity(g *gocui.Gui, v *gocui.View) error {
	m.ShowQuantity = false
	_, err := g.SetCurrentView("order_panel")
	return err
}

func (m *Manager) ToggleCoinPopup(g *gocui.Gui, v *gocui.View) error {
	m.mu.Lock()
	m.ShowCoin = !m.ShowCoin
	m.mu.Unlock()

	if !m.ShowCoin {
		_, err := g.SetCurrentView("order_panel")
		return err
	}
	return nil
}

func (m *Manager) ConfirmCoin(g *gocui.Gui, v *gocui.View) error {
	rawInput := strings.TrimSpace(v.Buffer())
	input := strings.ToUpper(rawInput)

	matches := m.CoinPopup.GetMatches(input)

	var finalCoin string

	if m.CoinPopup.IsValid(input) {
		finalCoin = input
	} else if len(matches) > 0 {
		// User typed a partial name (e.g., "BT") -> Auto-complete to first match ("BTCUSDT")
		finalCoin = matches[0]
	} else {
		// No match found at all: Clear and let them try again
		v.Clear()
		_ = v.SetCursor(0, 0)
		m.Logger.Error(fmt.Sprintf("No matches for: %s", input))
		return nil
	}

	// 3. Set the coin and close
	m.mu.Lock()
	m.CurrentCoin = finalCoin
	m.ShowCoin = false
	m.mu.Unlock()

	m.Logger.Info(fmt.Sprintf("Coin set to: %s (Success)", finalCoin))

	_, err := g.SetCurrentView("order_panel")
	return err
}

// ClearLogs resets the log panel
func (m *Manager) ClearLogs(g *gocui.Gui, v *gocui.View) error {
	m.Logger.Clear()
	return nil
}

// SetBalances is a thread-safe way to update the account available balance
func (m *Manager) SetBalances(spot, futures float64) {
	m.mu.Lock()
	m.SpotBalance = spot
	m.FuturesBalance = futures
	m.mu.Unlock()
}
