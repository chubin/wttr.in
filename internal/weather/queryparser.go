package weather

import (
	"context"
	"fmt"
	"net/http"

	"dario.cat/mergo"

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
func (p *strictQueryParser) Parse(ctx context.Context, r *http.Request) (*query.Options, error) {
	queryString := r.URL.RawQuery

	// Step 1: use the existing parser to get validated map[string]string
	rawMap, err := options.ParseQueryString(queryString, p.config)
	if err != nil {
		return nil, err
	}

	opts, err := query.FromRequest(r)
	if err != nil {
		return nil, err
	}

	opts, err = query.ApplyParsedMap(opts, rawMap)
	if err != nil {
		return nil, err
	}

	if opts.Output == "png" {
		filenameOpts, location, err := query.ParseOptionsInFilename(opts.Location, p.config)
		if err != nil {
			return nil, err
		}

		if err := mergo.Merge(opts, *filenameOpts); err != nil {
			return nil, fmt.Errorf("filename options merge error: %w", err)
		}

		opts.Location = location

	}

	query.ApplyAutoFixes(opts)

	return opts, nil
}

// MustParse implements QueryParser.MustParse
func (p *strictQueryParser) MustParse(ctx context.Context, r *http.Request) *query.Options {
	queryString := r.URL.RawQuery
	opts, err := p.Parse(ctx, r)
	if err != nil {
		panic(fmt.Sprintf("MustParse failed: %v (query: %q)", err, queryString))
	}
	return opts
}
