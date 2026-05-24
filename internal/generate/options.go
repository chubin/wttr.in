// Package generate creates query.Options struct + mapping logic from YAML spec.
package generate

import (
	"fmt"
	"os"
	"text/template"
	"unicode"

	"gopkg.in/yaml.v3"

	"github.com/chubin/wttr.in/internal/assets"
)

var (
	templateFile = "share/templates/options.go.tmpl"
	configFile   = "share/defs/options/options.yaml"
	outputFile   = "internal/options/options.go"
)

// OptionConfig mirrors the YAML structure (adjusted for generation needs)
type OptionConfig struct {
	Name        string            `yaml:"name"`
	Short       string            `yaml:"short,omitempty"`
	Description string            `yaml:"description"`
	Type        string            `yaml:"type"`
	Default     interface{}       `yaml:"default"`
	Active      bool              `yaml:"active"`
	Note        string            `yaml:"note,omitempty"`
	Values      []string          `yaml:"values,omitempty"` // simplified
	ValuesMap   map[string]string `yaml:"values_map,omitempty"`
	Range       struct {
		Min interface{} `yaml:"min"`
		Max interface{} `yaml:"max"`
	} `yaml:"range,omitempty"`
	Validate []string `yaml:"validate,omitempty"`
}

// OptionsConfig is the root of the YAML
type OptionsConfig struct {
	QueryOptions     []OptionConfig    `yaml:"query_options"`
	FormatSpecifiers []FormatSpecifier `yaml:"format_specifiers"`
}

type FormatSpecifier struct {
	Specifier   string `yaml:"specifier"`
	Description string `yaml:"description"`
	Active      bool   `yaml:"active"`
	Note        string `yaml:"note,omitempty"`
}

// GeneratedField holds template data for one field
type GeneratedField struct {
	Name        string
	FieldName   string
	Description string
	FieldType   string
	Note        string
	Default     interface{}
}

// toGoType maps YAML type → Go type
func toGoType(typ string) string {
	switch typ {
	case "boolean":
		return "bool"
	case "integer":
		return "int"
	case "string":
		return "string"
	default:
		return "string" // fallback – be conservative
	}
}

// toFieldName turns kebab/snake-case → PascalCase
func toFieldName(name string) string {
	if name == "" {
		return ""
	}
	var result []rune
	upperNext := true
	for _, r := range name {
		if r == '-' || r == '_' {
			upperNext = true
			continue
		}
		if upperNext {
			result = append(result, unicode.ToUpper(r))
			upperNext = false
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// GenerateOptionsAndParser is the main entry point
func GenerateOptionsAndParser() error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	var cfg OptionsConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("unmarshal yaml: %w", err)
	}

	var fields []GeneratedField
	for _, opt := range cfg.QueryOptions {
		if !opt.Active {
			continue
		}
		fields = append(fields, GeneratedField{
			Name:        opt.Name,
			FieldName:   toFieldName(opt.Name),
			Description: opt.Description,
			FieldType:   toGoType(opt.Type),
			Note:        opt.Note,
			Default:     opt.Default,
		})
	}

	tmplSrc, err := assets.GetFile(templateFile)
	if err != nil {
		return fmt.Errorf("loading template from assets: %w", err)
	}

	tmpl, err := template.New("query_options").Parse(string(tmplSrc))
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	f, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer f.Close()

	return tmpl.Execute(f, struct {
		QueryOptions []GeneratedField
	}{
		QueryOptions: fields,
	})
}
