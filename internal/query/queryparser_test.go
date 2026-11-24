package query_test

import (
	"testing"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"

	"github.com/chubin/wttr.in/internal/query"
	"github.com/chubin/wttr.in/internal/options"
	"github.com/chubin/wttr.in/internal/defs"
)

func TestWeatherUnits(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		parameters   string
		UseMetric    bool
		UseImperial  bool
		UseUscs      bool
		UseMsForWind bool
	}{
		{
			parameters:   "/Kiev?u",
			UseMetric:    false,
			UseImperial:  true,
			UseUscs:      true,
			UseMsForWind: false,
		},
		{
			parameters:   "/Kiev?m",
			UseMetric:    true,
			UseImperial:  false,
			UseUscs:      false,
			UseMsForWind: false,
		},
		{
			parameters:   "/Kiev?M",
			UseMetric:    true,
			UseImperial:  false,
			UseUscs:      false,
			UseMsForWind: true,
		},
		{
			parameters:   "/Kiev?Mu",
			UseMetric:    false,
			UseImperial:  true,
			UseUscs:      true,
			UseMsForWind: true,
		},
	}

	ipOpts := &options.Options{}
	wttrOps, load_err := defs.LoadDefsFromAssets()
	assert.NoError(load_err, "LoadDefs")
	p := query.NewQueryParser(wttrOps)
	for _, tt := range tests {
		r := httptest.NewRequest(http.MethodGet, tt.parameters, nil)
		result, err := p.Parse(nil, r, ipOpts)
		assert.NoError(err, "NoError: %s", tt.parameters)
		assert.Equal(result.UseMetric, tt.UseMetric, "UseMetric: %s", tt.parameters)
		assert.Equal(result.UseImperial, tt.UseImperial, "UseImperial: %s", tt.parameters)
		assert.Equal(result.UseUscs, tt.UseUscs, "UseUscs: %s", tt.parameters)
		assert.Equal(result.UseMsForWind, tt.UseMsForWind, "UseMsForWind: %s", tt.parameters)
	}
}

func TestTransparencyFlag(t *testing.T) {
	assert := assert.New(t)
	wttrOps, load_err := defs.LoadDefsFromAssets()
	assert.NoError(load_err, "LoadDefs")

	// `t` is documented as a bare flag meaning transparency=150 (see :help and README).
	filenames := []struct {
		filename     string
		Transparency int
	}{
		{filename: "Oymyakon_tqp0.png", Transparency: 150},
		{filename: "Oymyakon_t.png", Transparency: 150},
		{filename: "Oymyakon_t=200.png", Transparency: 200},
		{filename: "Oymyakon_t_transparency=100.png", Transparency: 100},
	}
	for _, tt := range filenames {
		result, location, err := query.ParseOptionsInFilename(tt.filename, wttrOps)
		if !assert.NoError(err, "NoError: %s", tt.filename) {
			continue
		}
		assert.Equal("oymyakon", location, "location: %s", tt.filename)
		assert.Equal(tt.Transparency, result.Transparency, "Transparency: %s", tt.filename)
	}

	queries := []struct {
		parameters   string
		Transparency int
	}{
		{parameters: "/Oymyakon?t", Transparency: 150},
		{parameters: "/Oymyakon?tqp0", Transparency: 150},
		{parameters: "/Oymyakon?t&transparency=100", Transparency: 100},
		{parameters: "/Oymyakon?transparency=100&t", Transparency: 100},
	}
	ipOpts := &options.Options{}
	p := query.NewQueryParser(wttrOps)
	// Query parameters are processed in map order, which is random;
	// repeat so that both orders of t and transparency=N are exercised.
	for i := 0; i < 20; i++ {
		for _, tt := range queries {
			r := httptest.NewRequest(http.MethodGet, tt.parameters, nil)
			result, err := p.Parse(nil, r, ipOpts)
			if !assert.NoError(err, "NoError: %s", tt.parameters) {
				continue
			}
			assert.Equal(tt.Transparency, result.Transparency, "Transparency: %s", tt.parameters)
		}
	}
}
