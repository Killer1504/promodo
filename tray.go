package main

import (
	"context"
	_ "embed"
	"log/slog"
	"runtime"

	"github.com/energye/systray"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed build/windows/icon.ico
var appIcon []byte

// startTray sets up the system tray icon with Show / Quit menu items.
// Pinned to an OS thread — required for the Windows message pump.
func startTray(ctx context.Context) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	systray.Run(func() {
		systray.SetIcon(appIcon)
		systray.SetTitle("Pomodoro Focus Timer")
		systray.SetTooltip("Pomodoro Focus Timer — running in background")

		mShow := systray.AddMenuItem("🍅 Show Window", "Restore the timer window")
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Quit", "Exit Pomodoro Focus Timer")

		mShow.Click(func() {
			wailsRuntime.WindowShow(ctx)
		})

		mQuit.Click(func() {
			slog.Info("quit from system tray")
			systray.Quit()
			wailsRuntime.Quit(ctx)
		})

	}, func() {})
}
