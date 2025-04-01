package main

import (
	"embed"
	"fmt"
	"log"

	"fyne.io/systray"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/skillptm/Bolt/internal/app"
)

var (
	//go:embed frontend/dist/*
	assets embed.FS
	//go:embed frontend/src/assets/images/*
	images embed.FS
	//go:embed build/appicon.png
	icon embed.FS
)

func main() {
	appInstance, err := app.NewApp(images)
	if err != nil {
		log.Fatalf("main: couldn't create app:\n--> %s", err.Error())
	}

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
				runtime.WindowShow(appInstance.CTX)
			case <-quit.ClickedCh:
				systray.Quit()
				runtime.Quit(appInstance.CTX)
			}
		}
	}

	go systray.Run(onReady, func() {})

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
