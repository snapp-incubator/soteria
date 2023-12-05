package topic

type HashData struct {
	Length   int    `koanf:"length"`
	Salt     string `koanf:"salt"`
	Alphabet string `koanf:"alphabet"`
}
