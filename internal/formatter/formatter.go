package formatter

import "github.com/chubin/wttr.in/internal/domain"

// Formatter Implementations (Stubs)
type TerminalFormatter struct{}

func (f *TerminalFormatter) Format(output domain.RenderOutput) (*domain.FormatOutput, error) {
	return &domain.FormatOutput{
		Content:     output.Content,
		ContentType: "application/text",
	}, nil
}

type PNGFormatter struct{}

func (f *PNGFormatter) Format(output domain.RenderOutput) (*domain.FormatOutput, error) {
	// Stub: To be implemented
	return &domain.FormatOutput{}, nil
}

type JSONFormatter struct{}

func (f *JSONFormatter) Format(output domain.RenderOutput) (*domain.FormatOutput, error) {
	return &domain.FormatOutput{
		Content:     output.Content,
		ContentType: "application/json",
	}, nil
}

type TextFormatter struct{}

func (f *TextFormatter) Format(output domain.RenderOutput) (*domain.FormatOutput, error) {
	return &domain.FormatOutput{
		Content:     output.Content,
		ContentType: "application/text",
	}, nil
}
