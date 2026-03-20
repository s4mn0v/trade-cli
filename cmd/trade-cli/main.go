package main

import (
	"log"
	"time"

	"trade-cli/internal/ui"

	"github.com/awesome-gocui/gocui"
	"github.com/s4mn0v/bitget/config"
)

func main() {
	// 1. Try to load existing session from ~/.bitget-trade-cli.json
	saved, err := ui.LoadSession()
	sessionLoaded := false

	if err == nil && saved.APIKey != "" {
		// Inject into SDK Memory
		config.APIKey = saved.APIKey
		config.SecretKey = saved.SecretKey
		config.PASSPHRASE = saved.Passphrase
		sessionLoaded = true
	}

	// 2. Start GUI
	ui.SetTerminalSize(20, 90)
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	// 3. Initialize Manager
	m := ui.NewManager()
	g.SetManagerFunc(m.Layout)
	_ = m.InitKeybindings(g)

	// 4. If we loaded a session, send a log message to the UI
	if sessionLoaded {
		m.RefreshUserInfo(g)
		m.Logger.Info("Session: Credentials loaded from local file")
	} else {
		m.Logger.Warning("Session: No saved keys found. Press Ctrl+A to configure.")
	}

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
