package profiler

type Config struct {
	Enabled bool   `json:"enabled,omitempty" koanf:"enabled"`
	URL     string `json:"url,omitempty"     koanf:"url"`
}
