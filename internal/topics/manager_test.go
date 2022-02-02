package topics_test

import (
	"fmt"
	"gitlab.snapp.ir/dispatching/snappids/v2"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/authenticator"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/config"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/user"
	"testing"
)

func TestTopic_GetType(t1 *testing.T) {
	tests := []struct {
		name   string
		arg    string
		issuer user.Issuer
		want   string
	}{
		{
			name:   "testing cab event",
			arg:    "passenger-event-152384980615c2bd16143cff29038b67",
			issuer: user.Passenger,
			want:   topics.CabEvent,
		},
		{
			name:   "testing cab event",
			arg:    "driver-event-152384980615c2bd16143cff29038b67",
			issuer: user.Driver,
			want:   topics.CabEvent,
		},
		{
			name:   "testing invalid event",
			arg:    "-event-123456789abcdefgABCDEFG",
			issuer: user.None,
			want:   "",
		},
		{
			name:   "testing driver location",
			arg:    "snapp/driver/DXKgaNQa7N5Y7bo/location",
			issuer: user.Driver,
			want:   topics.DriverLocation,
		},
		{
			name:   "testing passenger location",
			arg:    "snapp/passenger/DXKgaNQa7N5Y7bo/location",
			issuer: user.Passenger,
			want:   topics.PassengerLocation,
		},
		{
			name:   "testing invalid location",
			arg:    "snapp/thirdparty/DXKgaNQa7N5Y7bo/location",
			issuer: user.None,
			want:   "",
		},
		{
			name:   "testing superapp event",
			arg:    "snapp/passenger/DXKgaNQa7N5Y7bo/superapp",
			issuer: user.Passenger,
			want:   topics.SuperappEvent,
		},
		{
			name:   "testing shared passenger location",
			arg:    "snapp/passenger/DXKgaNQa7N5Y7bo/driver-location",
			issuer: user.Passenger,
			want:   topics.SharedLocation,
		},
		{
			name:   "testing shared driver location",
			arg:    "snapp/driver/DXKgaNQa7N5Y7bo/passenger-location",
			issuer: user.Driver,
			want:   topics.SharedLocation,
		},
		{
			name:   "testing passenger chat",
			arg:    "snapp/passenger/DXKgaNQa7N5Y7bo/chat",
			issuer: user.Passenger,
			want:   topics.Chat,
		},
		{
			name:   "testing driver chat",
			arg:    "snapp/driver/DXKgaNQa7N5Y7bo/chat",
			issuer: user.Driver,
			want:   topics.Chat,
		},
		{
			name:   "testing passenger general call entry",
			arg:    "shared/snapp/passenger/DXKgaNQa7N5Y7bo/call/send",
			issuer: user.Passenger,
			want:   topics.GeneralCallEntry,
		},
		{
			name:   "testing driver general call entry",
			arg:    "shared/snapp/driver/DXKgaNQa7N5Y7bo/call/send",
			issuer: user.Driver,
			want:   topics.GeneralCallEntry,
		},
		{
			name:   "testing passenger node call entry",
			arg:    "snapp/passenger/DXKgaNQa7N5Y7bo/call/heliograph-0/send",
			issuer: user.Passenger,
			want:   topics.NodeCallEntry,
		},
		{
			name:   "testing driver node call entry",
			arg:    "snapp/driver/DXKgaNQa7N5Y7bo/call/heliograph-1/send",
			issuer: user.Driver,
			want:   topics.NodeCallEntry,
		},
		{
			name:   "testing passenger call",
			arg:    "snapp/passenger/DXKgaNQa7N5Y7bo/call/receive",
			issuer: user.Passenger,
			want:   topics.CallOutgoing,
		},
		{
			name:   "testing driver call",
			arg:    "snapp/driver/DXKgaNQa7N5Y7bo/call/receive",
			issuer: user.Driver,
			want:   topics.CallOutgoing,
		},
		{
			name:   "testing box event",
			arg:    "bucks",
			issuer: user.Driver,
			want:   topics.BoxEvent,
		},
	}

	cfg := config.New()

	hid := &snappids.HashIDSManager{
		Salts: map[snappids.Audience]string{
			snappids.PassengerAudience:  "secret",
			snappids.DriverAudience:     "secret",
			snappids.ThirdPartyAudience: "secret",
		},
		Lengths: map[snappids.Audience]int{
			snappids.PassengerAudience:  15,
			snappids.DriverAudience:     15,
			snappids.ThirdPartyAudience: 15,
		},
	}
	auth := authenticator.Authenticator{
		Company:      "snapp",
		TopicManager: topics.NewTopicManager(cfg.Topics, hid, "snapp"),
	}

	sub := "DXKgaNQa7N5Y7bo"

	for i, tt := range tests {
		t1.Run(fmt.Sprintf("#%d %s", i, tt.name), func(t1 *testing.T) {
			t := tt.arg
			audience, audienceStr := topics.IssuerToAudience(tt.issuer)
			topicTemplate := auth.TopicManager.ValidateTopic(t, audienceStr, audience, sub)
			if topicTemplate != nil {
				if len(tt.want) == 0 {
					t1.Errorf("topic %s is invalid, must throw error.", tt.arg)
				} else if topicTemplate.Type != tt.want {
					t1.Errorf("GetType() = %v, want %v", topicTemplate.Type, tt.want)
				}
			} else {
				if len(tt.want) != 0 {
					t1.Errorf("failed to find topicTemplate for %s", tt.arg)
				}
			}
		})
	}
}