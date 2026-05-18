package location

type Config struct {
	LocationCacheType string `yaml:"locationCacheType"`
	LocationCacheDB   string `yaml:"locationCacheDb"`
	LocationCache     string `yaml:"locationCache"`
	NominatimServers  []struct {
		Name  string `yaml:"name"`
		Type  string `yaml:"type"`
		URL   string `yaml:"url"`
		Token string `yaml:"token"`
	} `yaml:"nominatim"`
}
