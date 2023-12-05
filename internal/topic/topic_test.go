package topic_test

import (
	"fmt"
	"testing"
	"text/template"

	"github.com/snapp-incubator/soteria/internal/topic"
	"github.com/snapp-incubator/soteria/pkg/acl"
	"github.com/stretchr/testify/require"
)

func TestCabEventTopic(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	tpc := topic.Topic{
		Type:     topic.CabEvent,
		Template: "^{{.iss}}-event-$",
		Accesses: map[string]acl.AccessType{
			topic.DriverIss:    acl.Sub,
			topic.PassengerIss: acl.Sub,
		},
	}

	temp := topic.Template{
		Type:     tpc.Type,
		Template: template.Must(template.New("").Parse(tpc.Template)),
		Accesses: tpc.Accesses,
	}

	s := temp.Parse(map[string]string{
		"iss": topic.PassengerIss,
	})

	require.Equal(fmt.Sprintf("^%s-event-$", topic.PassengerIss), s)

	require.True(temp.HasAccess(topic.PassengerIss, acl.Sub))
	require.False(temp.HasAccess(topic.PassengerIss, acl.Pub))
}

func TestChatTopic(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	tpc := topic.Topic{
		Type:     topic.Chat,
		Template: "^/{{.iss}}/{{.peer}}/chat$",
		Accesses: map[string]acl.AccessType{
			topic.DriverIss:    acl.PubSub,
			topic.PassengerIss: acl.PubSub,
		},
	}

	temp := topic.Template{
		Type:     tpc.Type,
		Template: template.Must(template.New("").Parse(tpc.Template)),
		Accesses: tpc.Accesses,
	}

	s := temp.Parse(map[string]string{
		"iss":  "passenger",
		"peer": "driver",
	})

	require.Equal("^/passenger/driver/chat$", s)

	require.True(temp.HasAccess(topic.PassengerIss, acl.Sub))
	require.True(temp.HasAccess(topic.PassengerIss, acl.Pub))
}
