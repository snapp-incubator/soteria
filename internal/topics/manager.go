package topics

import "regexp"

type Manager struct {
	regexes map[string]*regexp.Regexp
}

func NewManager(topics map[string]string) Manager {
	regexes := make(map[string]*regexp.Regexp)

	for key, template := range topics {
		regexes[key] = regexp.MustCompile(template)
	}

	return Manager{regexes: regexes}
}

func (m Manager) IsTopicValid(topic string) bool {

	for _, regex := range m.regexes {
		if regex.MatchString(topic) {
			return true
		}
	}

	return false
}
