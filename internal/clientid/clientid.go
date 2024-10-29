package clientid

import regexp "github.com/wasilibs/go-re2"

type Parser struct {
	regexs map[string]*regexp.Regexp
}

func NewParser(c Config) *Parser {
	return &Parser{
		regexs: c.Regexs(),
	}
}

func (p *Parser) Parse(clientID string) string {
	for name, regex := range p.regexs {
		if regex.MatchString(clientID) {
			return name
		}
	}

	return "-"
}
