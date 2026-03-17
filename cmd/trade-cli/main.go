package main

import (
	"fmt"
	"log"

	"trade-cli/internal/ui"

	"github.com/awesome-gocui/gocui"
)

// SetTerminalSize attempts to resize the physical terminal window
// Height is in rows, Width is in columns
func SetTerminalSize(rows, cols int) {
	fmt.Printf("\033[8;%d;%dt", rows, cols)
}

func main() {
	ui.SetTerminalSize(20, 62)
	g, _ := gocui.NewGui(gocui.OutputNormal, true)
	defer g.Close()

	m := ui.NewManager()
	g.SetManagerFunc(m.Layout)

	// Call the new separated keybindings file
	if err := m.InitKeybindings(g); err != nil {
		log.Panicln(err)
	}

	g.MainLoop()
}
