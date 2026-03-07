package server

type Config struct {
	PortHttp    int    `yaml:"portHttp"`
	PortHttps   int    `yaml:"portHttps"`
	TLSCertFile string `yaml:"tlsCertFile"`
	TLSKeyFile  string `yaml:"tlsKeyFile"`
}
