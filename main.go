package main

// <---------------------------------------------------------------------------------------------------->

import (
	"embed"
	"fmt"

	"github.com/skillptm/Quick-Search/internal/app"
	"github.com/skillptm/Quick-Search/internal/appmenu"
	"github.com/skillptm/bws"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

// <---------------------------------------------------------------------------------------------------->

//go:embed all:frontend/dist
var assets embed.FS

//go:embed frontend/src/assets/images/*
var images embed.FS

// <---------------------------------------------------------------------------------------------------->

func main() {
	bws.ForceUpdateCache()

	// Create an instance of the app structure
	app, err := app.NewApp(images)
	if err != nil {
		fmt.Println("couldn't create new app: ", err.Error())
		return
	}

	// Create application with options
	err = wails.Run(&options.App{
		Title:             "Quick-Search",
		Width:             570,
		Height:            45,
		DisableResize:     true,
		Frameless:         true,
		HideWindowOnClose: true,
		AlwaysOnTop:       true,
		StartHidden:       true,
		Menu:              appmenu.Get(app),
		LogLevel:          logger.INFO,
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
