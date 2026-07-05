// Package localization provides convenient, language-bound helpers
// around the Localizer interface.
//
// It helps avoid repeatedly passing the language code and reduces
// boilerplate in renderers and other components.
package localization

import (
	"github.com/chubin/wttr.in/internal/options"
)

// Localizer provides localized strings for a specific language.
type Localizer interface {
	// Text returns a translated string by key.
	// It should be used for messages, captions, v1/v2 keys, etc.
	// Do NOT pass condition codes here.
	Text(lang, key string) string

	// Condition returns translated weather condition by numeric code.
	Condition(lang string, code int) string

	// ConditionByName returns translated condition by its English name.
	ConditionByName(lang, englishName string) string

	// File returns raw file content (help.txt, etc.)
	File(lang, name string) (string, error)
}

// LocalizeFunc is a language-bound text localization function.
type LocalizeFunc func(key string) string

// ConditionFunc returns a translated weather condition by numeric code.
type ConditionFunc func(code int) string

// ConditionByNameFunc returns a translated condition by its English name.
type ConditionByNameFunc func(englishName string) string

// L10n is a convenient wrapper that provides pre-bound localization
// functions for a specific language from Options.
//
// Use `l10n` as the variable name to avoid conflict with `loc` (Location).
type L10n struct {
	// Text returns translated string for messages, captions, v1/v2 keys, etc.
	Text LocalizeFunc

	// Condition returns translated weather condition by numeric code.
	Condition ConditionFunc

	// ConditionByName returns translated condition by English name.
	ConditionByName ConditionByNameFunc

	// Lang is the active language code (e.g. "en", "de", "ru").
	Lang string

	// Underlying localizer (kept for advanced methods like File())
	localizer Localizer
}

// New creates a new L10n wrapper bound to the language specified in Options.
// Falls back to English if no language is set.
func New(l Localizer, opts *options.Options) L10n {
	lang := "en"
	if opts != nil && opts.Lang != "" {
		lang = opts.Lang
	}

	return L10n{
		Text: func(key string) string {
			if key == "" {
				return ""
			}
			return l.Text(lang, key)
		},

		Condition: func(code int) string {
			return l.Condition(lang, code)
		},

		ConditionByName: func(englishName string) string {
			if englishName == "" {
				return ""
			}
			return l.ConditionByName(lang, englishName)
		},

		Lang:      lang,
		localizer: l,
	}
}

// File returns the content of a localized file (e.g. help.txt, conditions.txt).
func (l10n L10n) File(name string) (string, error) {
	return l10n.localizer.File(l10n.Lang, name)
}

// WithLang returns a new L10n instance with a different language.
// Useful for forcing a specific language in some contexts.
func (l10n L10n) WithLang(lang string) L10n {
	if lang == "" {
		lang = "en"
	}
	return New(l10n.localizer, &options.Options{Lang: lang})
}

// IsRTL reports whether the active language is written right-to-left.
func (l10n L10n) IsRTL() bool {
	switch l10n.Lang {
	case "he", "ar", "fa":
		return true
	default:
		return false
	}
}
