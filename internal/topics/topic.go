package topics

import (
	"strings"
	"text/template"

	"github.com/snapp-incubator/soteria/pkg/acl"
)

type Topic struct {
	Type     string                    `json:"type,omitempty"     koanf:"type"`
	Template string                    `json:"template,omitempty" koanf:"template"`
	Accesses map[string]acl.AccessType `json:"accesses,omitempty" koanf:"accesses"`
}

type Template struct {
	Type     string
	Template *template.Template
	Accesses map[string]acl.AccessType
}

func (t Template) Parse(fields map[string]string) string {
	writer := new(strings.Builder)

	if err := t.Template.Execute(writer, fields); err != nil {
		return ""
	}

	return writer.String()
}

// HasAccess check if user has access on topic.
func (t Template) HasAccess(iss string, accessType acl.AccessType) bool {
	access := t.Accesses[iss]

	return access == acl.PubSub || access == accessType
}
