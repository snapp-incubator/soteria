package tracing

type Config struct {
	Enabled bool `json:"enabled,omitempty" koanf:"enabled"`
	Agent   `json:"agent,omitempty"   koanf:"agent"`
	Ratio   float64 `json:"ratio,omitempty"   koanf:"ratio"`
}

type Agent struct {
	Host string `json:"host,omitempty" koanf:"host"`
	Port string `json:"port,omitempty" koanf:"port"`
}
