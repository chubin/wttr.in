// internal/formatter/ansitopng/fonts.go
package ansitopng

import (
	"io/fs"
	"log"
	"os"
	"sync"

	"github.com/fogleman/gg"

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

// fontFiles holds the temporary file paths (extracted once at startup)
var (
	fontFiles = make(map[string]string) // category → real temp file path on disk
	fontMu    sync.RWMutex
)

// preloadFontFiles extracts fonts from embed.FS to temp files once at startup
func preloadFontFiles() {
	fontMu.Lock()
	defer fontMu.Unlock()

	for cat, relPath := range fontPaths {
		data, err := fs.ReadFile(assets.FS, "embed/"+relPath)
		if err != nil {
			log.Printf("ERROR: Could not read font %s from embed: %v", relPath, err)
			continue
		}

		tmp, err := os.CreateTemp("", "wttr-font-"+cat+"-*.ttf")
		if err != nil {
			log.Printf("ERROR: Could not create temp file for %s: %v", cat, err)
			continue
		}

		if _, err := tmp.Write(data); err != nil {
			tmp.Close()
			log.Printf("ERROR: Could not write font %s: %v", cat, err)
			continue
		}
		tmp.Close()

		fontFiles[cat] = tmp.Name()
		log.Printf("INFO: Extracted font for category '%s'", cat)
	}
}

func init() {
	preloadFontFiles()
}

func loadAndSetFont(dc *gg.Context, cat string) {
	fontMu.RLock()
	path, ok := fontFiles[cat]
	if !ok {
		path = fontFiles["default"]
	}
	fontMu.RUnlock()

	if path == "" {
		log.Printf("WARNING: No font path for '%s'", cat)
		dc.LoadFontFace("", FONT_SIZE)
		return
	}

	// CRITICAL FIX: Always load a fresh face -> no sharing
	if err := dc.LoadFontFace(path, FONT_SIZE); err != nil {
		log.Printf("ERROR: Failed to load font %s: %v", cat, err)
		dc.LoadFontFace("", FONT_SIZE)
	}
}

// Optional: cleanup on shutdown (good practice)
func CleanupFonts() {
	fontMu.Lock()
	defer fontMu.Unlock()
	for _, path := range fontFiles {
		os.Remove(path)
	}
}

func scriptCategory(r rune) string {
	if isEmoji(r) {
		return "Emoji"
	}

	// Box drawing, blocks, geometric shapes, powerline — very important for wttr.in
	if (r >= 0x2500 && r <= 0x257F) || // Box Drawing
		(r >= 0x2580 && r <= 0x259F) || // Block Elements
		(r >= 0x25A0 && r <= 0x25FF) || // Geometric Shapes
		(r >= 0xE0B0 && r <= 0xE0D7) { // Powerline
		return "default"
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
	case (r >= 0x0900 && r <= 0x097F) || (r >= 0xA8E0 && r <= 0xA8FF):
		return "Devanagari"
	case r >= 0x0980 && r <= 0x09FF:
		return "Bengali"
	case r >= 0x0A00 && r <= 0x0A7F:
		return "Gurmukhi"
	}
	return "default"
}

func isEmoji(r rune) bool {
	return (r >= 0x1F000 && r <= 0x1FAFF) ||
		(r >= 0x1F600 && r <= 0x1F64F) ||
		(r >= 0x2600 && r <= 0x26FF) ||
		(r >= 0x2700 && r <= 0x27BF)
}
