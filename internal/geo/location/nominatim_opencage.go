package location

import (
	"fmt"
	"net/url"

	"github.com/chubin/wttr.in/internal/types"
)

type locationOpenCage struct {
	Results []struct {
		Name     string `db:"name,key"`
		Geometry struct {
			Lat float64 `db:"lat"`
			Lng float64 `db:"lng"`
		}
		Fullname string `json:"formatted"`
	} `json:"results"`
}

func (data *locationOpenCage) Query(n *Nominatim, location string) (*Location, error) {
	url := fmt.Sprintf(
		"%s?q=%s&language=native&limit=1&key=%s",
		n.url, url.QueryEscape(location), n.token)

	err := makeQuery(url, data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", n.name, err)
	}

	if len(data.Results) != 1 {
		return nil, fmt.Errorf("%w: %s: invalid response", types.ErrUpstream, n.name)
	}

	nl := data.Results[0]

	return &Location{
		Lat:      fmt.Sprint(nl.Geometry.Lat),
		Lon:      fmt.Sprint(nl.Geometry.Lng),
		Fullname: nl.Fullname,
	}, nil
}
