// Package app ...
package app

// <---------------------------------------------------------------------------------------------------->

import (
	"context"
	"embed"
	"encoding/base64"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/skillptm/Quick-Search/internal/searchhandler"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.design/x/hotkey"
	"golang.org/x/sys/windows"
)

// <---------------------------------------------------------------------------------------------------->

// App holds all the main data and functions relevant to the front- and backend.
type App struct {
	CTX           context.Context
	hotkey        *hotkey.Hotkey
	images        map[string]string
	SearchHandler *searchhandler.SearchHandler
}

// NewApp creates a new App struct with all it's values.
func NewApp(images embed.FS) (*App, error) {
	imagePaths := map[string]string{
		"cross":  "frontend/src/assets/images/cross.png",
		"file":   "frontend/src/assets/images/file.png",
		"folder": "frontend/src/assets/images/folder.png",
		"tick":   "frontend/src/assets/images/tick.png",
	}

	imageMap := make(map[string]string)

	for name, path := range imagePaths {
		imageData, err := images.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("couldn't get image %s from embed: %s", path, err.Error())
		}

		imageMap[name] = "data:image/png;base64," + base64.StdEncoding.EncodeToString(imageData)
	}

	return &App{
		hotkey:        hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyS),
		images:        imageMap,
		SearchHandler: searchhandler.New(),
	}, nil
}

/*
Startup is called when the app starts. The context is saved so we can call the runtime methods.

It also starts the goroutine for emiting the search results to the frontend.
*/
func (a *App) Startup(CTX context.Context) {
	a.CTX = CTX
	go a.EmitSearchResult()
	go a.openOnHotKey()
	go a.windowHideOnUnselected()
}

// GetImageData receives a key (being the name of an image) and returns the base64 string data of that image
func (a *App) GetImageData(name string) string {
	if data, ok := a.images[name]; ok {
		return data
	}

	return ""
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

// OpenFileExplorer allows you to open the file explorer at any entry's location
func (a *App) OpenFileExplorer(filePath string) {
	cmd := exec.Command("explorer", "/select,", strings.TrimSuffix(strings.ReplaceAll(filePath, "/", "\\"), "\\"))
	cmd.Run()
}

// openOnHotKey will unhide and reload the app when ctrl+shift+s is pressed
func (a *App) openOnHotKey() {
	err := a.hotkey.Register()
	if err != nil {
		log.Fatalf("main hotkey failed to register: %s", err)
		return
	}

	for range a.hotkey.Keydown() {
		runtime.WindowShow(a.CTX)
	}
}

// windowHideOnUnselected will hide the window once you unselected it, by clicking somewhere else
func (a *App) windowHideOnUnselected() {
	recheckTicker := time.NewTicker(100 * time.Millisecond)

	for range recheckTicker.C {
		// The functonality here was copied from: https://gist.github.com/obonyojimmy/d6b263212a011ac7682ac738b7fb4c70
		mod := windows.NewLazyDLL("user32.dll")

		proc := mod.NewProc("GetForegroundWindow")
		hwnd, _, _ := proc.Call()

		proc = mod.NewProc("GetWindowTextLengthW")
		ret, _, _ := proc.Call(hwnd)

		buf := make([]uint16, int(ret)+1)
		proc = mod.NewProc("GetWindowTextW")
		proc.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), uintptr(int(ret)+1))

		title := syscall.UTF16ToString(buf)

		if title != "Quick-Search" {
			runtime.WindowHide(a.CTX)
			runtime.EventsEmit(a.CTX, "hidApp")
		}
	}
}
