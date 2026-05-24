// Package generate creates query.Options struct + mapping logic from YAML spec.
package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"gopkg.in/yaml.v3"

	"github.com/chubin/wttr.in/internal/assets"
)

var (
	templateFile = "share/templates/options.go.tmpl"
	configDir    = "share/defs/options"
	outputFile   = "internal/options/options.go"
)

// OptionConfig mirrors the YAML structure
type OptionConfig struct {
	Name        string            `yaml:"name"`
	Short       string            `yaml:"short,omitempty"`
	Description string            `yaml:"description"`
	Type        string            `yaml:"type"`
	Default     interface{}       `yaml:"default"`
	Active      bool              `yaml:"active"`
	Note        string            `yaml:"note,omitempty"`
	Values      []string          `yaml:"values,omitempty"`
	ValuesMap   map[string]string `yaml:"values_map,omitempty"`
	Range       struct {
		Min interface{} `yaml:"min"`
		Max interface{} `yaml:"max"`
	} `yaml:"range,omitempty"`
	Validate []string `yaml:"validate,omitempty"`
}

// OptionsConfig is the root of each YAML file
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
		return "string"
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

// loadAndMergeOptions reads all .yaml/.yml files in the directory and merges them
func loadAndMergeOptions(dir string) (OptionsConfig, error) {
	var merged OptionsConfig

	entries, err := os.ReadDir(dir)
	if err != nil {
		return merged, fmt.Errorf("read options directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
			continue
		}

		path := filepath.Join(dir, name)

		data, err := os.ReadFile(path)
		if err != nil {
			return merged, fmt.Errorf("read file %s: %w", path, err)
		}

		var cfg OptionsConfig
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return merged, fmt.Errorf("unmarshal %s: %w", path, err)
		}

		// Append lists from each file (as requested)
		merged.QueryOptions = append(merged.QueryOptions, cfg.QueryOptions...)
		merged.FormatSpecifiers = append(merged.FormatSpecifiers, cfg.FormatSpecifiers...)
	}

	if len(merged.QueryOptions) == 0 && len(merged.FormatSpecifiers) == 0 {
		return merged, fmt.Errorf("no options found in directory %s", dir)
	}

	return merged, nil
}

// GenerateOptionsAndParser is the main entry point
func GenerateOptionsAndParser() error {
	cfg, err := loadAndMergeOptions(configDir)
	if err != nil {
		return err
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
