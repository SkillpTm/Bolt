// Package app holds the wails app and all emit aswell as export functions that can be used in TS
package app

import (
	"embed"
	"fmt"
	"log"
	"os"

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
		open := systray.AddMenuItem("Open", "opens Bolt search")
		quit := systray.AddMenuItem("Quit", "quits Bolt search")

		for {
			select {
			case <-open.ClickedCh:
				runtime.WindowShow(a.CTX)
			case <-quit.ClickedCh:
				systray.Quit()
				runtime.Quit(a.CTX)
				os.Exit(0)
			}
		}
	}

	systray.Run(onReady, func() {})
}
