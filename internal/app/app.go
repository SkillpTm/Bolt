// Package app holds the wails app and all emit aswell as export functions that can be used in TS
package app

import (
	"context"
	"embed"
	"encoding/base64"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.design/x/hotkey"

	"github.com/skillptm/Bolt/internal/config"
	"github.com/skillptm/Bolt/internal/modules"
)

// App holds all the main data and functions relevant to the front- and backend.
type App struct {
	conf          *config.Config
	CTX           context.Context
	hotkey        hotkey.Key
	icon          embed.FS
	images        embed.FS
	SearchHandler *modules.SearchHandler
}

// NewApp is the constructor for App
func NewApp(images embed.FS, icon embed.FS) (*App, error) {
	conf, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("NewApp: couldn't create config:\n--> %w", err)
	}

	sh, err := modules.NewSearchHandler(conf)
	if err != nil {
		return nil, fmt.Errorf("NewApp: couldn't create SearchHandler:\n--> %w", err)
	}

	keyMap := map[string]hotkey.Key{
		"a": hotkey.KeyA, "b": hotkey.KeyB, "c": hotkey.KeyC, "d": hotkey.KeyD, "e": hotkey.KeyE,
		"f": hotkey.KeyF, "g": hotkey.KeyG, "h": hotkey.KeyH, "i": hotkey.KeyI, "j": hotkey.KeyJ,
		"k": hotkey.KeyK, "l": hotkey.KeyL, "m": hotkey.KeyM, "n": hotkey.KeyN, "o": hotkey.KeyO,
		"p": hotkey.KeyP, "q": hotkey.KeyQ, "r": hotkey.KeyR, "s": hotkey.KeyS, "t": hotkey.KeyT,
		"u": hotkey.KeyU, "v": hotkey.KeyV, "w": hotkey.KeyW, "x": hotkey.KeyX, "y": hotkey.KeyY,
		"z": hotkey.KeyZ, " ": hotkey.KeySpace, "space": hotkey.KeySpace,
	}

	if _, ok := keyMap[strings.ToLower(conf.ShortCutEnd)]; !ok {
		return nil, fmt.Errorf("NewApp: invalid hotkey input")
	}

	return &App{
		conf:          conf,
		icon:          icon,
		hotkey:        keyMap[strings.ToLower(conf.ShortCutEnd)],
		images:        images,
		SearchHandler: sh,
	}, nil
}

/*
Startup is called when the app starts. The context gets saved on the app.
It's responsible for launching the main loop, containing all the emit functions.
*/
func (a *App) Startup(CTX context.Context) {
	a.CTX = CTX
	go setupTray(a, a.icon)
	go a.emitSearchResult()
	go a.openOnHotKey()
}

// GetImageData emits a map[name]base64 png data to the frotend to bind in the images
func (a *App) GetImageData() map[string]string {
	imageData := map[string]string{
		"cross":            "frontend/src/assets/images/cross.png",
		"google":           "frontend/src/assets/images/google.png",
		"file":             "frontend/src/assets/images/file.png",
		"folder":           "frontend/src/assets/images/folder.png",
		"left":             "frontend/src/assets/images/left.png",
		"magnifying_glass": "frontend/src/assets/images/magnifying_glass.png",
		"not-left":         "frontend/src/assets/images/not_left.png",
		"right":            "frontend/src/assets/images/right.png",
		"not-right":        "frontend/src/assets/images/not_right.png",
		"tick":             "frontend/src/assets/images/tick.png",
	}

	for name, path := range imageData {
		imageBytes, err := a.images.ReadFile(path)
		if err != nil {
			return map[string]string{}
		}

		imageData[name] = fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(imageBytes))
	}

	return imageData
}

// emitSearchResult runs continuously and emits the search results with the "searchResult" event to the frontend
func (a *App) emitSearchResult() {
	for results := range a.SearchHandler.ResultsChan {
		runtime.EventsEmit(a.CTX, "searchResult", results)
	}
}

// openOnHotKey will unhide and reload the app when ctrl+shift+s is pressed
func (a *App) openOnHotKey() {
	openHotkey := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, a.hotkey)

	err := openHotkey.Register()
	if err != nil {
		log.Fatalf("openOnHotKey: couldn't register main hotkey:\n--> %s", err.Error())
	}

	for range openHotkey.Keydown() {
		a.ShowWindow()
	}
}

// HideWindow is a wrapper around runtime.WindowHide that ensures our cache data doesn't unnecessarily stay in memory
func (a *App) HideWindow() {
	runtime.WindowHide(a.CTX)
	a.SearchHandler.ClearImportedCache()
}

// LaunchSearch starts a search on the SearchHandler of the app
func (a *App) LaunchSearch(input string) {
	if len(input) < 1 {
		a.SearchHandler.ResultsChan <- []string{}
		return
	}

	go a.SearchHandler.Search(input)
}

// OpenFileExplorer allows you to open the file manager at any entry's location and select it (if the file manager is dolphin or nautilus)
func (a *App) OpenFileExplorer(filePath string) {
	var cmd *exec.Cmd

	if _, err := exec.LookPath("dolphin"); err == nil {
		cmd = exec.Command("dolphin", "--select", filePath)
	} else if _, err := exec.LookPath("nautilus"); err == nil {
		cmd = exec.Command("nautilus", "--select", filePath)
	} else {
		if index := strings.LastIndex(filePath, "/"); index != -1 {
			filePath = filePath[:index+1]
		}

		cmd = exec.Command("xdg-open", filePath)
	}

	err := cmd.Start()
	if err != nil {
		log.Fatalf("OpenFileExplorer: couldn't file manager:\n--> %s", err.Error())
	}
}

// ShowWindow is a wrapper around runtime.WindowShow that ensures we load our cache data into memory
func (a *App) ShowWindow() {
	a.SearchHandler.ImportCache()
	runtime.WindowShow(a.CTX)
}
