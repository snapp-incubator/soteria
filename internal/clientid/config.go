package clientid

import regexp "github.com/wasilibs/go-re2"

type Config struct {
	Patterns map[string]string `json:"patterns,omitempty" koanf:"patterns"`
}

func (c Config) Regexs() map[string]*regexp.Regexp {
	regexs := make(map[string]*regexp.Regexp)

	for name, pattern := range c.Patterns {
		regexs[name] = regexp.MustCompile(pattern)
	}

	return regexs
}
