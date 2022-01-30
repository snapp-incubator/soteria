package topics_test

import (
	"fmt"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/authenticator"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/config"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/topics"
	"testing"
)

func TestTopic_GetType(t1 *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "testing cab event",
			arg:  "passenger-event-123456789abcdefgABCDEFG",
			want: topics.CabEvent,
		},
		{
			name: "testing cab event",
			arg:  "driver-event-123456789abcdefgABCDEFG",
			want: topics.CabEvent,
		},
		{
			name: "testing invalid event",
			arg:  "-event-123456789abcdefgABCDEFG",
			want: "",
		},
		{
			name: "testing driver location",
			arg:  "snapp/driver/sfhsdkifs475sfhs/location",
			want: topics.DriverLocation,
		},
		{
			name: "testing passenger location",
			arg:  "snapp/passenger/sfhsdkifs475sfhs/location",
			want: topics.PassengerLocation,
		},
		{
			name: "testing invalid location",
			arg:  "snapp/thirdparty/sfhsdkifs475sfhs/location",
			want: "",
		},
		{
			name: "testing superapp event",
			arg:  "snapp/passenger/fhdyfuiksdf5456456adljada/superapp",
			want: topics.SuperappEvent,
		},
		{
			name: "testing shared passenger location",
			arg:  "snapp/passenger/py9kdjLYB35RP4q/driver-location",
			want: topics.SharedLocation,
		},
		{
			name: "testing shared driver location",
			arg:  "snapp/driver/py9kdjLYB35RP4q/passenger-location",
			want: topics.SharedLocation,
		},
		{
			name: "testing passenger chat",
			arg:  "snapp/passenger/py9kdjLYB35RP4q/chat",
			want: topics.Chat,
		},
		{
			name: "testing driver chat",
			arg:  "snapp/driver/py9kdjLYB35RP4q/chat",
			want: topics.Chat,
		},
		{
			name: "testing passenger general call entry",
			arg:  "shared/snapp/passenger/py9kdjLYB35RP4q/call/send",
			want: topics.GeneralCallEntry,
		},
		{
			name: "testing driver general call entry",
			arg:  "shared/snapp/driver/py9kdjLYB35RP4q/call/send",
			want: topics.GeneralCallEntry,
		},
		{
			name: "testing passenger node call entry",
			arg:  "snapp/passenger/py9kdjLYB35RP4q/call/heliograph-0/send",
			want: topics.NodeCallEntry,
		},
		{
			name: "testing driver node call entry",
			arg:  "snapp/driver/py9kdjLYB35RP4q/call/heliograph-1/send",
			want: topics.NodeCallEntry,
		},
		{
			name: "testing passenger call",
			arg:  "snapp/passenger/py9kdjLYB35RP4q/call/receive",
			want: topics.CallOutgoing,
		},
		{
			name: "testing driver call",
			arg:  "snapp/driver/py9kdjLYB35RP4q/call/receive",
			want: topics.CallOutgoing,
		},
		{
			name: "testing box event",
			arg:  "bucks",
			want: topics.BoxEvent,
		},
	}

	cfg := config.New()
	auth := authenticator.Authenticator{
		TopicManager: topics.NewTopicManager(cfg.Topics, nil, "snapp"),
	}

	for i, tt := range tests {
		t1.Run(fmt.Sprintf("#%d %s", i, tt.name), func(t1 *testing.T) {
			t := tt.arg
			if got := auth.TopicManager.GetTopicType(t, "snapp"); got != tt.want {
				t1.Errorf("GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}
