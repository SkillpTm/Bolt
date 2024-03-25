// Package appmenu ...
package appmenu

// <---------------------------------------------------------------------------------------------------->

import (
	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/skillptm/Quick-Search/internal/app"
)

// <---------------------------------------------------------------------------------------------------->

var AppIcon []byte
var AppInstance *app.App

// <---------------------------------------------------------------------------------------------------->

// Get provides our default menu
func Get() *menu.Menu {
	appMenu := menu.NewMenu()
	subMenu := appMenu.AddSubmenu("Shortcuts")
	subMenu.AddText("Hide", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		runtime.WindowHide(AppInstance.CTX)
	})
	subMenu.AddText("Page forwards", keys.CmdOrCtrl("f"), func(_ *menu.CallbackData) {
		runtime.EventsEmit(AppInstance.CTX, "pageForward")
	})
	subMenu.AddText("Page backwards", keys.CmdOrCtrl("b"), func(_ *menu.CallbackData) {
		runtime.EventsEmit(AppInstance.CTX, "pageBackward")
	})

	return appMenu
}

func OnReady() {
	systray.SetIcon(AppIcon)
	systray.SetTooltip("Quick Search")
	open := systray.AddMenuItem("Open", "opens quick search bar")
	quit := systray.AddMenuItem("Quit", "properly quits quick search")

	// OnReady gets started as a goroutine, so we can have an infinite for loop here
	for {
		select {
		case <-open.ClickedCh:
			runtime.WindowShow(AppInstance.CTX)
		case <-quit.ClickedCh:
			systray.Quit()
			runtime.Quit(AppInstance.CTX)
		}
	}
}
