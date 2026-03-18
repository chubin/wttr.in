package query

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/chubin/wttr.go/internal/options"
	"github.com/chubin/wttr.go/internal/spec"
)

// FromRequest creates an Options struct based on the provided HTTP request.
// It processes the URL path, domain name, and headers according to the specified rules.
func FromRequest(r *http.Request) (*Options, error) {
	opts := &Options{}

	// Step 1: Extract components from the URL path
	urlPath := r.URL.Path
	// Remove leading and trailing slashes for consistency
	urlPath = strings.Trim(urlPath, "/")

	// Split the path into components
	components := strings.Split(urlPath, "/")
	if len(components) == 0 {
		components = []string{""}
	}

	// Step 2: Process domain name for VIEW or LANG
	host := r.Host
	domainParts := strings.SplitN(host, ".", 2)
	if len(domainParts) > 0 {
		domainPrefix := strings.ToLower(domainParts[0])
		// Check if the domain prefix is a valid language code (e.g., "de" for German)
		// Assuming a simple check for 2-letter codes for languages
		if len(domainPrefix) == 2 && isValidLanguageCode(domainPrefix) {
			opts.Lang = domainPrefix
		} else if isValidView(domainPrefix) {
			// Check if it's a valid view name (e.g., "v2")
			opts.View = domainPrefix
		}
	}

	// Step 3: Process the first path component for VIEW or LANG if not set by domain
	firstComponent := ""
	if len(components) > 0 && components[0] != "" {
		firstComponent = strings.ToLower(components[0])
		if isValidLanguageCode(firstComponent) && opts.Lang == "" {
			opts.Lang = firstComponent
		} else if isValidView(firstComponent) && opts.View == "" {
			opts.View = firstComponent
		}
	}

	// Step 4: Extract location and format from the path
	location := ""
	if len(components) > 1 && components[1] != "" {
		location = components[1]
	} else if len(components) == 1 && firstComponent != "" && !isValidLanguageCode(firstComponent) && !isValidView(firstComponent) {
		location = components[0]
	}

	// Check for file extension in the location to determine output format
	if location != "" {
		ext := path.Ext(location)
		if ext != "" {
			ext = strings.ToLower(strings.TrimPrefix(ext, "."))
			switch ext {
			case "png", "jpg", "jpeg":
				opts.Output = ext
			case "json":
				opts.Output = "json"
			}
			// Remove the extension from the location
			location = strings.TrimSuffix(location, "."+ext)
		}

		opts.Location = location
	}
	opts.Location = strings.ReplaceAll(opts.Location, "+", " ")
	opts.Location = strings.TrimPrefix(opts.Location, "~")

	// Step 5: Set default language from header if not set by domain or URL
	if opts.Lang == "" {
		acceptLang := r.Header.Get("Accept-Language")
		if acceptLang != "" {
			// Take the first language from the Accept-Language header
			langs := strings.Split(acceptLang, ",")
			if len(langs) > 0 {
				lang := strings.Split(langs[0], ";")[0]
				lang = strings.TrimSpace(strings.ToLower(lang))
				if len(lang) >= 2 {
					lang = lang[:2] // Take first two letters (e.g., "en" from "en-US")
					if isValidLanguageCode(lang) {
						opts.Lang = lang
					}
				}
			}
		}
	}
	// Default language if still unset
	if opts.Lang == "" {
		opts.Lang = "en"
	}

	// Step 6: Set User-Agent from header
	userAgent := r.Header.Get("User-Agent")
	if userAgent != "" {
		opts.Agent = userAgent
		// Determine if it's a plain text client (e.g., curl, wget)
		if opts.Output == "" {
			if isPlainTextClient(userAgent) {
				opts.Output = "text"
			} else {
				opts.Output = "html"
			}
		}
	}

	// Step 8: Set default metric units (overridden by explicit options if present)
	opts.Metric = true

	// Step 9: Set default transparency for PNG output if applicable
	if opts.Output == "png" && opts.Transparency == 0 {
		opts.Transparency = 150
	}

	return opts, nil
}

func ApplyAutoFixes(opts *Options) {
	if opts.View == "" {
		if opts.Format == "j1" || opts.Format == "j2" {
			opts.View = opts.Format
		} else if opts.Format != "" {
			opts.View = "line"
		} else {
			opts.View = "v1"
		}
	}
	if opts.View == "j1" || opts.View == "j2" {
		opts.Output = "json"
	}
}

// ParseOptionsInFilename converts a wttr.in-style PNG filename into *Options
// using the same parsing/validation pipeline as normal ?query=string requests.
//
// Examples:
//
//	"Paris.png"                  → location=Paris
//	"Moscow_200x_m_q.png"        → location=Moscow, use_metric=true, no_caption=true, width=200
//	"Rome_0pq_lang=it_T.png"     → days=0, padding=true, no-caption=true, lang=it, no-terminal=true
//	"Berlin_u_300x150.png"       → use_imperial=true, width=300, height=150
//
// Returns location separately so caller can decide what to do with it.
func ParseOptionsInFilename(filename string, cfg *spec.WttrInOptions) (*Options, string, error) {
	if cfg == nil {
		return nil, "", fmt.Errorf("config required")
	}

	// 1. Normalize filename
	name := strings.TrimSuffix(strings.ToLower(filename), ".png")
	name = strings.TrimSuffix(name, ".PNG") // just in case

	if name == "" {
		return &Options{Output: "png"}, "", nil
	}

	// 2. Split into location + option parts
	parts := strings.Split(name, "_")
	if len(parts) == 0 {
		return &Options{Output: "png"}, "", nil
	}

	location := parts[0]
	location = strings.ReplaceAll(location, "+", " ") // standard wttr.in normalization

	// 3. Build fake query string from the remaining parts
	var qParams url.Values = make(url.Values)

	for _, part := range parts[1:] {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Special case: dimension tokens (most common in PNG URLs)
		if strings.Contains(part, "x") {
			dims := strings.SplitN(part, "x", 2)
			widthStr, heightStr := "", ""

			if len(dims) == 2 {
				widthStr = dims[0]
				heightStr = dims[1]
			} else if strings.HasPrefix(part, "x") {
				heightStr = part[1:]
			} else {
				widthStr = part[:len(part)-1] // 300x → width=300
			}

			if widthStr != "" {
				qParams.Add("width", widthStr) // assuming you support width=... (add to config if needed)
			}
			if heightStr != "" {
				qParams.Add("height", heightStr)
			}
			continue
		}

		// key=value long option
		if strings.Contains(part, "=") {
			kv := strings.SplitN(part, "=", 2)
			if len(kv) == 2 && kv[0] != "" {
				key := strings.ToLower(kv[0])
				val := kv[1]

				// Some common aliases / normalizations
				if key == "t" {
					key = "transparency"
				}
				if key == "view" || key == "format" {
					key = "view" // normalize
				}

				qParams.Add(key, val)
			}
			continue
		}

		// Otherwise: bundle of short flags (mMuIpq etc.)
		for _, ch := range part {
			if ch == 0 {
				continue
			}
			qParams.Add(string(ch), "") // empty value = flag
		}
	}

	// 4. Turn into query string and parse with the real parser
	queryStr := qParams.Encode()

	rawMap, err := options.ParseQueryString(queryStr, cfg)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse converted PNG query: %w (query was: %s)", err, queryStr)
	}

	// 5. Build Options the standard way
	opts := &Options{
		Output:       "png",
		Location:     location,
		Transparency: 150,  // PNG default
		Metric:       true, // global default
	}

	opts, err = ApplyParsedMap(opts, rawMap)
	if err != nil {
		return nil, "", err
	}

	return opts, location, nil
}

// isValidLanguageCode checks if the provided code is a valid 2-letter language code.
// This is a simplified check; in a real application, you'd have a list of supported languages.
func isValidLanguageCode(code string) bool {
	// Simplified: assume any 2-letter lowercase code is a language
	return len(code) == 2 && strings.ToLower(code) == code
}

// isValidView checks if the provided string is a valid view name.
// This is a placeholder; in a real application, you'd have a list of supported views.
func isValidView(view string) bool {
	// Example: support "v2", "j1", "j2" as valid views
	return view == "v2" || view == "j1" || view == "j2"
}

// isPlainTextClient determines if the User-Agent indicates a plain text client like curl or wget.
func isPlainTextClient(userAgent string) bool {
	// plainTextAgents contains signatures of the plain-text agents.
	plainTextAgents := []string{
		"curl",
		"httpie",
		"lwp-request",
		"wget",
		"python-httpx",
		"python-requests",
		"openbsd ftp",
		"powershell",
		"fetch",
		"aiohttp",
		"http_get",
		"xh",
		"nushell",
		"zig",
	}

	ua := strings.ToLower(userAgent)
	for _, textAgentSignature := range plainTextAgents {
		if strings.Contains(ua, textAgentSignature) {
			return true
		}
	}

	return false
}
