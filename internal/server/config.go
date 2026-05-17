package server

type Config struct {
	PortHTTP    int    `yaml:"portHttp"`
	PortHTTPS   int    `yaml:"portHttps"`
	TLSCertFile string `yaml:"tlsCertFile"`
	TLSKeyFile  string `yaml:"tlsKeyFile"`
}
