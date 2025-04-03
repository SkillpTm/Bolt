package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/skillptm/Bolt/internal/app"
)

var (
	//go:embed frontend/dist/*
	assets embed.FS
	//go:embed build/appicon.png
	icon embed.FS
	//go:embed frontend/src/assets/images/*
	images embed.FS
)

func main() {
	appInstance, err := app.NewApp(images, icon)
	if err != nil {
		log.Fatalf("main: couldn't create app:\n--> %s", err.Error())
	}

	err = wails.Run(&options.App{
		Title:             "Bolt",
		Width:             570,
		Height:            45,
		Frameless:         true,
		HideWindowOnClose: true,
		AlwaysOnTop:       true,
		StartHidden:       true,
		MinWidth:          570,
		MaxWidth:          570,
		MinHeight:         45,
		MaxHeight:         365,
		LogLevel:          logger.INFO,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: appInstance.Startup,
		Bind: []any{
			appInstance,
		},
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}
