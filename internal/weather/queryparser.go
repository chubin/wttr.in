package weather

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/chubin/wttr.go/internal/options"
	"github.com/chubin/wttr.go/internal/query"
)

// strictQueryParser implements QueryParser using the existing options.ParseQueryString
type strictQueryParser struct {
	config *options.WttrInOptions
}

// NewQueryParser returns a new QueryParser that uses the provided configuration
func NewQueryParser(config *options.WttrInOptions) QueryParser {
	if config == nil {
		panic("config must not be nil")
	}
	return &strictQueryParser{config: config}
}

// Parse implements QueryParser.Parse
func (p *strictQueryParser) Parse(ctx context.Context, queryString string) (*query.Options, error) {
	// Step 1: use the existing parser to get validated map[string]string
	rawMap, err := options.ParseQueryString(queryString, p.config)
	if err != nil {
		return nil, err
	}

	// Step 2: convert map → JSON → query.Options struct
	jsonBytes, err := json.Marshal(rawMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal parsed options to JSON: %w", err)
	}

	var opts query.Options
	decoder := json.NewDecoder(bytes.NewReader(jsonBytes))
	decoder.DisallowUnknownFields() // optional: fail if map has keys not in struct
	if err := decoder.Decode(&opts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal into query.Options: %w", err)
	}

	return &opts, nil
}

// MustParse implements QueryParser.MustParse
func (p *strictQueryParser) MustParse(ctx context.Context, query string) *query.Options {
	opts, err := p.Parse(ctx, query)
	if err != nil {
		panic(fmt.Sprintf("MustParse failed: %v (query: %q)", err, query))
	}
	return opts
}
