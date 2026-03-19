package ui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/awesome-gocui/gocui"
)

type QuantityPopup struct {
	CurrentVal int // 0 to 100
}

func NewQuantityPopup() *QuantityPopup {
	return &QuantityPopup{CurrentVal: 100} // Default 10%
}

func (q *QuantityPopup) Render(g *gocui.Gui, maxX, maxY int, balance float64, mode string) error {
	w, h := 45, 9
	x0, y0 := maxX/2-w/2, maxY/2-h/2
	x1, y1 := maxX/2+w/2, maxY/2+h/2

	v, err := g.SetView("quantity_pop", x0, y0, x1, y1, 0)
	if err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " Position Size % "
		v.FrameColor = gocui.ColorYellow
		v.Overlaps = 0
		v.FrameRunes = []rune{'═', '║', '╔', '╗', '╚', '╝', '╠', '╣', '╦', '╩', '╬'}
		_, _ = g.SetCurrentView("quantity_pop")
	}

	v.Clear()
	// Calculate the actual USDT amount based on percentage
	amount := balance * (float64(q.CurrentVal) / 100.0)

	_, _ = fmt.Fprintf(v, "\n    Mode: %s | Avbl: %.2f\n", mode, balance)
	_, _ = fmt.Fprintf(v, "    Set Size: \033[1;33m%d%%\033[0m (\033[32m%.2f USDT\033[0m)\n\n", q.CurrentVal, amount)

	barWidth := 35
	filled := int(float64(q.CurrentVal) / 100.0 * float64(barWidth))
	filled = max(filled, 0)

	bar := "\033[7m" + strings.Repeat(" ", filled) + "\033[0m" + strings.Repeat("░", barWidth-filled)
	_, _ = fmt.Fprintf(v, "    %s\n\n", bar)
	_, _ = fmt.Fprintf(v, "    \033[33m[k/j]\033[0m Adjust | \033[33m[r]\033[0m Reset | \033[33m[Enter]\033[0m Set")

	return nil
}
