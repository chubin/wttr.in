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
