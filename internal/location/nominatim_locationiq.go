package location

import (
	"fmt"
	"net/url"

	"github.com/chubin/wttr.in/internal/types"
)

type locationIQ []struct {
	Name string `db:"name,key"`
	Lat  string `db:"lat"`
	Lon  string `db:"lon"`
	//nolint:tagliatelle
	Fullname string `db:"displayName" json:"display_name"`
}

func (data *locationIQ) Query(n *Nominatim, location string) (*Location, error) {
	url := fmt.Sprintf(
		"%s?q=%s&format=json&language=native&limit=1&key=%s",
		n.url, url.QueryEscape(location), n.token)

	err := makeQuery(url, data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", n.name, err)
	}

	if len(*data) != 1 {
		return nil, fmt.Errorf("%w: %s: invalid response", types.ErrUpstream, n.name)
	}

	nl := &(*data)[0]

	return &Location{
		Lat:      nl.Lat,
		Lon:      nl.Lon,
		Fullname: nl.Fullname,
	}, nil
}
