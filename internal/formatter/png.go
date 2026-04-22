package formatter

import (
	"context"
	"fmt"
	"strings"

	"github.com/chubin/wttr.in/internal/domain"
	"github.com/chubin/wttr.in/internal/formatter/ansitopng"
	"github.com/chubin/wttr.in/internal/options"
	// Import your ANSI-to-PNG renderer from previous step
	// Example: "github.com/chubin/wttr.in/internal/renderer/ansitopng"
)

// PNGFormatter converts ANSI-rendered weather output into a PNG image.
type PNGFormatter struct{}

// NewPNGFormatter creates a new PNG formatter.
func NewPNGFormatter() *PNGFormatter {
	return &PNGFormatter{}
}

// Format implements the Formatter interface defined in weather.go
//
// It expects that the Renderer (usually the ANSI/text renderer) has already produced
// ANSI-colored output in renderOut.Content.
func (f *PNGFormatter) Format(
	ctx context.Context, // reserved for future cancellation / tracing
	query *domain.Query,
	renderOut *domain.RenderOutput,
) (*domain.FormatOutput, error) {
	if renderOut == nil || len(renderOut.Content) == 0 {
		return nil, fmt.Errorf("no render output provided for PNG formatting")
	}

	// Extract options that affect PNG rendering (background, transparency, inversion, etc.)
	opts := extractPNGOptions(query.Options)

	// Call your ANSI → PNG renderer (the Go port of the original Python code)
	pngBytes, err := ansitopng.RenderANSI(string(renderOut.Content), opts)
	if err != nil {
		return nil, fmt.Errorf("failed to render ANSI to PNG: %w", err)
	}

	return &domain.FormatOutput{
		Content:     pngBytes,
		ContentType: "image/png",
	}, nil
}

// extractPNGOptions converts wttr.in options into the map expected by RenderANSI.
// You can extend this as needed.
func extractPNGOptions(opts *options.Options) map[string]string {
	if opts == nil {
		return map[string]string{}
	}

	pngOpts := make(map[string]string)

	// Background color (e.g. ?background=000000)
	if opts.Background != "" {
		pngOpts["background"] = strings.TrimPrefix(opts.Background, "#")
	}

	// Transparency (e.g. ?t or ?transparency=150)
	if opts.Transparency != "" {
		pngOpts["transparency"] = opts.Transparency
	} else if opts.T { // short flag for transparency
		pngOpts["transparency"] = "150" // default value used in original wttr.in
	}

	// Inverted colors
	if opts.Inverted {
		pngOpts["inverted_colors"] = "true"
	}

	// You can add more PNG-specific options here (frame, padding, etc.)

	return pngOpts
}
