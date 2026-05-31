package query

import (
	"strconv"
	"strings"

	"github.com/chubin/wttr.in/internal/options"
)

// ParseTerminalMetadata parses the terminal metadata part from User-Agent
// Example input: "term=xterm-256color; cols=220; lines=50; attached=false; color=truecolor; graphics=sixel; lang=en_US.UTF-8"
func ParseTerminalMetadata(metadata string, opts *options.Options) error {
	if metadata == "" {
		return nil
	}

	// Split by semicolon
	pairs := strings.Split(metadata, ";")

	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		keyValue := strings.SplitN(pair, "=", 2)
		if len(keyValue) != 2 {
			continue // malformed pair
		}

		key := strings.TrimSpace(keyValue[0])
		value := strings.TrimSpace(keyValue[1])

		switch key {
		case "term":
			opts.AgentTerm = value

		case "cols":
			if n, err := strconv.Atoi(value); err == nil {
				opts.AgentCols = n
			}

		case "lines":
			if n, err := strconv.Atoi(value); err == nil {
				opts.AgentLines = n
			}

		case "attached":
			opts.AgentAttached = strings.EqualFold(value, "true")

		case "color":
			opts.AgentColor = value

		case "graphics":
			opts.AgentGraphics = value

		case "lang":
			opts.AgentLang = value
		}
	}

	return nil
}
