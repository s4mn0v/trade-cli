package main

import (
	"log"
	"time"

	"trade-cli/internal/ui"

	"github.com/awesome-gocui/gocui"
)

func main() {
	ui.SetTerminalSize(20, 90)
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	m := ui.NewManager()
	g.SetManagerFunc(m.Layout)
	_ = m.InitKeybindings(g)

	go func() {
		for range time.Tick(time.Second) {
			// This forces gocui to run the Layout function
			// even if no keys are pressed
			g.Update(func(g *gocui.Gui) error { return nil })
		}
	}()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
