package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"image"
	"image/color"
	"image/png"
	"log/slog"
	"runtime"

	"github.com/energye/systray"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// startTray sets up the system tray icon with Show / Quit menu items.
// Called in a goroutine with a locked OS thread (required for Windows message pump).
func startTray(ctx context.Context) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	systray.Run(func() {
		// Windows requires an ICO format icon for LoadImage(IMAGE_ICON)
		systray.SetIcon(trayIconICO())
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
		// onExit — nothing to clean up
	})
}

// trayIconICO returns a minimal valid Windows ICO file with a 16×16 tomato-red circle.
// ICO format: ICONDIR header + ICONDIRENTRY + image data (PNG is valid in modern ICO).
func trayIconICO() []byte {
	const size = 16

	// Build a 16×16 red circle as PNG first
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

	var pngBuf bytes.Buffer
	if err := png.Encode(&pngBuf, img); err != nil {
		return nil
	}
	pngBytes := pngBuf.Bytes()

	// Build ICO container (modern ICO can embed a PNG image directly)
	// ICONDIR: reserved(2) + type(2) + count(2)
	// ICONDIRENTRY: width(1) + height(1) + colorCount(1) + reserved(1) +
	//               planes(2) + bitCount(2) + bytesInRes(4) + imageOffset(4)
	const icoHeaderSize = 6
	const icoDirEntrySize = 16
	imageOffset := uint32(icoHeaderSize + icoDirEntrySize)
	imageSize := uint32(len(pngBytes))

	var buf bytes.Buffer
	// ICONDIR
	binary.Write(&buf, binary.LittleEndian, uint16(0)) // reserved
	binary.Write(&buf, binary.LittleEndian, uint16(1)) // type = 1 (icon)
	binary.Write(&buf, binary.LittleEndian, uint16(1)) // count = 1 image

	// ICONDIRENTRY
	buf.WriteByte(uint8(size))                           // width (0 = 256)
	buf.WriteByte(uint8(size))                           // height
	buf.WriteByte(0)                                     // color count (0 = no palette)
	buf.WriteByte(0)                                     // reserved
	binary.Write(&buf, binary.LittleEndian, uint16(1))   // planes
	binary.Write(&buf, binary.LittleEndian, uint16(32))  // bit count
	binary.Write(&buf, binary.LittleEndian, imageSize)   // bytes in resource
	binary.Write(&buf, binary.LittleEndian, imageOffset) // offset to image data

	// Image data (embedded PNG — valid for Vista+ ICO)
	buf.Write(pngBytes)

	return buf.Bytes()
}
