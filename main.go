package main

import (
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:     "Pomodoro Focus Timer",
		Width:     400,
		Height:    600,
		MinWidth:  400,
		MinHeight: 600,
		MaxWidth:  400,
		MaxHeight: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 250, G: 251, B: 252, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,

		// FR-011: Hide window (don't quit) when user clicks X.
		// The app keeps running; user restores via taskbar or the in-app Quit button.
		HideWindowOnClose: true,

		// OnBeforeClose: save paused session before the process actually exits.
		OnBeforeClose: func(ctx context.Context) bool {
			if app.timer != nil && app.timer.IsPaused() {
				app.savePausedSession()
			}
			return false // allow exit
		},

		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
