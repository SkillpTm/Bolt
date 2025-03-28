package main

// <---------------------------------------------------------------------------------------------------->

import (
	"embed"
	"fmt"

	"github.com/getlantern/systray"
	"github.com/skillptm/bws"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/skillptm/Bolt/internal/app"
	"github.com/skillptm/Bolt/internal/appmenu"
)

// <---------------------------------------------------------------------------------------------------->

//go:embed frontend/dist/*
var assets embed.FS

//go:embed frontend/src/assets/images/*
var images embed.FS

//go:embed build/windows/icon.ico
var icon embed.FS

// <---------------------------------------------------------------------------------------------------->

func main() {
	bws.ForceUpdateCache()

	appInstance, err := app.NewApp(images)
	if err != nil {
		fmt.Println("couldn't create new app: ", err.Error())
		return
	}

	appmenu.AppIcon, err = icon.ReadFile("build/windows/icon.ico")
	if err != nil {
		fmt.Println("couldn't get image build/windows/icon.ico from embed: ", err.Error())
		return
	}

	appmenu.AppInstance = appInstance

	go systray.Run(appmenu.OnReady, func() {})

	err = wails.Run(&options.App{
		Title:             "Quick Search",
		Width:             570,
		Height:            45,
		DisableResize:     true,
		Frameless:         true,
		HideWindowOnClose: true,
		AlwaysOnTop:       true,
		StartHidden:       true,
		Menu:              appmenu.Get(),
		LogLevel:          logger.INFO,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: appInstance.Startup,
		Bind: []interface{}{
			appInstance,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
