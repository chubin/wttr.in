package ansitopng

import (
	"io/fs"
	"log"
	"os"

	"github.com/chubin/wttr.in/internal/assets"
	"github.com/fogleman/gg"
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

// Cache to prevent loading the same font thousands of times
var loadedFontCategories = make(map[string]bool)

// === Font loading with caching ===
func loadAndSetFont(dc *gg.Context, cat string) {
	if loadedFontCategories[cat] {
		return // already loaded
	}

	path, ok := fontPaths[cat]
	if !ok {
		path = fontPaths["default"]
	}

	// IMPORTANT: embed/ prefix as per your build structure
	path = "embed/" + path

	data, err := fs.ReadFile(assets.FS, path)
	if err != nil {
		log.Printf("ERROR: Could not read font %s: %v", path, err)
		loadedFontCategories[cat] = false
		return
	}

	tmpfile, err := os.CreateTemp("", "wttr-font-*.ttf")
	if err != nil {
		log.Printf("ERROR: Could not create temp font file: %v", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(data); err != nil {
		log.Printf("ERROR: Could not write temp font: %v", err)
		tmpfile.Close()
		return
	}
	tmpfile.Close()

	if err := dc.LoadFontFace(tmpfile.Name(), FONT_SIZE); err != nil {
		log.Printf("ERROR: Failed to load font %s: %v", path, err)
		loadedFontCategories[cat] = false
		return
	}

	loadedFontCategories[cat] = true
	log.Printf("INFO: Font loaded for category '%s'", cat)
}
