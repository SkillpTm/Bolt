// Package app ...
package app

// <---------------------------------------------------------------------------------------------------->

import (
	"context"
	"log"
	"os/exec"
	"strings"

	"github.com/skillptm/Quick-Search/internal/searchhandler"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.design/x/hotkey"
)

// <---------------------------------------------------------------------------------------------------->

// App holds all the main data and functions relevant to the front- and backend.
type App struct {
	CTX           context.Context
	hotkey        *hotkey.Hotkey
	SearchHandler *searchhandler.SearchHandler
}

// NewApp creates a new App struct with all it's values.
func NewApp() *App {
	return &App{
		SearchHandler: searchhandler.New(),
	}
}

/*
Startup is called when the app starts. The context is saved so we can call the runtime methods.

It also starts the goroutine for emiting the search results to the frontend.
*/
func (a *App) Startup(CTX context.Context) {
	a.CTX = CTX
	go a.EmitSearchResult()
	go a.openOnHotKey()
}

// LaunchSearch starts a search on the SearchHandler of the app
func (a *App) LaunchSearch(input string) {
	if len(input) < 1 {
		a.SearchHandler.ResultsChan <- []string{}
		return
	}

	a.SearchHandler.StartSearch(input)
}

// EmitSearchResult runs continuously and emits the search results with the "searchResult" event to the frontend
func (a *App) EmitSearchResult() {
	for result := range a.SearchHandler.ResultsChan {
		runtime.EventsEmit(a.CTX, "searchResult", result)
	}
}

func (a *App) OpenFileExplorer(filepath string) {
	cmd := exec.Command("explorer", "/select,", strings.ReplaceAll(filepath, "/", "\\"))
	cmd.Run()
}

func (a *App) openOnHotKey() {
	a.hotkey = hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyS)
	err := a.hotkey.Register()
	if err != nil {
		log.Fatalf("main hotkey failed to register: %s", err)
		return
	}

	for range a.hotkey.Keydown() {
		runtime.WindowReload(a.CTX)
		runtime.WindowShow(a.CTX)
	}
}
