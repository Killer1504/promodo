//go:build ignore

// make_ico converts a PNG to a Windows ICO with 16, 32, 48, 256 px resolutions.
// Uses only stdlib — no extra dependencies.
// Usage: go run build/tools/make_ico.go <input.png> <output.ico>
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: make_ico <input.png> <output.ico>")
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	must(err, "open")
	defer f.Close()

	src, _, err := image.Decode(f)
	must(err, "decode")

	sizes := []int{16, 32, 48, 256}

	var pngs [][]byte
	for _, sz := range sizes {
		resized := resizeNearestNeighbor(src, sz, sz)
		var buf bytes.Buffer
		must(png.Encode(&buf, resized), "encode")
		pngs = append(pngs, buf.Bytes())
	}

	ico := buildICO(sizes, pngs)
	must(os.WriteFile(os.Args[2], ico, 0644), "write")
	fmt.Printf("✓ %s (%d bytes, sizes: %v)\n", os.Args[2], len(ico), sizes)
}

// resizeNearestNeighbor resizes src to w×h using bilinear interpolation via stdlib draw.
func resizeNearestNeighbor(src image.Image, w, h int) *image.NRGBA {
	dst := image.NewNRGBA(image.Rect(0, 0, w, h))
	srcB := src.Bounds()
	srcW := srcB.Max.X - srcB.Min.X
	srcH := srcB.Max.Y - srcB.Min.Y
	for y := 0; y < h; y++ {
		sy := int(math.Round(float64(y) * float64(srcH) / float64(h)))
		if sy >= srcH {
			sy = srcH - 1
		}
		for x := 0; x < w; x++ {
			sx := int(math.Round(float64(x) * float64(srcW) / float64(w)))
			if sx >= srcW {
				sx = srcW - 1
			}
			dst.Set(x, y, src.At(srcB.Min.X+sx, srcB.Min.Y+sy))
		}
	}
	return dst
}

// buildICO assembles a valid Windows ICO file with embedded PNG images.
func buildICO(sizes []int, pngs [][]byte) []byte {
	n := len(sizes)
	headerSize := 6
	entrySize := 16
	dataStart := uint32(headerSize + entrySize*n)

	var buf bytes.Buffer
	// ICONDIR header
	writeU16(&buf, 0) // reserved
	writeU16(&buf, 1) // type = icon
	writeU16(&buf, uint16(n))

	// Compute data offsets
	offset := dataStart
	for i, sz := range sizes {
		w := byte(sz)
		if sz == 256 {
			w = 0
		}
		h := w
		buf.WriteByte(w)
		buf.WriteByte(h)
		buf.WriteByte(0)   // color count
		buf.WriteByte(0)   // reserved
		writeU16(&buf, 1)  // color planes
		writeU16(&buf, 32) // bits per pixel
		writeU32(&buf, uint32(len(pngs[i])))
		writeU32(&buf, offset)
		offset += uint32(len(pngs[i]))
	}

	for _, p := range pngs {
		buf.Write(p)
	}
	return buf.Bytes()
}

func writeU16(b *bytes.Buffer, v uint16) { _ = binary.Write(b, binary.LittleEndian, v) }
func writeU32(b *bytes.Buffer, v uint32) { _ = binary.Write(b, binary.LittleEndian, v) }
func must(err error, ctx string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error (%s): %v\n", ctx, err)
		os.Exit(1)
	}
}

// Ensure image/png decoder is registered
var _ = draw.Src
