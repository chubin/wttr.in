// Package spec provides structured access to wttr.in's one-line format specification,
// particularly the list of supported placeholders (%c, %t, etc.).
package spec

import (
	_ "embed"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"dario.cat/mergo"
	"gopkg.in/yaml.v3"

	"github.com/chubin/wttr.in/internal/assets"
)

// WttrInOptions represents the configuration for wttr.in query options and format specifiers.
type WttrInOptions struct {
	QueryOptions         []QueryOption         `yaml:"query_options"`
	FormatSpecifiers     []FormatSpecifier     `yaml:"format_specifiers"`
	PreconfiguredFormats []PreconfiguredFormat `yaml:"preconfigured_formats"`
}

// QueryOption defines a single query option for wttr.in, such as `lang` or `format`.
type QueryOption struct {
	// Name is the long name of the option (e.g., "lang", "current_only").
	Name string `yaml:"name"`

	// Short is the short name of the option (e.g., "m"), omitted if not present.
	Short string `yaml:"short,omitempty"`

	// Description of the option's purpose (e.g., "Specify output language").
	Description string `yaml:"description"`

	// Type is the data type of the option (e.g., "boolean", "string", "integer").
	Type string `yaml:"type"`

	// ValuesMap is the map of possible values to their descriptions (e.g., for `lang`, `format`).
	ValuesMap map[string]string `yaml:"values_map,omitempty"`

	// Values is the list of supported values (e.g., ["true", "false"] for booleans).
	Values []string `yaml:"values,omitempty"`

	// Range is the numeric range for integer options (e.g., min and max for `transparency`).
	Range *Range `yaml:"range,omitempty"`

	// Default is the default value of the option (e.g., "en" for `lang`, null for `background`).
	Default interface{} `yaml:"default"`

	// Validate contains validation conditions in function call style (e.g., "length 6" for `background`).
	Validate []string `yaml:"validate,omitempty"`

	// Active indicates if the option is implemented (true) or proposed (false).
	Active bool `yaml:"active"`

	// Note includes additional notes or caveats (e.g., "Proposed but not officially supported").
	Note string `yaml:"note,omitempty"`
}

// Range defines a numeric range for integer-type query options.
type Range struct {
	// Min is the minimum value of the range (e.g., 0 for `transparency`).
	Min int `yaml:"min"`

	// Max is the maximum value of the range, nullable for unbounded (e.g., null for `date`).
	Max *int `yaml:"max"`
}

// FormatSpecifier represents a single format specifier in the one-line output mode.
type FormatSpecifier struct {
	// Name is the identifier or title of the format specifier.
	Name string `yaml:"name"`

	// Letter is the character or symbol used to represent this specifier in format strings.
	Letter string `yaml:"letter"`

	// Description provides a brief explanation of what the specifier does.
	Description string `yaml:"description"`

	// Example shows a sample usage or output of the specifier.
	Example string `yaml:"example"`

	// Active indicates if the specifier is implemented (true) or proposed (false).
	Active bool `yaml:"active"`
}

// PreconfiguredFormat represents one named preset for ?format=...
type PreconfiguredFormat struct {
	ID            string `yaml:"id"` // "1", "2", "69", etc.
	Name          string `yaml:"name"`
	Format        string `yaml:"format"` // the actual string like "%c %t"
	Description   string `yaml:"description"`
	ExampleOutput string `yaml:"example_output,omitempty"`
}

// LoadSpecFromAssets loads and merges all .yaml / .yml files under spec/ in the embedded FS.
// Slices are appended; other fields are overridden by later files.
func LoadSpecFromAssets() (*WttrInOptions, error) {
	const root = "embed/spec"

	var final WttrInOptions

	// Recursive walker using only fs.ReadDirFS + fs.ReadFileFS
	var walk func(dir string) error
	walk = func(dir string) error {
		entries, err := assets.FS.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("cannot read embedded directory %q: %w", dir, err)
		}

		for _, entry := range entries {
			name := entry.Name()
			if name == "gen" {
				continue
			}

			fullPath := path.Join(dir, name)

			if entry.IsDir() {
				if err := walk(fullPath); err != nil {
					return err
				}
				continue
			}

			// Only process YAML files
			ext := strings.ToLower(filepath.Ext(name))
			if ext != ".yaml" && ext != ".yml" {
				continue
			}

			data, err := assets.FS.ReadFile(fullPath)
			if err != nil {
				return fmt.Errorf("cannot read embedded file %q: %w", fullPath, err)
			}

			var current WttrInOptions
			if err := yaml.Unmarshal(data, &current); err != nil {
				return fmt.Errorf("cannot parse YAML %q: %w", fullPath, err)
			}

			// Merge: append slices + override other fields
			if err := mergo.Merge(&final, &current,
				mergo.WithOverride,
				mergo.WithAppendSlice,
			); err != nil {
				return fmt.Errorf("merge failed for %q: %w", fullPath, err)
			}
		}
		return nil
	}

	if err := walk(root); err != nil {
		return nil, err
	}

	// Optional: detect completely empty result
	if len(final.QueryOptions) == 0 &&
		len(final.FormatSpecifiers) == 0 &&
		len(final.PreconfiguredFormats) == 0 {
		return nil, fmt.Errorf("no valid spec data found under embed/%s (no items loaded)", root)
	}

	return &final, nil
}
