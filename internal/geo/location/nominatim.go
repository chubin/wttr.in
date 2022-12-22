package location

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/chubin/wttr.in/internal/types"
	log "github.com/sirupsen/logrus"
)

type Nominatim struct {
	name  string
	url   string
	token string
}

type NominatimLocation struct {
	Name string `db:"name,key"`
	Lat  string `db:"lat"`
	Lon  string `db:"lon"`
	//nolint:tagliatelle
	Fullname string `db:"displayName" json:"display_name"`
}

func NewNominatim(name, url, token string) *Nominatim {
	return &Nominatim{
		name:  name,
		url:   url,
		token: token,
	}
}

func (n *Nominatim) Query(location string) (*Location, error) {
	var (
		result []NominatimLocation

		errResponse struct {
			Error string
		}
	)

	urlws := fmt.Sprintf(
		"%s?q=%s&format=json&accept-language=native&limit=1&key=%s",
		n.url, url.QueryEscape(location), n.token)

	log.Debugln("nominatim:", urlws)
	resp, err := http.Get(urlws)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", n.name, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", n.name, err)
	}

	err = json.Unmarshal(body, &errResponse)
	if err == nil && errResponse.Error != "" {
		return nil, fmt.Errorf("%w: %s: %s", types.ErrUpstream, n.name, errResponse.Error)
	}

	log.Debugln("nominatim: response: ", string(body))
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", n.name, err)
	}

	if len(result) != 1 {
		return nil, fmt.Errorf("%w: %s: invalid response", types.ErrUpstream, n.name)
	}

	nl := &result[0]

	return &Location{
		Lat:      nl.Lat,
		Lon:      nl.Lon,
		Fullname: nl.Fullname,
	}, nil
}
