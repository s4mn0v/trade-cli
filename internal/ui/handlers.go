package ui

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/s4mn0v/bitget/config"
	v2 "github.com/s4mn0v/bitget/pkg/client/v2"
)

func (m *Manager) Quit(g *gocui.Gui, v *gocui.View) error { return gocui.ErrQuit }

func (m *Manager) ToggleFocus(g *gocui.Gui, v *gocui.View) error {
	// NEW: If a popup is open, ignore the Tab key
	if m.AnyPopupOpen() {
		return nil
	}

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
	baseAsset := strings.TrimSuffix(currentCoin, "USDT")
	price := 65000.0
	mode := m.Mode
	m.mu.Unlock()

	if mode == ModeSpot {
		// Calculate how much USDT spending
		costUSDT := m.SpotBalance * (float64(percent) / 100.0)

		if costUSDT <= 0 {
			m.Logger.Error("Insufficient USDT Balance")
			m.mu.Lock()
			m.OrderMode = false
			m.mu.Unlock()
			return nil
		}

		m.mu.Lock()
		m.SpotBalance -= costUSDT
		boughtAmount := costUSDT / price
		m.SpotAssets[baseAsset] += boughtAmount
		m.mu.Unlock()

		m.History.Add(HistoryEntry{
			Pair: currentCoin, Date: time.Now().Format("01-02 15:04"),
			Direction: "BUY", Price: fmt.Sprintf("%.2f", price),
			Total: fmt.Sprintf("%.2f", costUSDT), Status: "FILLED",
		})
		m.Logger.Info(fmt.Sprintf("Bought %.4f %s", boughtAmount, baseAsset))

	} else {
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

	m.mu.Lock()
	m.OrderMode = false
	m.mu.Unlock()
	return nil
}

func (m *Manager) HandleAction2(g *gocui.Gui, v *gocui.View) error {
	if !m.OrderMode {
		return nil
	}

	m.mu.Lock()
	currentCoin := m.CurrentCoin
	percent := m.PositionPercent
	baseAsset := strings.TrimSuffix(currentCoin, "USDT")
	price := 65000.0
	mode := m.Mode
	m.mu.Unlock()

	if mode == ModeSpot {
		// Calculate how much of the Asset are selling
		amountToSell := m.SpotAssets[baseAsset] * (float64(percent) / 100.0)

		if amountToSell <= 0 {
			m.Logger.Error(fmt.Sprintf("No %s to sell", baseAsset))
			m.mu.Lock()
			m.OrderMode = false
			m.mu.Unlock()
			return nil
		}

		m.mu.Lock()
		m.SpotAssets[baseAsset] -= amountToSell
		receivedUSDT := amountToSell * price
		m.SpotBalance += receivedUSDT
		m.mu.Unlock()

		m.History.Add(HistoryEntry{
			Pair: currentCoin, Date: time.Now().Format("01-02 15:04"),
			Direction: "SELL", Price: fmt.Sprintf("%.2f", price),
			Total: fmt.Sprintf("%.2f", receivedUSDT), Status: "FILLED",
		})
		m.Logger.Error(fmt.Sprintf("Sold %.4f %s", amountToSell, baseAsset))

	} else {
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

	m.mu.Lock()
	m.OrderMode = false
	m.mu.Unlock()
	return nil
}

// Helper function to handle thread-safe removal for both Timer and Manual Close
func (m *Manager) removePositionByID(id int, reason string) {
	m.Positions.mu.Lock()
	defer m.Positions.mu.Unlock()

	for i, p := range m.Positions.Active {
		if p.ID == id {
			// 1. Create a History Entry before deleting
			m.History.Add(HistoryEntry{
				Pair:      p.Pair,
				Date:      time.Now().Format("01-02 15:04"),
				Direction: p.Side,
				Price:     p.Entry,
				Total:     fmt.Sprintf("%+.2f%%", p.PnL),
				Status:    "CLOSED",
			})

			// 2. Remove from active positions slice
			m.Positions.Active = append(m.Positions.Active[:i], m.Positions.Active[i+1:]...)

			// 3. Reset selection index if out of bounds
			if m.Positions.SelectedIdx >= len(m.Positions.Active) && len(m.Positions.Active) > 0 {
				m.Positions.SelectedIdx = len(m.Positions.Active) - 1
			}

			// 4. Log the event
			m.Logger.Warning(fmt.Sprintf("%s: %s (PnL: %.2f%%)", reason, p.Pair, p.PnL))
			return
		}
	}
}

// CloseActivePosition

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

// SetBalances thread-safe way to update the account available balance
func (m *Manager) SetBalances(spot, futures float64) {
	m.mu.Lock()
	m.SpotBalance = spot
	m.FuturesBalance = futures
	m.mu.Unlock()
}

// ToggleSync Timers
func (m *Manager) ToggleSync(g *gocui.Gui, v *gocui.View) error {
	m.mu.Lock()
	m.ShowSync = !m.ShowSync
	m.mu.Unlock()

	if !m.ShowSync {
		_ = g.DeleteView("sync_pop")
		_, err := g.SetCurrentView("order_panel")
		return err
	}
	return nil
}

// ClearLogs resets the log panel
func (m *Manager) ClearLogs(g *gocui.Gui, v *gocui.View) error {
	m.Logger.Clear()
	return nil
}

func (m *Manager) ToggleAPIPopup(g *gocui.Gui, v *gocui.View) error {
	m.mu.Lock()
	m.ShowAPI = !m.ShowAPI
	m.APIPopup.FocusedField = 0 // Reset focus to first field
	m.mu.Unlock()

	if !m.ShowAPI {
		g.Cursor = false
		_, err := g.SetCurrentView("order_panel")
		return err
	}
	return nil
}

// NextApiField cycles through the 3 input boxes

func (m *Manager) NextAPIField(g *gocui.Gui, v *gocui.View) error {
	m.APIPopup.FocusedField = (m.APIPopup.FocusedField + 1) % 3
	return nil
}

// SaveApiConfig takes the data and pushes it to the SDK

// Helper struct to parse Bitget response
type bitgetAccountInfo struct {
	Code string `json:"code"`
	Data struct {
		UserId string `json:"userId"`
	} `json:"data"`
}

func (m *Manager) SaveAPIConfig(g *gocui.Gui, v *gocui.View) error {
	vKey, _ := g.View("api_key")
	vSec, _ := g.View("api_secret")
	vPas, _ := g.View("api_pass")

	key := strings.TrimSpace(vKey.Buffer())
	sec := strings.TrimSpace(vSec.Buffer())
	pas := strings.TrimSpace(vPas.Buffer())

	if key == "" || sec == "" || pas == "" {
		m.Logger.Error("API Setup: Missing Fields")
		return nil
	}

	// 1. Temporarily lock UI for validation
	m.mu.Lock()
	m.APIPopup.Validating = true
	m.mu.Unlock()

	// 2. Inject credentials into SDK memory
	config.APIKey = key
	config.SecretKey = sec
	config.PASSPHRASE = pas

	// 3. Test viability (Run in separate goroutine to avoid freezing UI)

	go func() {
		client := new(v2.SpotAccountClient).Init()
		resp, err := client.Info()

		g.Update(func(g *gocui.Gui) error {
			m.mu.Lock()
			m.APIPopup.Validating = false
			m.mu.Unlock()

			if err != nil {
				m.Logger.Error(fmt.Sprintf("Network Error: %v", err))
				return nil
			}

			var info bitgetAccountInfo
			_ = json.Unmarshal([]byte(resp), &info)

			if info.Code == "00000" {
				// SUCCESS: Update the UserID in the Manager
				m.mu.Lock()
				m.UserID = info.Data.UserId
				m.ShowAPI = false
				m.mu.Unlock()

				// Save to file
				_ = SaveSession(SavedConfig{
					APIKey:     key,
					SecretKey:  sec,
					Passphrase: pas,
				})

				m.Logger.Info(fmt.Sprintf("Bitget API: Connected (User ID: %s)", info.Data.UserId))

				m.mu.Lock()
				m.ShowAPI = false
				m.mu.Unlock()

				_ = g.DeleteView("api_pop")
				_ = g.DeleteView("api_key")
				_ = g.DeleteView("api_secret")
				_ = g.DeleteView("api_pass")

				g.Cursor = false
				_, _ = g.SetCurrentView("order_panel")
			} else {
				m.handleAPIError(info.Code, "")
			}
			return nil
		})
	}()

	return nil
}

// Maps Bitget error codes to human readable logs
func (m *Manager) handleAPIError(code, msg string) {
	switch code {
	case "40006", "40012", "40037":
		m.Logger.Error("Bitget Error: Invalid Key, Secret, or Passphrase")
	case "40038":
		m.Logger.Error("Bitget Error: IP Address not whitelisted in Bitget Settings")
	case "40008":
		m.Logger.Error("Bitget Error: System clock out of sync. Please sync your Windows/Linux time.")
	case "40014":
		m.Logger.Error("Bitget Error: Key is valid but lacks 'Spot/Futures' permissions")
	default:
		m.Logger.Error(fmt.Sprintf("Bitget Error [%s]: %s", code, msg))
	}
}

func (m *Manager) RefreshUserInfo(g *gocui.Gui) {
	if config.APIKey == "" {
		return
	}

	go func() {
		client := new(v2.SpotAccountClient).Init()
		resp, _ := client.Info()
		var info bitgetAccountInfo
		_ = json.Unmarshal([]byte(resp), &info)

		if info.Code == "00000" {
			g.Update(func(g *gocui.Gui) error {
				m.mu.Lock()
				m.UserID = info.Data.UserId
				m.mu.Unlock()
				return nil
			})
		}
	}()
}
