package uplink

type Config struct {
	Address1         string `yaml:"address1"`
	Address2         string `yaml:"address2"`
	Address3         string `yaml:"address3"`
	Address4         string `yaml:"address4"`
	Timeout          int    `yaml:"timeout"`
	PrefetchInterval int    `yaml:"prefetchInterval"`
}
