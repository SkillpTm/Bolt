// Package appmenu ...
package appmenu

// <---------------------------------------------------------------------------------------------------->

import (
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/skillptm/Quick-Search/internal/app"
)

// <---------------------------------------------------------------------------------------------------->

// Get provides our default menu
func Get(a *app.App) *menu.Menu {
	AppMenu := menu.NewMenu()
	FileMenu := AppMenu.AddSubmenu("Shortcuts")
	FileMenu.AddText("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		runtime.WindowHide(a.CTX)
	})

	return AppMenu
}
