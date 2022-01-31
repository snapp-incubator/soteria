package topics

import (
	"strings"
	"text/template"

	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/acl"
)

// HashType topic hashID type.
type HashType int

const (
	HashID HashType = iota
	MD5
)

type Topic struct {
	Type     string   `koanf:"type"`
	Template string   `koanf:"template"`
	HashType HashType `koanf:"hash_type"`
	Accesses map[string]acl.AccessType
}

type Template struct {
	Type     string
	Template *template.Template
	HashType HashType `koanf:"hash_type"`
	Accesses map[string]acl.AccessType
}

func (t Template) Parse(fields map[string]string) string {
	writer := new(strings.Builder)

	err := t.Template.Execute(writer, fields)
	if err != nil {
		return ""
	}

	return writer.String()
}
