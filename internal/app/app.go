// Package app ...
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

	"github.com/skillptm/Bolt/internal/modules"
)

// App holds all the main data and functions relevant to the front- and backend.
type App struct {
	CTX           context.Context
	images        embed.FS
	SearchHandler *modules.SearchHandler
}

// NewApp is the constructor for App
func NewApp(images embed.FS) (*App, error) {
	sh, err := modules.NewSearchHandler()
	if err != nil {
		return nil, fmt.Errorf("NewApp: couldn't create SearchHandler:\n--> %w", err)
	}

	return &App{
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
	openHotkey := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyS)

	err := openHotkey.Register()
	if err != nil {
		log.Fatalf("openOnHotKey: couldn't register main hotkey:\n--> %s", err.Error())
	}

	for range openHotkey.Keydown() {
		runtime.WindowShow(a.CTX)
		runtime.WindowReload(a.CTX)
	}
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
