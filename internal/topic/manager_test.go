package topic_test

import (
	"testing"

	"github.com/snapp-incubator/soteria/internal/config"
	"github.com/snapp-incubator/soteria/internal/topic"
	"go.uber.org/zap"
)

// nolint: funlen
func TestTopicManager(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		arg    string
		issuer string
		want   string
	}{
		{
			name:   "testing cab event",
			arg:    "passenger-event-152384980615c2bd16143cff29038b67",
			issuer: topic.PassengerIss,
			want:   topic.CabEvent,
		},
		{
			name:   "testing cab event",
			arg:    "driver-event-152384980615c2bd16143cff29038b67",
			issuer: topic.DriverIss,
			want:   topic.CabEvent,
		},
		{
			name:   "testing invalid event",
			arg:    "-event-123456789abcdefgABCDEFG",
			issuer: topic.NoneIss,
			want:   "",
		},
		{
			name:   "testing driver location",
			arg:    "snapp/driver/DXKgaNQa7N5Y7bo/location",
			issuer: topic.DriverIss,
			want:   topic.DriverLocation,
		},
		{
			name:   "testing passenger location",
			arg:    "snapp/passenger/DXKgaNQa7N5Y7bo/location",
			issuer: topic.PassengerIss,
			want:   topic.PassengerLocation,
		},
		{
			name:   "testing invalid location",
			arg:    "snapp/thirdparty/DXKgaNQa7N5Y7bo/location",
			issuer: topic.NoneIss,
			want:   "",
		},
		{
			name:   "testing superapp event",
			arg:    "snapp/passenger/DXKgaNQa7N5Y7bo/superapp",
			issuer: topic.PassengerIss,
			want:   topic.SuperappEvent,
		},
		{
			name:   "testing shared passenger location",
			arg:    "snapp/passenger/DXKgaNQa7N5Y7bo/driver-location",
			issuer: topic.PassengerIss,
			want:   topic.SharedLocation,
		},
		{
			name:   "testing shared driver location",
			arg:    "snapp/driver/DXKgaNQa7N5Y7bo/passenger-location",
			issuer: topic.DriverIss,
			want:   topic.SharedLocation,
		},
		{
			name:   "testing passenger chat",
			arg:    "snapp/passenger/DXKgaNQa7N5Y7bo/chat",
			issuer: topic.PassengerIss,
			want:   topic.Chat,
		},
		{
			name:   "testing driver chat",
			arg:    "snapp/driver/DXKgaNQa7N5Y7bo/chat",
			issuer: topic.DriverIss,
			want:   topic.Chat,
		},
		{
			name:   "testing passenger general call entry",
			arg:    "shared/snapp/passenger/DXKgaNQa7N5Y7bo/call/send",
			issuer: topic.PassengerIss,
			want:   topic.GeneralCallEntry,
		},
		{
			name:   "testing driver general call entry",
			arg:    "shared/snapp/driver/DXKgaNQa7N5Y7bo/call/send",
			issuer: topic.DriverIss,
			want:   topic.GeneralCallEntry,
		},
		{
			name:   "testing passenger node call entry",
			arg:    "snapp/passenger/DXKgaNQa7N5Y7bo/call/heliograph-0/send",
			issuer: topic.PassengerIss,
			want:   topic.NodeCallEntry,
		},
		{
			name:   "testing driver node call entry",
			arg:    "snapp/driver/DXKgaNQa7N5Y7bo/call/heliograph-1/send",
			issuer: topic.DriverIss,
			want:   topic.NodeCallEntry,
		},
		{
			name:   "testing passenger call",
			arg:    "snapp/passenger/DXKgaNQa7N5Y7bo/call/receive",
			issuer: topic.PassengerIss,
			want:   topic.CallOutgoing,
		},
		{
			name:   "testing driver call",
			arg:    "snapp/driver/DXKgaNQa7N5Y7bo/call/receive",
			issuer: topic.DriverIss,
			want:   topic.CallOutgoing,
		},
		{
			name:   "testing box event",
			arg:    "bucks",
			issuer: topic.DriverIss,
			want:   topic.BoxEvent,
		},
	}

	cfg := config.SnappVendor()

	hid, err := topic.NewHashIDManager(cfg.HashIDMap)
	if err != nil {
		t.Errorf("invalid default hash-id: %s", err)
	}

	// nolint: exhaustruct
	topicManager := topic.NewTopicManager(cfg.Topics, hid, "snapp", cfg.IssEntityMap, cfg.IssPeerMap, zap.NewNop())

	sub := "DXKgaNQa7N5Y7bo"

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			topic := tc.arg
			topicTemplate := topicManager.ParseTopic(topic, tc.issuer, sub)
			if topicTemplate != nil {
				if len(tc.want) == 0 {
					t.Errorf("topic %s is invalid, must throw error.", tc.arg)
				} else if topicTemplate.Type != tc.want {
					t.Errorf("GetType() = %v, want %v", topicTemplate.Type, tc.want)
				}
			} else {
				if len(tc.want) != 0 {
					t.Errorf("failed to find topicTemplate for %s", tc.arg)
				}
			}
		})
	}
}