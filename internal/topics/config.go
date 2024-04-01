package topics

type HashData struct {
	Length   int    `json:"length,omitempty"   koanf:"length"`
	Salt     string `json:"salt,omitempty"     koanf:"salt"`
	Alphabet string `json:"alphabet,omitempty" koanf:"alphabet"`
}
