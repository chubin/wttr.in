package renderer

import "github.com/chubin/wttr.in/internal/renderer/subprocess"

// Config holds configuration for all renderers that need external config.
type Config struct {
	Subprocess []subprocess.SubprocessRoute `yaml:"subprocess"`
}
