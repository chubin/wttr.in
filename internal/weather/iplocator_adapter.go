package weather

import (
	"fmt"

	"github.com/chubin/wttr.go/internal/ip"
)

type ipCacheLocator struct {
	cache *ip.Cache
}

func NewIPCacheLocator(cache *ip.Cache) IPLocator {
	return &ipCacheLocator{cache}
}

func (l *ipCacheLocator) GetIPData(ip string) (IPData, error) {
	addr, err := l.cache.Read(ip)
	if err != nil {
		return IPData{}, err
	}

	return IPData{
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
