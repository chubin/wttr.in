package formatter

import "github.com/chubin/wttr.in/internal/domain"

// Formatter Implementations (Stubs)
type TerminalFormatter struct{}

func (f *TerminalFormatter) Format(query *domain.Query, output *domain.RenderOutput) (*domain.FormatOutput, error) {
	return &domain.FormatOutput{
		Content:     output.Content,
		ContentType: "application/text",
	}, nil
}

type JSONFormatter struct{}

func (f *JSONFormatter) Format(query *domain.Query, output *domain.RenderOutput) (*domain.FormatOutput, error) {
	return &domain.FormatOutput{
		Content:     output.Content,
		ContentType: "application/json",
	}, nil
}

type TextFormatter struct{}

func (f *TextFormatter) Format(query *domain.Query, output *domain.RenderOutput) (*domain.FormatOutput, error) {
	return &domain.FormatOutput{
		Content:     output.Content,
		ContentType: "application/text",
	}, nil
}
