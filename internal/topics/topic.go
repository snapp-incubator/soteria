package topics

import (
	"regexp"
	"strings"
	"text/template"
)

type Topic struct {
	Type     string `koanf:"type"`
	Template string `koanf:"template"`
	Regex    string `koanf:"regex"`
}

type Template struct {
	Type     string
	Template *template.Template
	Regex    *regexp.Regexp
}

func (t Template) Parse(fields map[string]string) string {
	writer := new(strings.Builder)

	err := t.Template.Execute(writer, fields)
	if err != nil {
		return ""
	}

	return writer.String()
}
