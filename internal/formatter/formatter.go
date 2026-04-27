package formatter

import "github.com/chubin/wttr.in/internal/domain"

// Formatter Implementations
type TerminalFormatter struct{}

func (f *TerminalFormatter) Format(query *domain.Query, output *domain.RenderOutput) (*domain.FormatOutput, error) {
	return &domain.FormatOutput{
		Content:     output.Content,
		ContentType: "text/plain; charset=utf-8",
	}, nil
}

type JSONFormatter struct{}

func (f *JSONFormatter) Format(query *domain.Query, output *domain.RenderOutput) (*domain.FormatOutput, error) {
	return &domain.FormatOutput{
		Content:     output.Content,
		ContentType: "application/json; charset=utf-8",
	}, nil
}

type TextFormatter struct{}

func (f *TextFormatter) Format(query *domain.Query, output *domain.RenderOutput) (*domain.FormatOutput, error) {
	return &domain.FormatOutput{
		Content:     output.Content,
		ContentType: "text/plain; charset=utf-8",
	}, nil
}
