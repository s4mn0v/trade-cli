// Package ui Manager
package ui

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/awesome-gocui/gocui"
	"github.com/s4mn0v/bitget/config"
)

const (
	ModeSpot    = "SPOT"
	ModeFutures = "FUTURES"
)

type Manager struct {
	mu              sync.RWMutex
	OrderMode       bool
	Mode            string
	ShowLeverage    bool
	ShowQuantity    bool
	ShowCoin        bool
	ShowSync        bool
	ShowAPI         bool
	UserID          string
	CurrentCoin     string
	FuturesLeverage int
	PositionPercent int

	SpotBalance    float64
	SpotAssets     map[string]float64
	FuturesBalance float64

	History       *HistoryTable
	Logger        *UILogger
	LeveragePopup *LeveragePopup
	QuantityPopup *QuantityPopup
	APIPopup      *APIConfigPopup
	CoinPopup     *CoinPopup
	SyncPopup     *SyncPopup
	Positions     *PositionList
}

func NewManager() *Manager {
	return &Manager{
		History:         NewHistoryTable(),
		Logger:          NewUILogger(),
		Mode:            ModeSpot,
		LeveragePopup:   NewLeveragePopup(),
		QuantityPopup:   NewQuantityPopup(),
		CoinPopup:       NewCoinPopup(),
		SyncPopup:       NewSyncPopup(),
		Positions:       NewPositionList(),
		ShowLeverage:    false,
		ShowQuantity:    false,
		ShowCoin:        false,
		ShowSync:        false,
		ShowAPI:         false,
		UserID:          "",
		APIPopup:        &APIConfigPopup{FocusedField: 0},
		FuturesLeverage: 5,
		PositionPercent: 100,
		SpotBalance:     1250.50,
		SpotAssets: map[string]float64{
			"BTC": 0.052,
			"ETH": 1.2,
			"SOL": 15.0,
		},
		FuturesBalance: 500.00,
		CurrentCoin:    "BTCUSDT",
	}
}

func (m *Manager) Layout(g *gocui.Gui) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	maxX, maxY := g.Size()
	orderH := 10

	histW := int(float64(maxX) * 0.70)
	logX0 := histW + 1

	// 1. Calculate API Status (READY / OFFLINE)
	apiDisplay := "OFFLINE"
	if m.UserID != "" {
		apiDisplay = "READY: " + m.UserID
	} else if config.APIKey != "" {
		apiDisplay = "Connecting..."
	}
	// Combining into a styled string
	// apiDisplay := fmt.Sprintf("API: %s", m.UserID)

	// Helper to get the base asset
	baseAsset := strings.TrimSuffix(m.CurrentCoin, "USDT")
	assetAmount := m.SpotAssets[baseAsset]

	// 2. Order Panel (Clean title without Status)
	if v, err := g.SetView("order_panel", 0, 0, maxX-1, orderH, 0); err == nil {
		var title string
		if m.Mode == ModeSpot {
			title = fmt.Sprintf(" %s SPOT | Size: %d%% | Avbl: %.2f USDT | %.4f %s ", m.CurrentCoin, m.PositionPercent, m.SpotBalance, assetAmount, baseAsset)
		} else {
			title = fmt.Sprintf(" %s FUTURES | %dx | Size: %d%% | Avbl: %.2f USDT ", m.CurrentCoin, m.FuturesLeverage, m.PositionPercent, m.FuturesBalance)
		}
		v.Title = title
		v.Clear()

		if m.Mode == ModeSpot {
			fmt.Fprintf(v, "\n  \033[32m(b) BUY %s\033[0m (Spend USDT) | \033[31m(s) SELL %s\033[0m (Spend %s)", baseAsset, baseAsset, baseAsset)
			fmt.Fprint(v, "\n  (Ctrl+O) Order Mode | (p) Change Coin | (Q) Set Size")
			fmt.Fprint(v, "\n\n  \033[90mSpot mode enabled. View history below.\033[0m")
		} else {
			fmt.Fprint(v, "\n  \033[32m(l) LONG\033[0m | \033[31m(s) SHORT\033[0m | \033[33m(c) CLOSE SELECTED\033[0m")
			fmt.Fprint(v, "\n  (L) Leverage | (Q) Quantity | (p) Coin | (Tab) Switch Focus")
			m.Positions.Render(v, maxX)
		}
	}

	// 3. History Panel (Adding the Status here)
	if v, err := g.SetView("history", 0, orderH+1, histW, maxY-1, 0); err == nil || errors.Is(err, gocui.ErrUnknownView) {
		// We use a combination of " History " and the API status
		v.Subtitle = fmt.Sprintf(" History [%s] ", apiDisplay)
		m.History.Render(v, histW, m.Mode)
	}

	// 4. Logs Panel (Stays the same)
	if v, err := g.SetView("logs", logX0, orderH+1, maxX-1, maxY-1, 0); err == nil || errors.Is(err, gocui.ErrUnknownView) {
		v.Title = " Logs "
		v.Autoscroll = true
		v.Wrap = true
		m.Logger.Render(v)
	}

	// 3. Logs Panel
	if v, err := g.SetView("logs", logX0, orderH+1, maxX-1, maxY-1, 0); err == nil || errors.Is(err, gocui.ErrUnknownView) {
		v.Title = " Logs "
		v.Autoscroll = true
		v.Wrap = true
		m.Logger.Render(v)
	}

	// --- QUANTITY POPUP LAYER ---
	if m.ShowQuantity {
		balance := m.SpotBalance
		if m.Mode == ModeFutures {
			balance = m.FuturesBalance
		}
		if err := m.QuantityPopup.Render(g, maxX, maxY, balance, m.Mode); err != nil {
			return err
		}
	} else {
		_ = g.DeleteView("quantity_pop")
	}

	// --- LEVERAGEPOPUP LAYER ---
	if m.ShowLeverage {
		if err := m.LeveragePopup.Render(g, maxX, maxY); err != nil {
			return err
		}
	} else {
		if err := g.DeleteView("leverage_pop"); err != nil && !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
	}

	// --- COIN POPUP LAYER ---
	if m.ShowCoin {
		input := ""
		if v, err := g.View("coin_pop"); err == nil {
			input = v.Buffer()
		}
		if err := m.CoinPopup.Render(g, maxX, maxY, input); err != nil {
			return err
		}
		g.Cursor = true
	} else {
		_ = g.DeleteView("coin_pop")
		if !m.ShowQuantity && !m.ShowLeverage {
			g.Cursor = false
		}
	}

	// --- SYNC POPUP LAYER ---
	if m.ShowSync {
		if err := m.SyncPopup.Render(g, maxX, maxY); err != nil {
			return err
		}
	} else {
		_ = g.DeleteView("sync_pop")
	}

	// --- API CONFIG POPUP ---
	if m.ShowAPI {
		if err := m.APIPopup.Render(g, maxX, maxY); err != nil {
			return err
		}
	} else {
		// Clean up all 4 views when closed
		_ = g.DeleteView("api_pop")
		_ = g.DeleteView("api_key")
		_ = g.DeleteView("api_secret")
		_ = g.DeleteView("api_pass")
	}

	m.applyDynamicStyles(g)
	return nil
}

func (m *Manager) applyDynamicStyles(g *gocui.Gui) {
	curr := ""
	if v := g.CurrentView(); v != nil {
		curr = v.Name()
	}

	// Default Mode Colors
	modeColor := gocui.ColorGreen // SPOT
	if m.Mode == ModeFutures {
		modeColor = gocui.ColorRed // FUTURES
	}

	selectedRunes := []rune{'═', '║', '╔', '╗', '╚', '╝', '╠', '╣', '╦', '╩', '╬'}

	views := []string{"order_panel", "history", "logs"}
	for _, name := range views {
		if v, err := g.View(name); err == nil {
			if (name == "order_panel" && m.OrderMode) || curr == name {
				v.FrameColor = gocui.ColorYellow
				v.FrameRunes = selectedRunes
			} else {
				v.FrameColor = modeColor
				v.FrameRunes = nil
			}
		}
	}
}

func SetTerminalSize(rows, cols int) {
	fmt.Printf("\033[8;%d;%dt", rows, cols)
}

// AnyPopupOpen returns true if any modal window is currently being displayed
func (m *Manager) AnyPopupOpen() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.ShowLeverage || m.ShowQuantity || m.ShowCoin || m.ShowSync
}
