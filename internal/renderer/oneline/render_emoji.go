// internal/renderer/oneline/render_emoji.go
package oneline

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/chubin/wttr.in/internal/options"
)

// RenderConditionEmoji returns emoji/symbol + appropriate spacing.
// Now uses the new profile-based system (?emoji=...) while preserving
// backwards compatibility with ?view=v2d/v2n and ?StandardFont.
func RenderConditionEmoji(ctx *RenderContext) string {
	if ctx == nil || ctx.Data == nil || ctx.Data.ConditionCode == "" {
		fmt.Fprintln(os.Stderr, "WARNING: RenderConditionEmoji called with empty or nil ConditionCode")
		return "❓"
	}

	codeStr := ctx.Data.ConditionCode

	// Validate code is numeric
	if _, err := strconv.Atoi(codeStr); err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: invalid weather code (not numeric): %q\n", codeStr)
		return "❓"
	}

	// Get symbolic name
	name, ok := WWOCodeToName[codeStr]
	if !ok {
		fmt.Fprintf(os.Stderr, "WARNING: unknown weather code %q – not in WWOCodeToName\n", codeStr)
		name = "Unknown"
	}

	// Determine active profile (new ?emoji= option takes precedence)
	profileName := GetEmojiProfile(ctx.Options) // helper from emoji_profile.go
	profile, exists := EmojiProfiles[profileName]
	if !exists {
		log.Println(profileName, "not found")
		profile = EmojiProfiles["unicode"] // safe fallback
	}

	// Special case: respect v2d/v2n even if ?emoji= is set to something else
	view := ""
	if ctx.Options != nil {
		view = ctx.Options.View
	}
	if view == "v2n" {
		profile = EmojiProfiles["nerd"] // can be made smarter with night symbols
	} else if view == "v2d" {
		profile = EmojiProfiles["nerd"]
	}

	// Get symbol from the selected profile
	symbol, found := profile.Symbols[name]
	if !found {
		fmt.Fprintf(os.Stderr, "WARNING: no symbol defined for %q (code %s) in profile %s\n", name, codeStr, profile.Name)
		symbol = "✨"
	}

	// Apply profile-specific width for alignment
	width := profile.Width(symbol)
	padding := 0
	if width < 3 { // keep classic 3-cell target for nice spacing
		padding = 3 - width
		if padding < 0 {
			padding = 0
		}
	}

	if padding == 2 {
		return " " + symbol + strings.Repeat(" ", padding-1)
	}
	return symbol + strings.Repeat(" ", padding)
}

// GetEmojiProfile returns the active emoji profile name based on options.
func GetEmojiProfile(opts *options.Options) string {
	if opts == nil || opts.Emoji == "" {
		if opts.View == "v2d" || opts.View == "v2n" {
			return "nerd"
		}
		return "unicode"
	}
	return opts.Emoji
}
