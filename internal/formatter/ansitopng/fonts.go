// internal/formatter/ansitopng/fonts.go
package ansitopng

import (
	"io/fs"
	"log"
	"os"
	"sync"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"

	"github.com/chubin/wttr.in/internal/assets"
)

var fontPaths = map[string]string{
	"default":    "fonts/DejaVuSansMono.ttf",
	"Cyrillic":   "fonts/DejaVuSansMono.ttf",
	"Greek":      "fonts/DejaVuSansMono.ttf",
	"Arabic":     "fonts/DejaVuSansMono.ttf",
	"Hebrew":     "fonts/DejaVuSansMono.ttf",
	"Han":        "fonts/wqy-zenhei.ttc",
	"Hiragana":   "fonts/MTLc3m.ttf",
	"Katakana":   "fonts/MTLc3m.ttf",
	"Hangul":     "fonts/LexiGulim.ttf",
	"Braille":    "fonts/Symbola_hint.ttf",
	"Emoji":      "fonts/Symbola_hint.ttf",
	"Devanagari": "fonts/NotoSansDevanagari-Regular.ttf",
	"Bengali":    "fonts/NotoSansBengali-Regular.ttf",
}

// Preloaded font data – initialized once at startup
var (
	preloadOnce sync.Once
	fontFiles   = make(map[string]string)    // category → temp file path
	fontFaces   = make(map[string]font.Face) // category → pre-parsed font face
	fontMu      sync.RWMutex
)

// preloadFonts loads all fonts once at package initialization.
func preloadFonts() {
	for cat, relPath := range fontPaths {
		fullPath := "embed/" + relPath

		data, err := fs.ReadFile(assets.FS, fullPath)
		if err != nil {
			log.Printf("ERROR: Could not read font %s: %v", fullPath, err)
			continue
		}

		tmp, err := os.CreateTemp("", "wttr-font-"+cat+"-*.ttf")
		if err != nil {
			log.Printf("ERROR: Could not create temp file for font %s: %v", cat, err)
			continue
		}

		if _, err := tmp.Write(data); err != nil {
			tmp.Close()
			log.Printf("ERROR: Could not write font %s: %v", cat, err)
			continue
		}
		tmp.Close()

		// Keep the temp file for the lifetime of the process
		fontFiles[cat] = tmp.Name()

		// Pre-parse the font face (best performance)
		face, err := gg.LoadFontFace(tmp.Name(), FONT_SIZE)
		if err != nil {
			log.Printf("ERROR: Failed to load font face for %s: %v", cat, err)
			continue
		}

		fontFaces[cat] = face
		log.Printf("INFO: Preloaded font for category '%s'", cat)
	}
}

func init() {
	preloadOnce.Do(preloadFonts)
}

// loadAndSetFont switches the font on the given context.
// Works correctly on every RenderANSI call.
func loadAndSetFont(dc *gg.Context, cat string) {
	fontMu.RLock()
	defer fontMu.RUnlock()

	if face, ok := fontFaces[cat]; ok {
		dc.SetFontFace(face)
		return
	}

	// Fallback: load by path
	if path, ok := fontFiles[cat]; ok {
		if err := dc.LoadFontFace(path, FONT_SIZE); err != nil {
			log.Printf("ERROR: Failed to load font %s on context: %v", cat, err)
		}
		return
	}

	// Ultimate fallback to default
	if path, ok := fontFiles["default"]; ok {
		_ = dc.LoadFontFace(path, FONT_SIZE)
	}
}
