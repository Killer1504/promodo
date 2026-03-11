//go:build ignore

// gen_menu_icons creates transparent 16x16 PNG icons for tray menu items.
// Usage: go run build/tools/gen_menu_icons.go
package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

func main() {
	// Show Window icon: a small window outline (white/light)
	showIcon := image.NewNRGBA(image.Rect(0, 0, 64, 64))
	drawWindowIcon(showIcon, color.NRGBA{R: 76, G: 175, B: 80, A: 255}) // green

	// Quit icon: an X mark (red)
	quitIcon := image.NewNRGBA(image.Rect(0, 0, 64, 64))
	drawXIcon(quitIcon, color.NRGBA{R: 231, G: 76, B: 60, A: 255}) // red

	savePNG("build/windows/show-icon.png", showIcon)
	savePNG("build/windows/quit-icon.png", quitIcon)
	fmt.Println("✓ generated show-icon.png and quit-icon.png")
}

// drawWindowIcon draws a window/restore icon shape
func drawWindowIcon(img *image.NRGBA, c color.NRGBA) {
	w, h := 64, 64
	// Outer frame
	drawRect(img, 8, 12, w-8, h-8, c, 4)
	// Title bar (filled top strip)
	fillRect(img, 8, 12, w-8, 22, c)
}

// drawXIcon draws a bold X/close icon
func drawXIcon(img *image.NRGBA, c color.NRGBA) {
	w, h := 64, 64
	// Draw two diagonal thick lines forming an X
	thickness := 5.0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			fx, fy := float64(x), float64(y)
			// Distance to main diagonal (top-left to bottom-right)
			d1 := math.Abs(fx-fy) / math.Sqrt(2)
			// Distance to anti-diagonal (top-right to bottom-left)
			d2 := math.Abs(fx+fy-float64(w)) / math.Sqrt(2)

			// Margin from edges
			margin := 10.0
			if fx < margin || fx > float64(w)-margin || fy < margin || fy > float64(h)-margin {
				continue
			}

			if d1 < thickness || d2 < thickness {
				// Anti-alias: smooth edge
				dist := math.Min(d1, d2)
				alpha := 1.0
				if dist > thickness-1.5 {
					alpha = math.Max(0, (thickness-dist)/1.5)
				}
				img.SetNRGBA(x, y, color.NRGBA{
					R: c.R, G: c.G, B: c.B,
					A: uint8(float64(c.A) * alpha),
				})
			}
		}
	}
}

func drawRect(img *image.NRGBA, x1, y1, x2, y2 int, c color.NRGBA, thickness int) {
	for t := 0; t < thickness; t++ {
		// Top
		for x := x1; x < x2; x++ {
			img.SetNRGBA(x, y1+t, c)
		}
		// Bottom
		for x := x1; x < x2; x++ {
			img.SetNRGBA(x, y2-1-t, c)
		}
		// Left
		for y := y1; y < y2; y++ {
			img.SetNRGBA(x1+t, y, c)
		}
		// Right
		for y := y1; y < y2; y++ {
			img.SetNRGBA(x2-1-t, y, c)
		}
	}
}

func fillRect(img *image.NRGBA, x1, y1, x2, y2 int, c color.NRGBA) {
	for y := y1; y < y2; y++ {
		for x := x1; x < x2; x++ {
			img.SetNRGBA(x, y, c)
		}
	}
}

func savePNG(path string, img *image.NRGBA) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		fmt.Fprintf(os.Stderr, "encode %s: %v\n", path, err)
		os.Exit(1)
	}
	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write %s: %v\n", path, err)
		os.Exit(1)
	}
}
