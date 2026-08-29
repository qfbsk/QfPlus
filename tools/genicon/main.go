// Command genicon renders the QfPlus mark and writes every icon artifact the
// build needs. The artwork is generated rather than drawn by hand so the same
// geometry scales from a 16px taskbar tile to the 1024px Wails app icon.
//
// Usage: go run ./tools/genicon [project-root]   (defaults to the working dir)
// Optional second argument writes an 8x magnified preview PNG of the small sizes.
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
)

const (
	fieldInk = 13 // #0D0D0D, the app's --ink token

	qcX = 0.410
	qcY = 0.565
	qOuter = 0.240
	qInner = 0.135

	tailFrom = 0.115
	tailTo   = 0.345
	tailHalf = 0.0525

	pcX = 0.730
	pcY = 0.275
	armHalf = 0.105
	armThick = 0.045
)

// roundedRectCoverage reports whether (x, y) lands inside the icon field.
func inField(x, y float64) bool {
	const radius = 0.2237
	cx := clamp(x, radius, 1-radius)
	cy := clamp(y, radius, 1-radius)
	dx := x - cx
	dy := y - cy
	return dx*dx+dy*dy <= radius*radius
}

func inRing(x, y float64) bool {
	dx := x - qcX
	dy := y - qcY
	d2 := dx*dx + dy*dy
	return d2 >= qInner*qInner && d2 <= qOuter*qOuter
}

// inTail is the Q's diagonal stem, rotated 45 degrees and pointing down-right.
func inTail(x, y float64) bool {
	const invRoot2 = 0.7071067811865476
	dx := x - qcX
	dy := y - qcY
	along := (dx + dy) * invRoot2
	across := (dx - dy) * invRoot2
	return along >= tailFrom && along <= tailTo && abs(across) <= tailHalf
}

func inPlus(x, y float64) bool {
	dx := abs(x - pcX)
	dy := abs(y - pcY)
	return (dx <= armThick && dy <= armHalf) || (dy <= armThick && dx <= armHalf)
}

func inMark(x, y float64) bool {
	return inRing(x, y) || inTail(x, y) || inPlus(x, y)
}

func render(size int) *image.RGBA {
	samples := 4
	if size < 256 {
		samples = 8
	}
	step := 1.0 / float64(samples)
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	inv := 1 / float64(size)
	for py := 0; py < size; py++ {
		for px := 0; px < size; px++ {
			fieldHits, markHits := 0, 0
			for sy := 0; sy < samples; sy++ {
				fy := (float64(py) + (float64(sy)+0.5)*step) * inv
				for sx := 0; sx < samples; sx++ {
					fx := (float64(px) + (float64(sx)+0.5)*step) * inv
					if inField(fx, fy) {
						fieldHits++
						if inMark(fx, fy) {
							markHits++
						}
					}
				}
			}
			if fieldHits == 0 {
				continue
			}
			alpha := uint8(255 * fieldHits / (samples * samples))
			lum := uint8(255 * markHits / (samples * samples))
			v := uint8(fieldInk) + uint8((int(255-fieldInk)*int(lum))/255)
			i := img.PixOffset(px, py)
			scaled := uint16(v) * uint16(alpha) / 255
			img.Pix[i+0] = uint8(scaled)
			img.Pix[i+1] = uint8(scaled)
			img.Pix[i+2] = uint8(scaled)
			img.Pix[i+3] = alpha
		}
	}
	return img
}

func pngBytes(img image.Image) []byte {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// writeICO packs full-colour PNG payloads, which Vista and makensis both accept.
func writeICO(path string, sizes []int) {
	var entries, payloads bytes.Buffer
	dir := make([]byte, 6)
	binary.LittleEndian.PutUint16(dir[2:], 1)
	binary.LittleEndian.PutUint16(dir[4:], uint16(len(sizes)))

	offset := int32(6 + 16*len(sizes))
	for _, size := range sizes {
		data := pngBytes(render(size))
		header := make([]byte, 16)
		if size >= 256 {
			header[0], header[1] = 0, 0
		} else {
			header[0], header[1] = byte(size), byte(size)
		}
		binary.LittleEndian.PutUint16(header[4:], 1)
		binary.LittleEndian.PutUint16(header[6:], 32)
		binary.LittleEndian.PutUint32(header[8:], uint32(len(data)))
		binary.LittleEndian.PutUint32(header[12:], uint32(offset))
		entries.Write(header)
		payloads.Write(data)
		offset += int32(len(data))
	}
	buf := make([]byte, 0, 6+entries.Len()+payloads.Len())
	buf = append(buf, dir...)
	buf = append(buf, entries.Bytes()...)
	buf = append(buf, payloads.Bytes()...)
	if err := os.WriteFile(path, buf, 0o644); err != nil {
		panic(err)
	}
}

func writeFile(path string, data []byte) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		panic(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		panic(err)
	}
	fmt.Printf("%-44s %8.1f KiB\n", path, float64(len(data))/1024)
}

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	writeFile(filepath.Join(root, "build", "appicon.png"), pngBytes(render(1024)))
	writeFile(filepath.Join(root, "frontend", "src", "assets", "icons", "icon.png"), pngBytes(render(64)))

	icoPath := filepath.Join(root, "build", "windows", "icon.ico")
	if err := os.MkdirAll(filepath.Dir(icoPath), 0o755); err != nil {
		panic(err)
	}
	writeICO(icoPath, []int{16, 24, 32, 48, 64, 128, 256})
	info, err := os.Stat(icoPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%-44s %8.1f KiB\n", icoPath, float64(info.Size())/1024)

	if len(os.Args) > 2 {
		writePreview(os.Args[2])
	}
}

// writePreview blows the small ICO sizes up so pixel-level legibility is checkable.
func writePreview(path string) {
	sizes := []int{16, 24, 32, 48}
	const zoom = 8
	gap := 12 * zoom
	width := 0
	for _, size := range sizes {
		width += size*zoom + gap
	}
	out := image.NewRGBA(image.Rect(0, 0, width, 64*zoom))
	for y := range out.Pix {
		out.Pix[y] = 235
	}
	for i := 0; i < width*64*zoom*4; i += 4 {
		out.Pix[i+3] = 255
	}
	x := gap
	for _, size := range sizes {
		small := render(size)
		for py := 0; py < size*zoom; py++ {
			for px := 0; px < size*zoom; px++ {
				src := small.PixOffset(px/zoom, py/zoom)
				dst := out.PixOffset(x+px, (64-size)*zoom/2+py)
				copy(out.Pix[dst:dst+4], small.Pix[src:src+4])
			}
		}
		x += size*zoom + gap
	}
	writeFile(path, pngBytes(out))
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
