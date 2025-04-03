package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/skillptm/Bolt/internal/app"
	"github.com/skillptm/Bolt/internal/logger"
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
	lg := &logger.Logger{}

	appInstance, err := app.NewApp(lg, images, icon)
	if err != nil {
		if len(lg.ErrorLogPath) > 0 {
			lg.Fatal("main: couldn't create app:\n--> %s", err.Error())
		} else {
			lg.Panic("main: couldn't create app:\n--> %s", err.Error())
		}
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
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: appInstance.Startup,
		Bind: []any{
			appInstance,
		},
	})

	if err != nil {
		lg.Fatal("main: wails.Run had an erroer while running:\n--> %s", err.Error())
	}
}
