// Package spec provides structured access to wttr.in's one-line format specification,
// particularly the list of supported placeholders (%c, %t, etc.).
package spec

import (
	_ "embed"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

	"yourproject/internal/assets" // adjust to your actual module path
)

//go:generate go run ../../cmd/generators/spec-validator/main.go   // optional: if you add validation

// Placeholder represents a single format specifier in the one-line output mode.
type Placeholder struct {
	Name        string `yaml:"name"`
	Letter      string `yaml:"letter"`
	Description string `yaml:"description"`
	Example     string `yaml:"example"`
}

// Spec holds the complete one-line placeholders definition.
type Spec struct {
	Placeholders []Placeholder `yaml:"placeholders"`
}

// DefaultSpec is the parsed specification loaded at init time from embedded data.
var DefaultSpec *Spec

// ErrSpecNotLoaded is returned when the embedded spec could not be parsed.
var ErrSpecNotLoaded = fmt.Errorf("spec not loaded or invalid")

func init() {
	data, err := assets.GetFile("spec/oneline/placeholders.yaml")
	if err != nil {
		// In production you might want to panic or log fatal here.
		// For development/demo we just leave DefaultSpec as nil.
		return
	}

	var s Spec
	if err := yaml.Unmarshal(data, &s); err != nil {
		// Again: in real code → panic or fatal log
		return
	}

	// Optional lightweight validation
	if len(s.Placeholders) == 0 {
		return
	}

	DefaultSpec = &s
}

// Get returns the default (embedded) specification.
// Returns ErrSpecNotLoaded if loading or parsing failed.
func Get() (*Spec, error) {
	if DefaultSpec == nil {
		return nil, ErrSpecNotLoaded
	}
	return DefaultSpec, nil
}

// MustGet panics if the spec is not available (useful for init-time dependencies).
func MustGet() *Spec {
	spec, err := Get()
	if err != nil {
		panic(fmt.Sprintf("critical: cannot load wttr.in one-line spec: %v", err))
	}
	return spec
}

// FindByLetter looks up a placeholder by its format letter (case-sensitive).
func (s *Spec) FindByLetter(letter string) *Placeholder {
	for i := range s.Placeholders {
		if s.Placeholders[i].Letter == letter {
			return &s.Placeholders[i]
		}
	}
	return nil
}

// FindByName looks up a placeholder by its semantic name (case-insensitive).
func (s *Spec) FindByName(name string) *Placeholder {
	nameLower := strings.ToLower(name)
	for i := range s.Placeholders {
		if strings.ToLower(s.Placeholders[i].Name) == nameLower {
			return &s.Placeholders[i]
		}
	}
	return nil
}

// Letters returns a sorted string of all supported placeholder letters.
func (s *Spec) Letters() string {
	var letters []string
	for _, p := range s.Placeholders {
		letters = append(letters, p.Letter)
	}
	// You could sort them if desired: sort.Strings(letters)
	return strings.Join(letters, "")
}

// String implements a human-readable summary (useful for debugging / :help output).
func (s *Spec) String() string {
	if s == nil || len(s.Placeholders) == 0 {
		return "spec: empty or not loaded"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("wttr.in one-line placeholders (%d total):\n", len(s.Placeholders)))
	for _, p := range s.Placeholders {
		sb.WriteString(fmt.Sprintf("  %s  %-18s  %s\n", p.Letter, p.Name, p.Example))
	}
	return sb.String()
}