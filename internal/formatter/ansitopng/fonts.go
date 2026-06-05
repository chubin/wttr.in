package ansitopng

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/fogleman/gg"

	"github.com/chubin/wttr.in/internal/assets"
)

// ScriptConfig represents one script entry from scripts.json
type ScriptConfig struct {
	Filename string  `json:"filename"`
	Package  string  `json:"package"`
	Ranges   [][]int `json:"ranges"`
}

// FontConfig holds the full configuration
type FontConfig struct {
	Scripts map[string]ScriptConfig `json:"scripts"`
}

var (
	fontConfig FontConfig
	fontFiles  = make(map[string]string) // category → temp file on disk
	fontMu     sync.RWMutex
)

// loadFontConfig loads configuration from embedded scripts.json
func loadFontConfig() {
	data, err := assets.GetFile("share/defs/fonts/scripts.json")
	if err != nil {
		log.Printf("ERROR: Could not read embedded scripts.json: %v", err)
		return
	}

	if err := json.Unmarshal(data, &fontConfig); err != nil {
		log.Printf("ERROR: Failed to parse scripts.json: %v", err)
	}
}

// getFontBasename extracts filename from full system path
func getFontBasename(fullPath string) string {
	for i := len(fullPath) - 1; i >= 0; i-- {
		if fullPath[i] == '/' || fullPath[i] == '\\' {
			return fullPath[i+1:]
		}
	}
	return fullPath
}

// preloadFontFiles extracts fonts from embed.FS to temporary files
func preloadFontFiles() {
	fontMu.Lock()
	defer fontMu.Unlock()

	for cat, cfg := range fontConfig.Scripts {
		if cfg.Filename == "" {
			continue
		}

		basename := getFontBasename(cfg.Filename)
		relPath := "fonts/" + basename

		data, err := assets.GetFile(relPath)
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

	// Ensure default exists
	if _, ok := fontFiles["default"]; !ok {
		log.Printf("WARNING: No default font loaded")
	}
}

func init() {
	loadFontConfig()
	preloadFontFiles()
}

// loadAndSetFont loads the appropriate font for a category
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

	// Always load fresh face
	if err := dc.LoadFontFace(path, FONT_SIZE); err != nil {
		log.Printf("ERROR: Failed to load font %s: %v", cat, err)
		dc.LoadFontFace("", FONT_SIZE)
	}
}

// CleanupFonts removes temporary font files (call on shutdown if needed)
func CleanupFonts() {
	fontMu.Lock()
	defer fontMu.Unlock()
	for _, path := range fontFiles {
		os.Remove(path)
	}
}

// scriptCategory returns the font category for a given rune
func scriptCategory(r rune) string {
	if isEmoji(r) {
		return "Emoji"
	}

	// Special categories that use default font
	if (r >= 0x2500 && r <= 0x257F) || // Box Drawing
		(r >= 0x2580 && r <= 0x259F) || // Block Elements
		(r >= 0x25A0 && r <= 0x25FF) || // Geometric Shapes
		(r >= 0xE0B0 && r <= 0xE0D7) { // Powerline
		return "default"
	}

	// Dynamic range check from JSON
	for cat, cfg := range fontConfig.Scripts {
		for _, rng := range cfg.Ranges {
			if len(rng) == 2 {
				if r >= rune(rng[0]) && r <= rune(rng[1]) {
					return cat
				}
			}
		}
	}

	return "default"
}

func isEmoji(r rune) bool {
	return (r >= 0x1F000 && r <= 0x1FAFF) ||
		(r >= 0x1F600 && r <= 0x1F64F) ||
		(r >= 0x2600 && r <= 0x26FF) ||
		(r >= 0x2700 && r <= 0x27BF)
}
