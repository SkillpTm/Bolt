// Package app ...
package app

import (
	"embed"
	"fmt"
	"log"

	"fyne.io/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// setupTray sets up the tray, with the icon, open and quit options and blocks indefenitly
func setupTray(a *App, icon embed.FS) {
	onReady := func() {
		appIcon, err := icon.ReadFile("build/appicon.png")
		if err != nil {
			log.Fatal(fmt.Errorf("main: couldn't get image build/appicon.png from embed:\n--> %w", err))
		}

		systray.SetIcon(appIcon)
		systray.SetTooltip("Bolt")
		open := systray.AddMenuItem("Open", "opens bolt search")
		quit := systray.AddMenuItem("Quit", "quits bolt search")

		for {
			select {
			case <-open.ClickedCh:
				runtime.WindowShow(a.CTX)
			case <-quit.ClickedCh:
				systray.Quit()
				runtime.Quit(a.CTX)
			}
		}
	}

	systray.Run(onReady, func() {})
}
