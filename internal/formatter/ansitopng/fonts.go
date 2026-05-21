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
	"Gurmukhi":   "fonts/NotoSansGurmukhi-Regular.ttf",
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
	case r >= 0x0900 && r <= 0x097F || r >= 0xA8E0 && r <= 0xA8FF: // Devanagari + Extended
		return "Devanagari"
	case r >= 0x0980 && r <= 0x09FF: // Bengali
		return "Bengali"
	case r >= 0x0A00 && r <= 0x0A7F: // Gurmukhi
		return "Gurmukhi"
	}
	return "default"
}

func isEmoji(r rune) bool {
	return (r >= 0x1F000 && r <= 0x1FAFF) ||
		(r >= 0x1F600 && r <= 0x1F64F) ||
		(r >= 0x2600 && r <= 0x26FF)
}
