package location

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

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

	// Check if the location is a GPS coordinate
	if isGPSCoordinate(location) {
		// Find the nearest known location for the given GPS coordinate
		nearestLocation, err := findNearestKnownLocation(location)
		if err != nil {
			return nil, err
		}
		location = nearestLocation
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

// isGPSCoordinate checks if the given location is a GPS coordinate.
func isGPSCoordinate(location string) bool {
	parts := strings.Split(location, ",")
	if len(parts) != 2 {
		return false
	}

	_, err1 := strconv.ParseFloat(parts[0], 64)
	_, err2 := strconv.ParseFloat(parts[1], 64)

	return err1 == nil && err2 == nil
}

// findNearestKnownLocation finds the nearest known location for the given GPS coordinate.
func findNearestKnownLocation(gps string) (string, error) {
	// Implement the logic to find the nearest known location for the given GPS coordinate.
	// This is a placeholder implementation and should be replaced with the actual logic.
	return "Montreal", nil
}
