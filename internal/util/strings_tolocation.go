package util

import (
	"strings"
	"unicode"
)

// ToLocationCase converts a lowercase string to proper location/title case.
// Only changes casing. Preserves all original punctuation, hyphens, and structure.
func ToLocationCase(s string) string {
	if s == "" {
		return ""
	}

	// If string already has any uppercase letter, preserve original casing
	if hasUppercase(s) {
		return s
	}

	// Common words that usually remain lowercase in place names
	lowercaseWords := map[string]bool{
		"of": true, "the": true, "and": true, "or": true,
		"in": true, "on": true, "at": true, "to": true,
		"for": true, "with": true, "by": true,
		"de": true, "la": true, "las": true, "el": true, "los": true,
		"da": true, "das": true, "do": true, "dos": true,
		"van": true, "von": true, "der": true, "den": true,
		"an": true, "am": true,
		"upon": true, "sur": true, "sous": true,
	}

	// Split on spaces (preserving hyphenated words as single tokens)
	words := strings.Fields(s)
	result := make([]string, len(words))

	for i, word := range words {
		if strings.Contains(word, "-") {
			// Handle hyphenated parts (e.g. "saint-denis" → "Saint-Denis")
			result[i] = handleHyphenated(word, lowercaseWords)
		} else {
			lower := strings.ToLower(word)

			if i == 0 {
				// Always capitalize first word
				result[i] = capitalize(lower)
			} else if lowercaseWords[lower] {
				// Keep common connectors lowercase
				result[i] = lower
			} else {
				result[i] = capitalize(lower)
			}
		}
	}

	return strings.Join(result, " ")
}

// hasUppercase returns true if string contains any uppercase character
func hasUppercase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

// handleHyphenated processes hyphenated words like "new-york", "saint-denis"
func handleHyphenated(word string, lowercaseWords map[string]bool) string {
	parts := strings.Split(word, "-")
	for j, part := range parts {
		lower := strings.ToLower(part)
		if j == 0 || !lowercaseWords[lower] {
			parts[j] = capitalize(lower)
		} else {
			parts[j] = lower
		}
	}
	return strings.Join(parts, "-")
}

// capitalize returns the word with first letter uppercased
func capitalize(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
