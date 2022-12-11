package location

import (
	"encoding/json"
	"log"

	"github.com/chubin/wttr.in/internal/config"
)

type Location struct {
	Name     string `db:"name,key"`
	Lat      string `db:"lat"`
	Lon      string `db:"lon"`
	Timezone string `db:"timezone"`
	//nolint:tagliatelle
	Fullname string `db:"displayName" json:"display_name"`
}

// String returns string representation of location.
func (l *Location) String() string {
	bytes, err := json.Marshal(l)
	if err != nil {
		// should never happen
		log.Fatalln(err)
	}

	return string(bytes)
}

type Provider interface {
	Query(location string) (*Location, error)
}

type Searcher struct {
	providers []Provider
}

// NewSearcher returns a new Searcher for the specified config.
func NewSearcher(config *config.Config) *Searcher {
	providers := []Provider{}
	for _, p := range config.Geo.Nominatim {
		providers = append(providers, NewNominatim(p.Name, p.URL, p.Token))
	}

	return &Searcher{
		providers: providers,
	}
}

// Search makes queries through all known providers,
// and returns response, as soon as it is not nil.
// If all responses were nil, the last response is returned.
func (s *Searcher) Search(location string) (*Location, error) {
	var (
		err    error
		result *Location
	)

	for _, p := range s.providers {
		result, err = p.Query(location)
		if result != nil && err == nil {
			return result, nil
		}
	}

	return result, err
}
