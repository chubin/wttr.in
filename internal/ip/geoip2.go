package ip

import (
	"errors"
	"fmt"
	"net/netip"
	"strconv"

	"github.com/chubin/wttr.go/internal/weather"
	"github.com/oschwald/geoip2-golang/v2"
)

// IPLocatorGeoIP2 is an implementation of IPLocator using MaxMind GeoIP2 / GeoLite2 City database.
type IPLocatorGeoIP2 struct {
	db *geoip2.Reader
}

// NewIPLocatorGeoIP2 creates a new IPLocator using the given GeoIP2-City.mmdb (or GeoLite2-City.mmdb) file.
// Call Close() when you're done with it (or use it as a long-lived service component).
func NewIPLocatorGeoIP2(dbPath string) (*IPLocatorGeoIP2, error) {
	db, err := geoip2.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open GeoIP2 database: %w", err)
	}
	return &IPLocatorGeoIP2{db: db}, nil
}

// Close releases the underlying database reader resources.
func (l *IPLocatorGeoIP2) Close() error {
	if l.db != nil {
		return l.db.Close()
	}
	return nil
}

// GetIPData implements the IPLocator interface.
func (l *IPLocatorGeoIP2) GetIPData(ipStr string) (*weather.IPData, error) {
	ip, err := netip.ParseAddr(ipStr)
	if err != nil {
		return nil, fmt.Errorf("invalid IP address: %w", err)
	}

	record, err := l.db.City(ip)
	if err != nil {
		return nil, fmt.Errorf("lookup failed: %w", err)
	}

	if !record.HasData() {
		return nil, errors.New("no geolocation data found for this IP")
	}

	data := &weather.IPData{
		IP: ipStr,
	}

	// Country
	if record.Country.ISOCode != "" {
		data.CountryCode = record.Country.ISOCode
	}
	if record.Country.Names.English != "" {
		data.Country = record.Country.Names.English
	}

	// Region (first subdivision — usually state/province)
	if len(record.Subdivisions) > 0 {
		if record.Subdivisions[0].Names.English != "" {
			data.Region = record.Subdivisions[0].Names.English
		}
	}

	// City
	if record.City.Names.English != "" {
		data.City = record.City.Names.English
	}

	// Coordinates (only if present)
	if record.Location.HasCoordinates() {
		if record.Location.Latitude != nil {
			data.Latitude = strconv.FormatFloat(*record.Location.Latitude, 'f', 6, 64)
		}
		if record.Location.Longitude != nil {
			data.Longitude = strconv.FormatFloat(*record.Location.Longitude, 'f', 6, 64)
		}
	}

	return data, nil
}
