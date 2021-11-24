package tracing

type Config struct {
	Enabled bool `koanf:"enabled"`
	Agent   `koanf:"agent"`
	Ratio   float64 `koanf:"ratio"`
}

type Agent struct {
	Host string `koanf:"host"`
	Port string `koanf:"port"`
}
