package ui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/awesome-gocui/gocui"
)

type LeveragePopup struct {
	CurrentVal int
}

func NewLeveragePopup() *LeveragePopup {
	return &LeveragePopup{CurrentVal: 10}
}

func (l *LeveragePopup) Render(g *gocui.Gui, maxX, maxY int) error {
	w, h := 45, 9
	x0, y0 := maxX/2-w/2, maxY/2-h/2
	x1, y1 := maxX/2+w/2, maxY/2+h/2

	v, err := g.SetView("leverage_pop", x0, y0, x1, y1, 0)
	if err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		// Initialization settings
		v.Title = " Leverage Setting "
		v.FrameColor = gocui.ColorYellow
		v.Overlaps = 0
		v.FrameRunes = []rune{'═', '║', '╔', '╗', '╚', '╝', '╠', '╣', '╦', '╩', '╬'}
		g.SetCurrentView("leverage_pop")
	}

	// DRAWING LOGIC (Outside the error check to ensure it shows immediately)
	v.Clear()
	fmt.Fprintf(v, "\n    Adjust Leverage: \033[1;33m%dx\033[0m\n\n", l.CurrentVal)

	barWidth := 35
	filled := int(float64(l.CurrentVal) / 125.0 * float64(barWidth))
	if filled < 1 {
		filled = 1
	}

	// Progress bar using Invert for high visibility
	bar := "\033[7m" + strings.Repeat(" ", filled) + "\033[0m" + strings.Repeat("░", barWidth-filled)
	fmt.Fprintf(v, "    %s\n\n", bar)
	fmt.Fprintf(v, "    \033[33m[←/→]\033[0m Adjust | \033[33m[R]\033[0m Reset | \033[33m[Enter]\033[0m Set")

	return nil
}
