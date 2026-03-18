// Package spec provides structured access to wttr.in's one-line format specification,
// particularly the list of supported placeholders (%c, %t, etc.).
package spec

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/chubin/wttr.go/internal/assets"
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

// NewFromAssets reads a QueryOption from an embeeded YAML file and returns a pointer to it.
func NewFromAssets() (*WttrInOptions, error) {
	optionsSpecFile := "spec/options/options.yaml"
	onelineSpecFile := "spec/oneline/oneline.yaml"

	// Read Options description
	data, err := assets.GetFile(optionsSpecFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Unmarshal the YAML content into a QueryOption struct
	var option WttrInOptions
	if err := yaml.Unmarshal(data, &option); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	// Read Format specifiers
	data, err = assets.GetFile(onelineSpecFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Unmarshal the YAML content into a QueryOption struct
	var oneline WttrInOptions
	if err := yaml.Unmarshal(data, &option); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	option.FormatSpecifiers = oneline.FormatSpecifiers

	return &option, nil
}
