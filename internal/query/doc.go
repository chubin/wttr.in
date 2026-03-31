// Package query provides configuration-driven parsing and validation of
// wttr.in-style HTTP query strings.
//
// The package is built around a YAML specification file (typically
// spec/options/options.yaml) that defines all supported query parameters,
// their types, allowed values, ranges, validation rules, default values,
// implementation status (active/proposed), and – for the special "format"
// parameter – the list of supported format specifiers (%c, %t, %w, etc.).
//
// Main features:
//   - Strict validation of parameter names, types and values
//   - Support for short flags (both single and bundled: ?0pq)
//   - Boolean flags without values (?T → T=true)
//   - Special rules for certain parameters (e.g. background must be 6-char hex)
//   - Distinction between active (implemented) and inactive (proposed) options
//   - Batch validation of historical access logs to detect invalid or
//     deprecated usage patterns
//
// Typical usage:
//
//	parsed, err := query.ParseQueryString("lang=fr&format=3&T&0pq", config)
//	if err != nil { ... }  // err contains detailed reason
//
//	// or process access log to find bad queries:
//	err = query.ProcessLogFile("access.log", "invalid.log", config)
//
// The package is primarily intended for:
//   - wttr.in backend / proxy implementations
//   - API gateways or WAF rules that want to reject malformed probes
//   - Log analysis tools that monitor usage of undocumented parameters
//   - Tools that generate up-to-date documentation / CLI help from the YAML spec
//
// # Compatibility note
//
// This implementation is stricter than the original wttr.in service in several
// ways (case sensitivity, rejection of unknown parameters, no support for
// comma-separated lang lists, strict format specifier check when format=custom,
// etc.).
package query
