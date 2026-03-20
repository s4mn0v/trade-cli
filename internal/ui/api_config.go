package ui

import (
	"errors"

	"github.com/awesome-gocui/gocui"
)

type APIConfigPopup struct {
	FocusedField int // 0: Key, 1: Secret, 2: Passphrase
	Validating   bool
}

func (a *APIConfigPopup) Render(g *gocui.Gui, maxX, maxY int) error {
	w, h := 64, 11
	x0, y0 := maxX/2-w/2, maxY/2-h/2
	x1, y1 := maxX/2+w/2, maxY/2+h/2

	// 1. Background Container
	if v, err := g.SetView("api_pop", x0, y0, x1, y1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " Bitget API Configuration "
		v.FrameColor = gocui.ColorMagenta
	}

	// Change color based on validation status
	color := gocui.ColorMagenta
	if a.Validating {
		color = gocui.ColorYellow
	}

	vMain, _ := g.View("api_pop")
	vMain.FrameColor = color

	// 2. Define the three input fields
	fieldNames := []string{"api_key", "api_secret", "api_pass"}
	labels := []string{" API KEY ", " SECRET  ", " PASS    "}

	for i, name := range fieldNames {
		// Calculate position for each box
		fy0 := y0 + 2 + (i * 2)
		if v, err := g.SetView(name, x0+12, fy0, x1-2, fy0+2, 0); err != nil {
			if !errors.Is(err, gocui.ErrUnknownView) {
				return err
			}
			v.Title = labels[i]
			v.Editable = true
			v.Editor = gocui.DefaultEditor
		}
	}

	if a.Validating {
		vMain.Subtitle = " Validating credentials with Bitget... "
	} else {
		vMain.Subtitle = " [Tab] Switch | [Enter] Save & Test | [Esc] Cancel "
	}

	// Set focus and cursor to the current field
	activeField := fieldNames[a.FocusedField]
	_, _ = g.SetCurrentView(activeField)
	g.Cursor = !a.Validating

	return nil
}
