package main

import (
	"log"

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

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
