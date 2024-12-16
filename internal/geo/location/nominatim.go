package location

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/chubin/wttr.in/internal/types"
	log "github.com/sirupsen/logrus"
)

type Nominatim struct {
	name  string
	url   string
	token string
	typ   string
}

type locationQuerier interface {
	Query(*Nominatim, string) (*Location, error)
}

func NewNominatim(name, typ, url, token string) *Nominatim {
	return &Nominatim{
		name:  name,
		url:   url,
		token: token,
		typ:   typ,
	}
}

func (n *Nominatim) Query(location string) (*Location, error) {
	var data locationQuerier

	switch n.typ {
	case "iq":
		data = &locationIQ{}
	case "opencage":
		data = &locationOpenCage{}
	default:
		return nil, fmt.Errorf("%s: %w", n.name, types.ErrUnknownLocationService)
	}

	return data.Query(n, location)
}

func makeQuery(url string, result interface{}) error {
	var errResponse struct {
		Error string
	}

	log.Debugln("nominatim:", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &errResponse)
	if err == nil && errResponse.Error != "" {
		return fmt.Errorf("%w: %s", types.ErrUpstream, errResponse.Error)
	}

	log.Debugln("nominatim: response: ", string(body))
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	return nil
}
