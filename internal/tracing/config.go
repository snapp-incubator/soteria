package tracing

type Config struct {
	Enabled  bool    `json:"enabled,omitempty"  koanf:"enabled"`
	Endpoint string  `json:"endpoint,omitempty" koanf:"endpoint"`
	Ratio    float64 `json:"ratio,omitempty"    koanf:"ratio"`
}
