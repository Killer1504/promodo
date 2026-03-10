package main

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"log/slog"

	"github.com/energye/systray"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// startTray sets up the system tray icon with Show / Quit menu items.
// Must be called in a goroutine; blocks until systray exits.
func startTray(ctx context.Context) {
	systray.Run(func() {
		// Tomato-red circle icon — generated in memory, no file needed
		systray.SetIcon(trayIconPNG())
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

	}, func() {
		// onExit
	})
}

// trayIconPNG generates a 16×16 tomato-red circle as a PNG byte slice.
func trayIconPNG() []byte {
	const size = 16
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	cx, cy := float64(size)/2, float64(size)/2
	r := cx - 1

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := float64(x) + 0.5 - cx
			dy := float64(y) + 0.5 - cy
			if dx*dx+dy*dy <= r*r {
				img.Set(x, y, color.RGBA{R: 219, G: 53, B: 53, A: 255})
			}
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil
	}
	return buf.Bytes()
}
