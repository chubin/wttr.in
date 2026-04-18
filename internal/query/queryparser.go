package query

import (
	"context"
	"fmt"
	"net/http"

	"dario.cat/mergo"

	"github.com/chubin/wttr.in/internal/options"
	"github.com/chubin/wttr.in/internal/spec"
	"github.com/chubin/wttr.in/internal/weather"
)

// strictQueryParser implements QueryParser using the existing options.ParseQueryString
type strictQueryParser struct {
	config *spec.WttrInOptions
}

// NewQueryParser returns a new QueryParser that uses the provided configuration
func NewQueryParser(config *spec.WttrInOptions) weather.QueryParser {
	if config == nil {
		panic("config must not be nil")
	}
	return &strictQueryParser{config: config}
}

// Parse implements QueryParser.Parse
func (p *strictQueryParser) Parse(ctx context.Context, r *http.Request, ipOpts *options.Options) (*options.Options, error) {
	queryString := r.URL.RawQuery

	// Step 1: use the existing parser to get validated map[string]string
	rawMap, err := ParseQueryString(queryString, p.config)
	if err != nil {
		return nil, err
	}

	requestOpts, err := FromRequest(r)
	if err != nil {
		return nil, err
	}

	opts, err := options.ApplyParsedMap(requestOpts, rawMap)
	if err != nil {
		return nil, err
	}

	if opts.Output == "png" {
		filenameOpts, location, err := ParseOptionsInFilename(opts.Location, p.config)
		if err != nil {
			return nil, err
		}

		if err := mergo.Merge(opts, *filenameOpts); err != nil {
			return nil, fmt.Errorf("filename options merge error: %w", err)
		}

		opts.Location = location

	}

	// If there is no strong preference for specific units,
	// use IP-based defaults.
	if !opts.UseMetric && !opts.UseImperial && !opts.UseUscs {
		opts.UseMetric = ipOpts.UseMetric
		opts.UseImperial = ipOpts.UseImperial
		opts.UseUscs = ipOpts.UseUscs
	}

	ApplyAutoFixes(opts)

	err = Validate(opts)
	if err != nil {
		return nil, fmt.Errorf("invalid option: %w", err)
	}

	return opts, nil
}

// MustParse implements QueryParser.MustParse
func (p *strictQueryParser) MustParse(ctx context.Context, r *http.Request, ipOpts *options.Options) *options.Options {
	queryString := r.URL.RawQuery
	opts, err := p.Parse(ctx, r, ipOpts)
	if err != nil {
		panic(fmt.Sprintf("MustParse failed: %v (query: %q)", err, queryString))
	}
	return opts
}
