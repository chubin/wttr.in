package oneline

import (
	"fmt"
	"log"
	"strconv"

	"github.com/itchyny/gojq"
)

// GetValue executes a jq query and returns the first result as interface{}
// Use this when you expect exactly one value (most common for weather fields)
func GetValue(data interface{}, jqQuery string) (interface{}, error) {
	query, err := gojq.Parse(jqQuery)
	if err != nil {
		return nil, fmt.Errorf("invalid jq query %q: %w", jqQuery, err)
	}

	iter := query.Run(data)
	v, ok := iter.Next()
	if !ok {
		return nil, fmt.Errorf("query returned no results: %s", jqQuery)
	}

	if err, ok := v.(error); ok {
		return nil, fmt.Errorf("query execution error: %w", err)
	}

	// Optional: check if there are more results (warn in some cases)
	if _, hasMore := iter.Next(); hasMore {
		log.Printf("Warning: query %q returned multiple results, using first", jqQuery)
	}

	return v, nil
}

// GetStringValue — convenience wrapper for the most common case (strings & numbers)
func GetStringValue(data interface{}, jqQuery string) (string, error) {
	val, err := GetValue(data, jqQuery)
	if err != nil {
		return "", err
	}

	switch v := val.(type) {
	case string:
		return v, nil
	case float64:
		// Most weather numbers are integers anyway
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case bool:
		return fmt.Sprintf("%t", v), nil
	case nil:
		return "", fmt.Errorf("value is null")
	default:
		return "", fmt.Errorf("unexpected type %T from query %q", v, jqQuery)
	}
}
