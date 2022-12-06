package ip

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseCacheEntry(t *testing.T) {
	tests := []struct {
		addr     string
		input    string
		expected Location
		err      error
	}{
		{
			"1.2.3.4",
			"DE;Germany;Free and Hanseatic City of Hamburg;Hamburg;53.5736;9.9782",
			Location{
				IP:          "1.2.3.4",
				CountryCode: "DE",
				Country:     "Germany",
				Region:      "Free and Hanseatic City of Hamburg",
				City:        "Hamburg",
				Latitude:    53.5736,
				Longitude:   9.9782,
			},
			nil,
		},

		{
			"1.2.3.4",
			"ES;Spain;Madrid, Comunidad de;Madrid;40.4165;-3.70256;28223;Orange Espagne SA;orange.es",
			Location{
				IP:          "1.2.3.4",
				CountryCode: "ES",
				Country:     "Spain",
				Region:      "Madrid, Comunidad de",
				City:        "Madrid",
				Latitude:    40.4165,
				Longitude:   -3.70256,
			},
			nil,
		},

		{
			"1.2.3.4",
			"US;United States of America;California;Mountain View",
			Location{
				IP:          "1.2.3.4",
				CountryCode: "US",
				Country:     "United States of America",
				Region:      "California",
				City:        "Mountain View",
				Latitude:    -1000,
				Longitude:   -1000,
			},
			nil,
		},

		// Invalid entries
		{
			"1.2.3.4",
			"DE;Germany;Free and Hanseatic City of Hamburg;Hamburg;53.5736;XXX",
			Location{},
			ErrInvalidCacheEntry,
		},
	}

	for _, tt := range tests {
		result, err := NewLocationFromString(tt.addr, tt.input)
		if tt.err == nil {
			require.NoError(t, err)
			require.Equal(t, *result, tt.expected)
		} else {
			require.ErrorIs(t, err, tt.err)
		}
	}
}
