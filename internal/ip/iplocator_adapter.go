package ip

import (
	"fmt"

	"github.com/chubin/wttr.in/internal/weather"
)

type ipCacheLocator struct {
	cache *Cache
}

func NewIPCacheLocator(cache *Cache) weather.IPLocator {
	return &ipCacheLocator{cache}
}

func (l *ipCacheLocator) GetIPData(ip string) (*weather.IPData, error) {
	addr, err := l.cache.Read(ip)
	if err != nil {
		return &weather.IPData{}, err
	}

	return &weather.IPData{
		IP:          addr.IP,
		CountryCode: addr.CountryCode,
		Country:     addr.Country,
		Region:      addr.Region,
		City:        addr.City,
		Latitude:    fmt.Sprint(addr.Latitude),
		Longitude:   fmt.Sprint(addr.Longitude),
		// FullAddress:  fmt.Sprintf("%s, %s, %s", addr.City, addr.Region, addr.Country),
	}, nil
}
