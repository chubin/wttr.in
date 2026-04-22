package formatter

import (
	"context"
	"fmt"

	"github.com/chubin/wttr.in/internal/domain"
	"github.com/chubin/wttr.in/internal/formatter/ansitopng"
)

// PNGFormatter converts ANSI-rendered weather output into a PNG image.
type PNGFormatter struct{}

// NewPNGFormatter creates a new PNG formatter.
func NewPNGFormatter() *PNGFormatter {
	return &PNGFormatter{}
}

// Format implements the Formatter interface defined in weather.go
func (f *PNGFormatter) Format(
	ctx context.Context,
	query *domain.Query,
	renderOut *domain.RenderOutput,
) (*domain.FormatOutput, error) {
	if renderOut == nil || len(renderOut.Content) == 0 {
		return nil, fmt.Errorf("no render output provided for PNG formatting")
	}

	// Convert domain options to PNG-specific options (strongly typed)
	pngOpts := ansitopng.FromDomainOptions(query.Options)

	// Call the ANSI → PNG renderer
	pngBytes, err := ansitopng.RenderANSI(string(renderOut.Content), pngOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to render ANSI to PNG: %w", err)
	}

	return &domain.FormatOutput{
		Content:     pngBytes,
		ContentType: "image/png",
	}, nil
}
