package topics_test

import (
	"testing"

	"gitlab.snapp.ir/dispatching/soteria/internal/config"
	"gitlab.snapp.ir/dispatching/soteria/internal/topics"
)

// nolint: funlen
func TestTopic_GetType(t *testing.T) {
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
			issuer: topics.PassengerIss,
			want:   topics.CabEvent,
		},
		{
			name:   "testing cab event",
			arg:    "driver-event-152384980615c2bd16143cff29038b67",
			issuer: topics.DriverIss,
			want:   topics.CabEvent,
		},
		{
			name:   "testing invalid event",
			arg:    "-event-123456789abcdefgABCDEFG",
			issuer: topics.NoneIss,
			want:   "",
		},
		{
			name:   "testing driver location",
			arg:    "snapp/driver/DXKgaNQa7N5Y7bo/location",
			issuer: topics.DriverIss,
			want:   topics.DriverLocation,
		},
		{
			name:   "testing passenger location",
			arg:    "snapp/passenger/DXKgaNQa7N5Y7bo/location",
			issuer: topics.PassengerIss,
			want:   topics.PassengerLocation,
		},
		{
			name:   "testing invalid location",
			arg:    "snapp/thirdparty/DXKgaNQa7N5Y7bo/location",
			issuer: topics.NoneIss,
			want:   "",
		},
		{
			name:   "testing superapp event",
			arg:    "snapp/passenger/DXKgaNQa7N5Y7bo/superapp",
			issuer: topics.PassengerIss,
			want:   topics.SuperappEvent,
		},
		{
			name:   "testing shared passenger location",
			arg:    "snapp/passenger/DXKgaNQa7N5Y7bo/driver-location",
			issuer: topics.PassengerIss,
			want:   topics.SharedLocation,
		},
		{
			name:   "testing shared driver location",
			arg:    "snapp/driver/DXKgaNQa7N5Y7bo/passenger-location",
			issuer: topics.DriverIss,
			want:   topics.SharedLocation,
		},
		{
			name:   "testing passenger chat",
			arg:    "snapp/passenger/DXKgaNQa7N5Y7bo/chat",
			issuer: topics.PassengerIss,
			want:   topics.Chat,
		},
		{
			name:   "testing driver chat",
			arg:    "snapp/driver/DXKgaNQa7N5Y7bo/chat",
			issuer: topics.DriverIss,
			want:   topics.Chat,
		},
		{
			name:   "testing passenger general call entry",
			arg:    "shared/snapp/passenger/DXKgaNQa7N5Y7bo/call/send",
			issuer: topics.PassengerIss,
			want:   topics.GeneralCallEntry,
		},
		{
			name:   "testing driver general call entry",
			arg:    "shared/snapp/driver/DXKgaNQa7N5Y7bo/call/send",
			issuer: topics.DriverIss,
			want:   topics.GeneralCallEntry,
		},
		{
			name:   "testing passenger node call entry",
			arg:    "snapp/passenger/DXKgaNQa7N5Y7bo/call/heliograph-0/send",
			issuer: topics.PassengerIss,
			want:   topics.NodeCallEntry,
		},
		{
			name:   "testing driver node call entry",
			arg:    "snapp/driver/DXKgaNQa7N5Y7bo/call/heliograph-1/send",
			issuer: topics.DriverIss,
			want:   topics.NodeCallEntry,
		},
		{
			name:   "testing passenger call",
			arg:    "snapp/passenger/DXKgaNQa7N5Y7bo/call/receive",
			issuer: topics.PassengerIss,
			want:   topics.CallOutgoing,
		},
		{
			name:   "testing driver call",
			arg:    "snapp/driver/DXKgaNQa7N5Y7bo/call/receive",
			issuer: topics.DriverIss,
			want:   topics.CallOutgoing,
		},
		{
			name:   "testing box event",
			arg:    "bucks",
			issuer: topics.DriverIss,
			want:   topics.BoxEvent,
		},
	}

	cfg := config.SnappVendor()

	hid, err := topics.NewHashIDManager(cfg.HashIDMap)
	if err != nil {
		t.Errorf("invalid default hash-id: %s", err)
	}

	// nolint: exhaustruct
	topicManager := topics.NewTopicManager(cfg.Topics, hid, "snapp", cfg.IssEntityMap, cfg.IssPeerMap)

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
