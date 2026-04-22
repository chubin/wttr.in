package ansitopng

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io/fs"
	"log"
	"path/filepath"
	"strconv"
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

var fontPaths = map[string]string{
	"default":  "fonts/DejaVuSansMono.ttf",
	"Cyrillic": "fonts/DejaVuSansMono.ttf",
	"Greek":    "fonts/DejaVuSansMono.ttf",
	"Arabic":   "fonts/DejaVuSansMono.ttf",
	"Hebrew":   "fonts/DejaVuSansMono.ttf",
	"Han":      "fonts/wqy-zenhei.ttc",
	"Hiragana": "fonts/MTLc3m.ttf",
	"Katakana": "fonts/MTLc3m.ttf",
	"Hangul":   "fonts/LexiGulim.ttf",
	"Braille":  "fonts/Symbola_hint.ttf",
	"Emoji":    "fonts/Symbola_hint.ttf",
}

// You can adjust this map if needed (from your constants.WEATHER_SYMBOL_WIDTH_VTE)
var weatherSymbolWidth = map[string]int{
	// Example: add your wide symbols here, e.g. "☀": 2,
}

type Options map[string]string

// RenderANSI renders ANSI text to PNG and returns the image bytes.
func RenderANSI(text string, opts Options) ([]byte, error) {
	if opts == nil {
		opts = make(Options)
	}

	// Fix graphemes (replace complex ones with placeholder "!")
	text, graphemes := fixGraphemes(text)

	// Create screen and stream (mimics pyte)
	screen := te.NewScreen(COLS, ROWS)
	screen.SetMode(te.ModeLNM) // Line Feed / New Line Mode

	stream := te.NewStream(screen, false)
	stream.Feed([]byte(text))

	// Convert screen buffer to 2D slice of cells (similar to pyte buf)
	buf := screenToBuffer(screen)

	// Strip empty lines and trailing spaces
	buf = stripBuffer(buf)

	return genTerm(buf, graphemes, opts)
}

// screenToBuffer converts go-te screen to [][]te.Char (row-major)
func screenToBuffer(screen *te.Screen) [][]te.Char {
	rows := make([][]te.Char, screen.Height)
	for y := 0; y < screen.Height; y++ {
		row := make([]te.Char, screen.Width)
		for x := 0; x < screen.Width; x++ {
			cell := screen.Cell(x, y)
			row[x] = cell // te.Char has Data, FG, BG, etc.
		}
		rows[y] = row
	}
	return rows
}

// fixGraphemes replaces complex graphemes with "!" and returns them separately
func fixGraphemes(text string) (string, []string) {
	var builder strings.Builder
	var graphemes []string

	gr := uniseg.NewGraphemes(text)
	for gr.Next() {
		gra := gr.Str()
		if uniseg.GraphemeClusterCount([]byte(gra)) > 1 || len([]rune(gra)) > 1 {
			builder.WriteByte('!')
			graphemes = append(graphemes, gra)
		} else {
			builder.WriteString(gra)
		}
	}
	return builder.String(), graphemes
}

func stripBuffer(buf [][]te.Char) [][]te.Char {
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

	// Trim trailing spaces per line and find max width
	maxLen := 0
	for _, line := range buf {
		l := lineLength(line)
		if l > maxLen {
			maxLen = l
		}
	}

	for i := range buf {
		buf[i] = buf[i][:maxLen]
	}

	return buf
}

func isEmptyLine(line []te.Char) bool {
	for _, c := range line {
		if c.Data != " " && c.Data != "" {
			return false
		}
	}
	return true
}

func lineLength(line []te.Char) int {
	for i := len(line) - 1; i >= 0; i-- {
		if line[i].Data != " " && line[i].Data != "" {
			return i + 1
		}
	}
	return 0
}

func colorMapping(c string, inverse bool) color.Color {
	if c == "default" {
		if inverse {
			return color.Black
		}
		return color.RGBA{211, 211, 211, 255} // lightgray
	}

	switch c {
	case "black":
		return color.Black
	case "green":
		return color.RGBA{0, 128, 0, 255}
	case "cyan":
		return color.Cyan
	case "blue":
		return color.Blue
	case "brown":
		return color.RGBA{165, 42, 42, 255}
	}

	// Try RGB hex like "ff0000"
	if len(c) == 6 {
		if r, err := strconv.ParseUint(c[0:2], 16, 8); err == nil {
			if g, err := strconv.ParseUint(c[2:4], 16, 8); err == nil {
				if b, err := strconv.ParseUint(c[4:6], 16, 8); err == nil {
					return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
				}
			}
		}
	}

	// Fallback to black
	return color.Black
}

func loadEmojiLib() (map[string]image.Image, error) {
	emojilib := make(map[string]image.Image)

	emojiFS, err := fs.Sub(assets.FS, "share/emoji")
	if err != nil {
		return nil, err
	}

	err = fs.WalkDir(emojiFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if !strings.HasSuffix(path, ".png") {
			return nil
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

		// Resize to CHAR_HEIGHT x CHAR_HEIGHT
		resized := resizeImage(img, CHAR_HEIGHT, CHAR_HEIGHT)
		char := strings.TrimSuffix(filepath.Base(path), ".png")
		emojilib[char] = resized
		return nil
	})

	return emojilib, err
}

func resizeImage(src image.Image, w, h int) image.Image {
	// Simple nearest-neighbor resize (you can improve with draw or external lib)
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	// For production, consider github.com/nfnt/resize or gg's own scaling
	// Here we use a simple approach for brevity
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			srcX := x * src.Bounds().Dx() / w
			srcY := y * src.Bounds().Dy() / h
			dst.Set(x, y, src.At(srcX, srcY))
		}
	}
	return dst
}

func genTerm(buf [][]te.Char, graphemes []string, opts Options) ([]byte, error) {
	cols := 0
	for _, line := range buf {
		if len(line) > cols {
			cols = len(line)
		}
	}
	rows := len(buf)
	if rows == 0 {
		rows = 1
	}

	inverted := opts["inverted_colors"] == "true" || opts["inverted_colors"] == "1"

	bg := colorMapping(opts["background"], inverted)
	if opts["background"] == "" {
		bg = color.Black // default dark background
	}

	dc := gg.NewContext(cols*CHAR_WIDTH, rows*CHAR_HEIGHT)
	dc.SetColor(bg)
	dc.Clear()

	// Load fonts (you must embed or provide real font files)
	fonts := make(map[string]*gg.FontFace)
	for cat, path := range fontPaths {
		// In real code: load from embed.FS or os.Open
		// Example placeholder:
		f, err := assets.FS.Open(path)
		if err != nil {
			log.Printf("Warning: font %s not found: %v", cat, err)
			continue
		}
		// gg.LoadFontFace would be used here in practice
		_ = f.Close() // placeholder
		// fonts[cat] = ... (implement proper font loading)
	}

	emojilib, _ := loadEmojiLib()

	currentGrapheme := 0
	yPos := 0

	for _, line := range buf {
		xPos := 0
		for _, cell := range line {
			fg := colorMapping(cell.FG, inverted)

			// Background rectangle
			if cell.BG != "default" {
				bgCol := colorMapping(cell.BG, inverted)
				dc.SetColor(bgCol)
				dc.DrawRectangle(float64(xPos), float64(yPos), float64(CHAR_WIDTH), float64(CHAR_HEIGHT))
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

				if cat == "Emoji" {
					if img, ok := emojilib[data]; ok {
						dc.DrawImage(img, xPos, yPos)
					} else {
						// fallback to text
						drawText(dc, data, xPos, yPos, fg, fonts, cat)
					}
				} else {
					drawText(dc, data, xPos, yPos, fg, fonts, cat)
				}
			}

			width := CHAR_WIDTH * getSymbolWidth(data)
			xPos += width
		}
		yPos += CHAR_HEIGHT
	}

	// Transparency support
	if alphaStr := opts["transparency"]; alphaStr != "" {
		alpha, _ := strconv.Atoi(alphaStr)
		if alpha < 0 {
			alpha = 0
		}
		if alpha > 255 {
			alpha = 255
		}
		// gg context is RGB; for RGBA transparency you may need to post-process or use another lib.
		// For simplicity, we skip full alpha here or implement via image manipulation.
	}

	var bufBytes bytes.Buffer
	err := png.Encode(&bufBytes, dc.Image())
	if err != nil {
		return nil, err
	}
	return bufBytes.Bytes(), nil
}

func drawText(dc *gg.Context, text string, x, y int, col color.Color, fonts map[string]*gg.FontFace, cat string) {
	// Use appropriate font or default
	// dc.SetFontFace(fonts[cat] or default)
	dc.SetColor(col)
	dc.DrawStringAnchored(text, float64(x), float64(y), 0, 1) // adjust anchoring
}

func getSymbolWidth(s string) int {
	if w, ok := weatherSymbolWidth[s]; ok {
		return w
	}
	return 1
}

func scriptCategory(r rune) string {
	// Implement your _script_category logic here (emoji check + unicodedata script)
	// For brevity, placeholder:
	if isEmoji(r) {
		return "Emoji"
	}
	// ... map to "Han", "default", etc.
	return "default"
}

func isEmoji(r rune) bool {
	// Simple check; improve with proper emoji detection
	return r >= 0x1F000 && r <= 0x1FAFF
}
