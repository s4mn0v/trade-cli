package ui

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/awesome-gocui/gocui"
)

type Timeframe struct {
	Label string
	Secs  int
	Color string
}

type SyncPopup struct {
	Timeframes []Timeframe
}

func NewSyncPopup() *SyncPopup {
	return &SyncPopup{
		Timeframes: []Timeframe{
			{"1M", 60, "\033[36m"},
			{"5M", 300, "\033[36m"},
			{"15M", 900, "\033[36m"},
			{"30M", 1800, "\033[36m"},
			{"1H", 3600, "\033[34m"},
			{"4H", 14400, "\033[32m"},
			{"1D", 86400, "\033[35m"},
		},
	}
}

func (s *SyncPopup) Render(g *gocui.Gui, maxX, maxY int) error {
	w, h := 42, 16 // Adjusted width for better fit
	x0, y0 := maxX/2-w/2, maxY/2-h/2
	x1, y1 := maxX/2+w/2, maxY/2+h/2

	v, err := g.SetView("sync_pop", x0, y0, x1, y1, 0)
	if err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = " Sync & Sessions "
		v.FrameColor = gocui.ColorCyan
		_, _ = g.SetCurrentView("sync_pop")
	}

	v.Clear()
	now := time.Now().Unix()

	_, _ = fmt.Fprintln(v, "\033[1m  TIMEFRAME      REMAINING\033[0m")
	for _, tf := range s.Timeframes {
		rem := int64(tf.Secs) - (now % int64(tf.Secs))

		var timeStr string
		if rem >= 3600 {
			timeStr = fmt.Sprintf("%02dh %02dm", rem/3600, (rem%3600)/60)
		} else {
			timeStr = fmt.Sprintf("%02dm %02ds", rem/60, rem%60)
		}

		labelPadded := fmt.Sprintf("%-10s", tf.Label)
		_, _ = fmt.Fprintf(v, "  %s%s\033[0m   %s\n", tf.Color, labelPadded, timeStr)
	}

	_, _ = fmt.Fprintln(v, "\n "+strings.Repeat("─", w-4))
	s.renderSessions(v)

	return nil
}

func (s *SyncPopup) renderSessions(v *gocui.View) {
	nyTime := time.Now().UTC().Add(-4 * time.Hour)
	curM := nyTime.Hour()*60 + nyTime.Minute()

	sessions := []struct {
		Name  string
		Start int
		End   int
		Color string
	}{
		{"ASIA", 20 * 60, 3 * 60, "\033[36m"},     // Cyan
		{"LONDON", 3 * 60, 8*60 + 30, "\033[33m"}, // Yellow
		{"NY", 8*60 + 30, 16 * 60, "\033[32m"},    // Green
	}

	var activeName, activeColor string
	var remSecs int

	for _, sess := range sessions {
		isOpen := false
		if sess.Start < sess.End {
			isOpen = curM >= sess.Start && curM < sess.End
		} else {
			isOpen = curM >= sess.Start || curM < sess.End
		}

		if isOpen {
			activeName = sess.Name
			activeColor = sess.Color
			remMins := (sess.End - curM + 1440) % 1440
			remSecs = (remMins * 60) - nyTime.Second()
			break
		}
	}

	if activeName != "" {
		h, m, s := remSecs/3600, (remSecs%3600)/60, remSecs%60
		// Aligned with the timeframe table above
		_, _ = fmt.Fprintf(v, "  ACTIVE SESSION: %s%s\033[0m\n", activeColor, activeName)
		_, _ = fmt.Fprintf(v, "  ENDS IN:        %02dh %02dm %02ds\n", h, m, s)
	} else {
		_, _ = fmt.Fprintln(v, "  ACTIVE SESSION: \033[90mNONE\033[0m")
	}
}
