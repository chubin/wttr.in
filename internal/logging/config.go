package logging

type Config struct {
	AccessLog string `yaml:"accessLog"`
	ErrorsLog string `yaml:"errorsLog"`
	Interval  int    `yaml:"interval"`
}
