package ui

import (
	"errors"
	"slices"
	"strings"

	"github.com/awesome-gocui/gocui"
)

type CoinPopup struct {
	Suggestions []string
}

func NewCoinPopup() *CoinPopup {
	return &CoinPopup{
		Suggestions: []string{"BTCUSDT", "ETHUSDT", "SOLUSDT", "BNBUSDT", "LINKUSDT", "DOTUSDT", "MATICUSDT", "XRPUSDT"},
	}
}

func (c *CoinPopup) Render(g *gocui.Gui, maxX, maxY int, currentInput string) error {
	w, h := 40, 4
	x0, y0 := maxX/2-w/2, maxY/2-h/2
	x1, y1 := maxX/2+w/2, maxY/2+h/2

	v, err := g.SetView("coin_pop", x0, y0, x1, y1, 0)
	if err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " Enter Coin (e.g. BTCUSDT) "
		v.FrameColor = gocui.ColorCyan
		v.Editable = true // Enable text input
		v.Editor = gocui.DefaultEditor

		_, _ = g.SetCurrentView("coin_pop")
	}

	// Displaying suggestions based on current input
	matches := []string{}
	upperInput := strings.ToUpper(strings.TrimSpace(currentInput))
	if upperInput != "" {
		for _, s := range c.Suggestions {
			if strings.Contains(s, upperInput) {
				matches = append(matches, s)
			}
		}
	}

	// Display matches in the subtitle for "autocomplete" visibility
	if len(matches) > 0 {
		v.Subtitle = " Suggestions: " + strings.Join(matches, ", ") + " "
	} else {
		v.Subtitle = " Type to search... "
	}

	return nil
}

// IsValid checks if the coin exists in our dataset
func (c *CoinPopup) IsValid(coin string) bool {
	coin = strings.ToUpper(strings.TrimSpace(coin))
	return slices.Contains(c.Suggestions, coin)
}

// GetMatches returns a list of coins that start with or contain the input string
func (c *CoinPopup) GetMatches(input string) []string {
	upperInput := strings.ToUpper(strings.TrimSpace(input))
	if upperInput == "" {
		return nil
	}

	matches := []string{}
	for _, s := range c.Suggestions {
		if strings.HasPrefix(s, upperInput) {
			matches = append(matches, s)
		}
	}
	return matches
}
