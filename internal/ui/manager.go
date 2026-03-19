// Package ui Manager
package ui

import (
	"errors"
	"fmt"
	"sync"

	"github.com/awesome-gocui/gocui"
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
	FuturesLeverage int
	PositionPercent int

	SpotBalance    float64
	FuturesBalance float64

	History       *HistoryTable
	Logger        *UILogger
	LeveragePopup *LeveragePopup
	QuantityPopup *QuantityPopup
}

func NewManager() *Manager {
	return &Manager{
		History:         NewHistoryTable(),
		Logger:          NewUILogger(),
		Mode:            ModeSpot,
		LeveragePopup:   NewLeveragePopup(),
		QuantityPopup:   NewQuantityPopup(),
		ShowLeverage:    false,
		ShowQuantity:    false,
		FuturesLeverage: 5,
		PositionPercent: 100,
		SpotBalance:     1250.50,
		FuturesBalance:  500.00,
	}
}

func (m *Manager) Layout(g *gocui.Gui) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	maxX, maxY := g.Size()
	orderH := 5

	histW := int(float64(maxX) * 0.70)
	logX0 := histW + 1

	// 1. Order Panel
	if v, err := g.SetView("order_panel", 0, 0, maxX-1, orderH, 0); err == nil {
		var title string
		if m.Mode == ModeSpot {
			title = fmt.Sprintf(" Place order [%s] [%d%%] [Avbl: %.2f USDT] ", m.Mode, m.PositionPercent, m.SpotBalance)
		} else {
			title = fmt.Sprintf(" Place order [%s] [%dx] [%d%%] [Avbl: %.2f USDT] ", m.Mode, m.FuturesLeverage, m.PositionPercent, m.FuturesBalance)
		}
		v.Title = title
		v.Clear()
		if m.Mode == ModeSpot {
			_, _ = fmt.Fprint(v, "\n  (Ctrl+O, b) = Buy | (Ctrl+O, s) = Sell | (Ctrl+S) Spot | (Ctrl+F) Futures")
		} else {
			_, _ = fmt.Fprint(v, "\n  (Ctrl+O, l) = Long | (Ctrl+O, s) = Short | (L) Leverage | (Ctrl+S) Spot | (Ctrl+F) Futures")
		}
	}

	// 2. History Panel
	if v, err := g.SetView("history", 0, orderH+1, histW, maxY-1, 0); err == nil || errors.Is(err, gocui.ErrUnknownView) {
		v.Subtitle = " History "
		m.History.Render(v, histW, m.Mode)
	}

	// 3. Logs Panel
	if v, err := g.SetView("logs", logX0, orderH+1, maxX-1, maxY-1, 0); err == nil || errors.Is(err, gocui.ErrUnknownView) {
		v.Title = " Logs "
		v.Autoscroll = true
		v.Wrap = true // Your requested wrap
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
