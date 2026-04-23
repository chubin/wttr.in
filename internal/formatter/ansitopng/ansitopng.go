package ansitopng

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
	"github.com/rcarmo/go-te/pkg/te"
	"github.com/rivo/uniseg"

	"github.com/chubin/wttr.in/internal/assets"
)

const (
	COLS        = 180
	ROWS        = 100
	CHAR_WIDTH  = 8
	CHAR_HEIGHT = 14
	FONT_SIZE   = 13
)

var weatherSymbolWidth = map[string]int{}

// RenderANSI renders ANSI text to PNG using strongly-typed PNGOptions.
func RenderANSI(text string, opts PNGOptions) ([]byte, error) {
	// Use defaults if zero value is passed
	if opts.Transparency == 0 {
		opts.Transparency = 255
	}

	// Fix graphemes (replace complex ones with placeholder "!")
	text, graphemes := fixGraphemes(text)

	// Create screen and stream
	screen := te.NewScreen(COLS, ROWS)
	screen.SetMode([]int{te.ModeLNM}, false)

	stream := te.NewStream(screen, false)
	if err := stream.Feed(text); err != nil {
		return nil, err
	}

	// Convert screen buffer
	buf := screenToBuffer(screen)

	// Strip trailing empty lines and spaces
	buf = stripBuffer(buf)

	return genTerm(buf, graphemes, opts)
}

// screenToBuffer returns the screen content as [][]te.Cell
func screenToBuffer(screen *te.Screen) [][]te.Cell {
	return screen.LinesCells()
}

// fixGraphemes replaces complex graphemes with "!" and returns them separately
func fixGraphemes(text string) (string, []string) {
	var builder strings.Builder
	var graphemes []string
	gr := uniseg.NewGraphemes(text)
	for gr.Next() {
		gra := gr.Str()
		if uniseg.GraphemeClusterCount(gra) > 1 || len([]rune(gra)) > 1 {
			builder.WriteByte('!')
			graphemes = append(graphemes, gra)
		} else {
			builder.WriteString(gra)
		}
	}
	return builder.String(), graphemes
}

func stripBuffer(buf [][]te.Cell) [][]te.Cell {
	// Remove trailing empty lines
	for len(buf) > 0 {
		if !isEmptyLine(buf[len(buf)-1]) {
			break
		}
		buf = buf[:len(buf)-1]
	}
	if len(buf) == 0 {
		return buf
	}

	// Trim trailing spaces and find max width
	maxLen := 0
	for _, line := range buf {
		l := lineLength(line)
		if l > maxLen {
			maxLen = l
		}
	}

	for i := range buf {
		if maxLen <= len(buf[i]) {
			buf[i] = buf[i][:maxLen]
		}
	}
	return buf
}

func isEmptyLine(line []te.Cell) bool {
	for _, c := range line {
		if c.Data != " " && c.Data != "" {
			return false
		}
	}
	return true
}

func lineLength(line []te.Cell) int {
	for i := len(line) - 1; i >= 0; i-- {
		if line[i].Data != " " && line[i].Data != "" {
			return i + 1
		}
	}
	return 0
}

func loadEmojiLib() (map[string]image.Image, error) {
	emojilib := make(map[string]image.Image)
	emojiFS, err := fs.Sub(assets.FS, "share/emoji")
	if err != nil {
		log.Printf("Warning: could not open emoji assets: %v", err)
		return emojilib, err
	}

	err = fs.WalkDir(emojiFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".png") {
			return err
		}
		f, err := emojiFS.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		img, _, err := image.Decode(f)
		if err != nil {
			return err
		}

		resized := resizeImage(img, CHAR_HEIGHT, CHAR_HEIGHT)
		char := strings.TrimSuffix(filepath.Base(path), ".png")
		emojilib[char] = resized
		return nil
	})

	return emojilib, err
}

func resizeImage(src image.Image, w, h int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			srcX := x * src.Bounds().Dx() / w
			srcY := y * src.Bounds().Dy() / h
			dst.Set(x, y, src.At(srcX, srcY))
		}
	}
	return dst
}

func genTerm(buf [][]te.Cell, graphemes []string, opts PNGOptions) ([]byte, error) {
	if len(buf) == 0 {
		buf = [][]te.Cell{{}}
	}

	cols := 0
	for _, line := range buf {
		if len(line) > cols {
			cols = len(line)
		}
	}
	rows := len(buf)

	dc := gg.NewContext(cols*CHAR_WIDTH, rows*CHAR_HEIGHT)

	// Set background
	bg := opts.Background
	if opts.Inverted {
		r, g, b, a := bg.RGBA()
		bg = color.RGBA{uint8(255 - r>>8), uint8(255 - g>>8), uint8(255 - b>>8), uint8(a >> 8)}
	}
	dc.SetColor(bg)
	dc.Clear()

	emojilib, _ := loadEmojiLib()

	currentGrapheme := 0
	yPos := 0.0

	for _, line := range buf {
		xPos := 0.0
		for _, cell := range line {
			fg := colorFromANSI(cell.Attr.Fg, opts.Inverted)

			if cell.Attr.Bg.Mode != te.ColorDefault && cell.Attr.Bg.Name != "default" {
				bgCol := colorFromANSI(cell.Attr.Bg, opts.Inverted)
				dc.SetColor(bgCol)
				dc.DrawRectangle(xPos, yPos, float64(CHAR_WIDTH), float64(CHAR_HEIGHT))
				dc.Fill()
			}

			data := cell.Data
			if data == "!" {
				if currentGrapheme < len(graphemes) {
					data = graphemes[currentGrapheme]
					currentGrapheme++
				}
			}

			if data != "" && data != " " {
				cat := scriptCategory([]rune(data)[0])

				// Load and set font right before drawing
				loadAndSetFont(dc, cat)

				if cat == "Emoji" {
					if img, ok := emojilib[data]; ok {
						dc.DrawImage(img, int(xPos), int(yPos))
					} else {
						drawText(dc, data, xPos, yPos, fg)
					}
				} else {
					drawText(dc, data, xPos, yPos, fg)
				}
			}

			width := float64(CHAR_WIDTH * getSymbolWidth(data))
			xPos += width
		}
		yPos += float64(CHAR_HEIGHT)
	}

	// Apply transparency if requested (post-process)
	if opts.Transparency < 255 {
		img := dc.Image()
		rgba := image.NewRGBA(img.Bounds())
		for y := 0; y < img.Bounds().Dy(); y++ {
			for x := 0; x < img.Bounds().Dx(); x++ {
				r, g, b, _ := img.At(x, y).RGBA()
				rgba.Set(x, y, color.RGBA{
					R: uint8(r >> 8),
					G: uint8(g >> 8),
					B: uint8(b >> 8),
					A: uint8(opts.Transparency),
				})
			}
		}
		dc = gg.NewContextForRGBA(rgba)
	}

	var bufBytes bytes.Buffer
	if err := png.Encode(&bufBytes, dc.Image()); err != nil {
		return nil, err
	}
	return bufBytes.Bytes(), nil
}

func drawText(dc *gg.Context, text string, x, y float64, col color.Color) {
	dc.SetColor(col)
	// TODO: Load proper fonts per script category using assets.FS + gg.LoadFontFace
	dc.DrawStringAnchored(text, x, y+float64(CHAR_HEIGHT)-2, 0, 1)
}

func getSymbolWidth(s string) int {
	if w, ok := weatherSymbolWidth[s]; ok {
		return w
	}
	return 1
}

func scriptCategory(r rune) string {
	if isEmoji(r) {
		return "Emoji"
	}
	switch {
	case r >= 0x0400 && r <= 0x04FF:
		return "Cyrillic"
	case r >= 0x0370 && r <= 0x03FF:
		return "Greek"
	case r >= 0x0600 && r <= 0x06FF:
		return "Arabic"
	case r >= 0x0590 && r <= 0x05FF:
		return "Hebrew"
	case (r >= 0x4E00 && r <= 0x9FFF) || (r >= 0x3400 && r <= 0x4DBF):
		return "Han"
	case r >= 0x3040 && r <= 0x309F:
		return "Hiragana"
	case r >= 0x30A0 && r <= 0x30FF:
		return "Katakana"
	case r >= 0xAC00 && r <= 0xD7AF:
		return "Hangul"
	case r >= 0x2800 && r <= 0x28FF:
		return "Braille"
	}
	return "default"
}

func isEmoji(r rune) bool {
	return (r >= 0x1F000 && r <= 0x1FAFF) ||
		(r >= 0x1F600 && r <= 0x1F64F) ||
		(r >= 0x2600 && r <= 0x26FF)
}
