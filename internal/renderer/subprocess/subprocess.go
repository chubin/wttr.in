// Package subprocess provides a Renderer that delegates rendering to an
// external subprocess.
//
// Data is passed exclusively via environment variables (WTTR_* prefix).
// Routing is configured via a list of SubprocessRoute rules (first match wins).
// Location (options.Location) is the primary selector; any option from
// options.ToMap() can be used for additional filtering.
//
// Which parts of domain.Query are exported is declared per-route.
package subprocess

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/chubin/wttr.in/internal/domain"
)

// SubprocessRoute defines a single routing rule for the subprocess renderer.
type SubprocessRoute struct {
	// Location matches query.Options.Location exactly.
	// Use "*" or leave empty to match any location.
	Location string `yaml:"location,omitempty"`

	// Selectors are additional exact matches against options.ToMap().
	// All listed keys must match for the rule to apply.
	Selectors map[string]string `yaml:"selectors,omitempty"`

	// Program is the absolute or $PATH-resolved executable to run.
	Program string `yaml:"program"`

	// Data lists which parts of the domain.Query should be serialized
	// into WTTR_* environment variables.
	// Supported values (case-insensitive):
	//   - "options"        → WTTR_OPTION_*
	//   - "location"       → WTTR_LOCATION_*
	//   - "ipdata", "ip"   → WTTR_IP_*
	//   - "clientdata", "client" → WTTR_CLIENT_*
	//   - "weather", "weatherraw" → WTTR_WEATHER_RAW (base64)
	Data []string `yaml:"data"`
}

// Renderer implements weather.Renderer by executing an external subprocess.
type Renderer struct {
	routes []SubprocessRoute
}

// NewRenderer creates a new subprocess renderer with the given routing rules.
// The slice order determines matching priority (first match wins).
func NewRenderer(routes []SubprocessRoute) *Renderer {
	if routes == nil {
		routes = []SubprocessRoute{}
	}
	return &Renderer{routes: routes}
}

// Render implements the Renderer interface.
func (r *Renderer) Render(query domain.Query) (domain.RenderOutput, error) {
	if query.Options == nil {
		return domain.RenderOutput{}, fmt.Errorf("subprocess renderer: query.Options is required")
	}

	loc := query.Options.Location
	optsMap := query.Options.ToMap()

	for _, route := range r.routes {
		if !r.matchesLocation(route, loc) {
			continue
		}
		if !r.matchesSelectors(route.Selectors, optsMap) {
			continue
		}

		// first matching route wins
		return r.execute(route, &query)
	}

	return domain.RenderOutput{}, fmt.Errorf("subprocess renderer: no route matched for location %q", loc)
}

func (r *Renderer) matchesLocation(route SubprocessRoute, loc string) bool {
	if route.Location == "" || route.Location == "*" {
		return true
	}
	return route.Location == loc
}

func (r *Renderer) matchesSelectors(selectors map[string]string, opts map[string]string) bool {
	for k, want := range selectors {
		if got, ok := opts[k]; !ok || got != want {
			return false
		}
	}
	return true
}

// execute runs the external program with the selected data in env vars.
func (r *Renderer) execute(route SubprocessRoute, q *domain.Query) (domain.RenderOutput, error) {
	cmd := exec.Command(route.Program)

	// Start with current environment and append WTTR_* variables
	env := os.Environ()

	for _, datum := range route.Data {
		switch strings.ToLower(datum) {
		case "options":
			if q.Options != nil {
				m := q.Options.ToMap()
				for k, v := range m {
					key := "WTTR_OPTION_" + strings.ToUpper(strings.ReplaceAll(k, "-", "_"))
					env = append(env, key+"="+v)
				}
			}

		case "location":
			if q.Location != nil {
				l := q.Location
				env = append(env,
					fmt.Sprintf("WTTR_LOCATION_NAME=%s", l.Name),
					fmt.Sprintf("WTTR_LOCATION_COUNTRY=%s", l.Country),
					fmt.Sprintf("WTTR_LOCATION_COUNTRY_CODE=%s", l.CountryCode),
					fmt.Sprintf("WTTR_LOCATION_LATITUDE=%.6f", l.Latitude),
					fmt.Sprintf("WTTR_LOCATION_LONGITUDE=%.6f", l.Longitude),
					fmt.Sprintf("WTTR_LOCATION_FULL_ADDRESS=%s", l.FullAddress),
					fmt.Sprintf("WTTR_LOCATION_TIMEZONE=%s", l.TimeZone),
				)
			}

		case "ipdata", "ip":
			if q.IPData != nil {
				i := q.IPData
				env = append(env,
					fmt.Sprintf("WTTR_IP=%s", i.IP),
					fmt.Sprintf("WTTR_IP_COUNTRY=%s", i.Country),
					fmt.Sprintf("WTTR_IP_COUNTRY_CODE=%s", i.CountryCode),
					fmt.Sprintf("WTTR_IP_REGION=%s", i.Region),
					fmt.Sprintf("WTTR_IP_CITY=%s", i.City),
					fmt.Sprintf("WTTR_IP_LATITUDE=%s", i.Latitude),
					fmt.Sprintf("WTTR_IP_LONGITUDE=%s", i.Longitude),
				)
			}

		case "clientdata", "client":
			if q.ClientData != nil {
				c := q.ClientData
				env = append(env,
					fmt.Sprintf("WTTR_CLIENT_IP=%s", c.ClientIP),
					fmt.Sprintf("WTTR_CLIENT_AGENT=%s", c.ClientAgent),
				)
			}

		case "weather", "weatherraw":
			if q.Weather != nil && len(*q.Weather) > 0 {
				b64 := base64.StdEncoding.EncodeToString(*q.Weather)
				env = append(env, "WTTR_WEATHER_RAW="+b64)
			}
		}
	}

	cmd.Env = env

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return domain.RenderOutput{}, fmt.Errorf("subprocess %q failed (exit code %d): %w\nstderr:\n%s",
			route.Program, cmd.ProcessState.ExitCode(), err, stderr.String())
	}

	return domain.RenderOutput{Content: stdout.Bytes()}, nil
}
