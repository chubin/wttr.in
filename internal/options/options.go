package options

// WttrInOptions represents the configuration for wttr.in query options and format specifiers.
type WttrInOptions struct {
	QueryOptions     []QueryOption     `yaml:"query_options"`
	FormatSpecifiers []FormatSpecifier `yaml:"format_specifiers"`
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

// FormatSpecifier defines a format specifier for the `format` option (e.g., %c, %t).
type FormatSpecifier struct {
	// Specifier is the format specifier (e.g., "%c").
	Specifier string `yaml:"specifier"`

	// Description provides the output description of the specifier (e.g., "Weather condition").
	Description string `yaml:"description"`

	// Active indicates if the specifier is implemented (true) or proposed (false).
	Active bool `yaml:"active"`

	// Note contains additional notes, if any (e.g., "Proposed in issue #585").
	Note string `yaml:"note,omitempty"`
}
