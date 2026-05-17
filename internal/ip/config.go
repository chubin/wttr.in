package ip

type Config struct {
	IPCacheDB   string `yaml:"ipCacheDb"`
	IPCache     string `yaml:"ipCache"`
	IPCacheType string `yaml:"ipCacheType"`
	GeoIP2      string `yaml:"geoip2"`
}
