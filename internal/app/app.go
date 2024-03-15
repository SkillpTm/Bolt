// Package app ...
package app

// <---------------------------------------------------------------------------------------------------->

import (
	"context"

	"github.com/skillptm/ModSearch/internal/searchhandler"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// <---------------------------------------------------------------------------------------------------->

// App holds all the main data and functions relevant to the front- and backend.
type App struct {
	CTX           context.Context
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
}

// LaunchSearch starts a search on the SearchHandler of the app
func (a *App) LaunchSearch(input string) {
	if len(input) < 1 {
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
