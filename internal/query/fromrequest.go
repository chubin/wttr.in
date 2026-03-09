package query

import (
	"net/http"
	"path"
	"strings"
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
		if isPlainTextClient(userAgent) {
			opts.Output = "text"
		} else if opts.Output == "" {
			opts.Output = "html"
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
