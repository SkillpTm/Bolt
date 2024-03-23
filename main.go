package main

// <---------------------------------------------------------------------------------------------------->

import (
	"embed"

	"github.com/skillptm/bws"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/skillptm/Quick-Search/internal/app"
	"github.com/skillptm/Quick-Search/internal/appmenu"
)

// <---------------------------------------------------------------------------------------------------->

var assets embed.FS

// <---------------------------------------------------------------------------------------------------->

func main() {
	bws.ForceUpdateCache()

	// Create an instance of the app structure
	app := app.NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:             "Quick-Search",
		Width:             570,
		Height:            45,
		DisableResize:     true,
		Frameless:         true,
		HideWindowOnClose: true,
		AlwaysOnTop:       true,
		StartHidden:       true,
		Menu:              appmenu.Get(app),
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.Startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
