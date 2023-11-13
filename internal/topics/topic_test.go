package topics_test

import (
	"testing"
	"text/template"

	"github.com/snapp-incubator/soteria/internal/topics"
	"github.com/snapp-incubator/soteria/pkg/acl"
	"github.com/stretchr/testify/require"
)

func TestTopic(t *testing.T) {
	require := require.New(t)

	topic := topics.Topic{
		Type:     topics.CabEvent,
		Template: "^{{.iss}}-event-$",
		Accesses: map[string]acl.AccessType{
			topics.DriverIss:    acl.Sub,
			topics.PassengerIss: acl.Sub,
		},
	}

	temp := topics.Template{
		Type:     topic.Type,
		Template: template.Must(template.New("").Parse(topic.Template)),
		Accesses: topic.Accesses,
	}

	s := temp.Parse(map[string]string{
		"iss": "passenger",
	})

	require.Equal("^passenger-event-$", s)
}
