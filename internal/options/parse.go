package options

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// Errors for query parsing and validation
var (
	ErrInvalidQueryString      = errors.New("ERR001: failed to parse query string")
	ErrMultipleValues          = errors.New("ERR002: option has multiple values")
	ErrUnknownOption           = errors.New("ERR003: unknown option")
	ErrOptionNotImplemented    = errors.New("ERR004: option is not implemented")
	ErrBooleanValueRequired    = errors.New("ERR005: option requires boolean value (true/false)")
	ErrInvalidBooleanValue     = errors.New("ERR006: invalid boolean value")
	ErrValueRequired           = errors.New("ERR007: option requires a value")
	ErrInvalidStringValue      = errors.New("ERR008: invalid string value")
	ErrInvalidValidationRule   = errors.New("ERR009: invalid validation rule")
	ErrInvalidLengthRule       = errors.New("ERR010: invalid length rule")
	ErrInvalidLength           = errors.New("ERR011: invalid length")
	ErrInvalidRegexpRule       = errors.New("ERR012: invalid regexp rule")
	ErrInvalidRegexp           = errors.New("ERR013: invalid regexp match")
	ErrUnsupportedRule         = errors.New("ERR014: unsupported validation rule")
	ErrInvalidInteger          = errors.New("ERR015: option requires an integer")
	ErrBelowMinimum            = errors.New("ERR016: value below minimum")
	ErrAboveMaximum            = errors.New("ERR017: value exceeds maximum")
	ErrUnsupportedType         = errors.New("ERR018: unsupported option type")
	ErrInactiveFormatSpecifier = errors.New("ERR019: format specifier is not implemented")
)

// ParseQueryString parses a wttr.in query string into a map of option names to values.
// It validates option names, types, content, and values based on the provided WttrInOptions.
// Returns the parsed map or an error if parsing or validation fails.
func ParseQueryString(query string, config *WttrInOptions) (map[string]string, error) {
	// Initialize result map and lookup tables
	result := make(map[string]string)
	shortToName, nameToOption := buildOptionLookups(config.QueryOptions)

	// Parse query string
	query = strings.ReplaceAll(
		strings.ReplaceAll(query, "%25", "%"),
		"%", "%25")

	parsed, err := url.ParseQuery(query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidQueryString, err)
	}

	// Process query parameters
	for key, values := range parsed {
		if len(values) == 0 || (len(values) == 1 && values[0] == "") {
			if err := handleFlagOption(key, shortToName, nameToOption, result); err != nil {
				return nil, err
			}
			continue
		}
		if len(values) > 1 {
			return nil, fmt.Errorf("%w: %s: %v", ErrMultipleValues, key, values)
		}
		if err := handleValueOption(key, values[0], shortToName, nameToOption, result); err != nil {
			return nil, err
		}
	}

	// Validate format specifiers if format=custom
	if formatValue, exists := result["format"]; exists && formatValue == "custom" {
		for _, spec := range config.FormatSpecifiers {
			if !spec.Active {
				return nil, fmt.Errorf("%w: %s", ErrInactiveFormatSpecifier, spec.Specifier)
			}
		}
	}

	return result, nil
}

// buildOptionLookups creates lookup maps for short-to-long names and option details.
func buildOptionLookups(options []QueryOption) (map[string]string, map[string]QueryOption) {
	shortToName := make(map[string]string)
	nameToOption := make(map[string]QueryOption)
	for _, opt := range options {
		nameToOption[opt.Name] = opt
		if opt.Short != "" {
			shortToName[opt.Short] = opt.Name
		}
	}
	return shortToName, nameToOption
}

// handleFlagOption processes flag options (e.g., "0pq" or "m") without values.
func handleFlagOption(key string, shortToName map[string]string, nameToOption map[string]QueryOption, result map[string]string) error {
	// Handle bundled short options (e.g., "0pq")
	if len(key) > 1 && !strings.Contains(key, "=") {
		for _, short := range strings.Split(key, "") {
			if err := validateAndSetFlag(short, shortToName, nameToOption, result); err != nil {
				return err
			}
		}
		return nil
	}
	// Handle single flag
	return validateAndSetFlag(key, shortToName, nameToOption, result)
}

// validateAndSetFlag validates and sets a single flag option.
func validateAndSetFlag(key string, shortToName map[string]string, nameToOption map[string]QueryOption, result map[string]string) error {
	name, exists := shortToName[key]
	if !exists {
		name = key // Try long name
	}
	opt, exists := nameToOption[name]
	if !exists {
		return fmt.Errorf("%w: %s", ErrUnknownOption, key)
	}
	if !opt.Active {
		return fmt.Errorf("%w: %s", ErrOptionNotImplemented, name)
	}
	if opt.Type != "boolean" {
		return fmt.Errorf("%w: %s", ErrValueRequired, key)
	}
	result[name] = "true"
	return nil
}

// handleValueOption processes options with values (e.g., "lang=fr").
func handleValueOption(key, value string, shortToName map[string]string, nameToOption map[string]QueryOption, result map[string]string) error {
	// Resolve option name
	name, exists := shortToName[key]
	if !exists {
		name = key // Assume long name
	}
	opt, exists := nameToOption[name]
	if !exists {
		return fmt.Errorf("%w: %s", ErrUnknownOption, key)
	}
	if !opt.Active {
		return fmt.Errorf("%w: %s", ErrOptionNotImplemented, name)
	}

	// Validate based on type
	switch opt.Type {
	case "boolean":
		return validateBooleanOption(name, value, opt, result)
	case "string":
		return validateStringOption(name, value, opt, result)
	case "integer":
		return validateIntegerOption(name, value, opt, result)
	default:
		return fmt.Errorf("%w: %s: %s", ErrUnsupportedType, name, opt.Type)
	}
}

// validateBooleanOption validates and sets boolean option values.
func validateBooleanOption(name, value string, opt QueryOption, result map[string]string) error {
	if value != "true" && value != "false" {
		return fmt.Errorf("%w: %s: %s", ErrBooleanValueRequired, name, value)
	}
	if len(opt.Values) > 0 {
		for _, v := range opt.Values {
			if v == value {
				result[name] = value
				return nil
			}
		}
		return fmt.Errorf("%w: %s: %s, expected one of %v", ErrInvalidBooleanValue, name, value, opt.Values)
	}
	result[name] = value
	return nil
}

// validateStringOption validates and sets string option values, including background-specific rules.
func validateStringOption(name, value string, opt QueryOption, result map[string]string) error {
	if len(opt.ValuesMap) > 0 {
		if _, valid := opt.ValuesMap[value]; !valid {
			return fmt.Errorf("%w: %s: %s, expected one of %v", ErrInvalidStringValue, name, value, getMapKeys(opt.ValuesMap))
		}
	}
	if len(opt.Values) > 0 {
		for _, v := range opt.Values {
			if v == value {
				result[name] = value
				return nil
			}
		}
		return fmt.Errorf("%w: %s: %s, expected one of %v", ErrInvalidStringValue, name, value, opt.Values)
	}
	if name == "background" && len(opt.Validate) > 0 {
		for _, rule := range opt.Validate {
			if err := applyValidationRule(name, value, rule); err != nil {
				return err
			}
		}
	}
	result[name] = value
	return nil
}

// applyValidationRule applies a single validation rule (e.g., "length 6", "regexp [0-9a-fA-F]{6}").
func applyValidationRule(name, value, rule string) error {
	parts := strings.SplitN(rule, " ", 2)
	if len(parts) < 2 {
		return fmt.Errorf("%w: %s: %s", ErrInvalidValidationRule, name, rule)
	}
	switch parts[0] {
	case "length":
		length, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("%w: %s: %s", ErrInvalidLengthRule, name, rule)
		}
		if len(value) != length {
			return fmt.Errorf("%w: %s: %s, expected length %d", ErrInvalidLength, name, value, length)
		}
	case "regexp":
		re, err := regexp.Compile(parts[1])
		if err != nil {
			return fmt.Errorf("%w: %s: %s", ErrInvalidRegexpRule, name, rule)
		}
		if !re.MatchString(value) {
			return fmt.Errorf("%w: %s: %s, expected match for %s", ErrInvalidRegexp, name, value, parts[1])
		}
	default:
		return fmt.Errorf("%w: %s: %s", ErrUnsupportedRule, name, rule)
	}
	return nil
}

// validateIntegerOption validates and sets integer option values.
func validateIntegerOption(name, value string, opt QueryOption, result map[string]string) error {
	num, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("%w: %s: %s", ErrInvalidInteger, name, value)
	}
	if opt.Range != nil {
		if num < opt.Range.Min {
			return fmt.Errorf("%w: %s: %d, minimum is %d", ErrBelowMinimum, name, num, opt.Range.Min)
		}
		if opt.Range.Max != nil && num > *opt.Range.Max {
			return fmt.Errorf("%w: %s: %d, maximum is %d", ErrAboveMaximum, name, num, *opt.Range.Max)
		}
	}
	result[name] = value
	return nil
}

// getMapKeys returns a sorted slice of keys from a map for error messages.
func getMapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// Sort for consistent error messages
	for i := 0; i < len(keys)-1; i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
	return keys
}
