package main

// <---------------------------------------------------------------------------------------------------->

import (
	"embed"

	"github.com/skillptm/bws"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/skillptm/ModSearch/internal/app"
	"github.com/skillptm/ModSearch/internal/appmenu"
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
		Title:             "ModSearch",
		Width:             571,
		Height:            46,
		DisableResize:     true,
		Frameless:         true,
		HideWindowOnClose: true,
		AlwaysOnTop:       true,
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
