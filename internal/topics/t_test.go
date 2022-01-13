package topics

import (
	"fmt"
	"testing"
)

func TestTopic_GetType(t1 *testing.T) {
	tests := []struct {
		name string
		arg  Topic
		want Type
	}{
		{
			name: "testing cab event",
			arg:  Topic("passenger-event-123456789abcdefgABCDEFG"),
			want: CabEvent,
		},
		{
			name: "testing cab event",
			arg:  "driver-event-123456789abcdefgABCDEFG",
			want: CabEvent,
		},
		{
			name: "testing invalid event",
			arg:  "-event-123456789abcdefgABCDEFG",
			want: "",
		},
		{
			name: "testing driver location",
			arg:  "snapp/driver/sfhsdkifs475sfhs/location",
			want: DriverLocation,
		},
		{
			name: "testing passenger location",
			arg:  "snapp/passenger/sfhsdkifs475sfhs/location",
			want: PassengerLocation,
		},
		{
			name: "testing invalid location",
			arg:  "snapp/thirdparty/sfhsdkifs475sfhs/location",
			want: "",
		},
		{
			name: "testing superapp event",
			arg:  "snapp/passenger/fhdyfuiksdf5456456adljada/superapp",
			want: SuperappEvent,
		},
		{
			name: "testing shared passenger location",
			arg:  "snapp/passenger/py9kdjLYB35RP4q/driver-location",
			want: SharedLocation,
		},
		{
			name: "testing shared driver location",
			arg:  "snapp/driver/py9kdjLYB35RP4q/passenger-location",
			want: SharedLocation,
		},
		{
			name: "testing passenger chat",
			arg:  "snapp/passenger/py9kdjLYB35RP4q/chat",
			want: Chat,
		},
		{
			name: "testing driver chat",
			arg:  "snapp/driver/py9kdjLYB35RP4q/chat",
			want: Chat,
		},
		{
			name: "testing passenger general call entry",
			arg:  "shared/snapp/passenger/py9kdjLYB35RP4q/call/send",
			want: GeneralCallEntry,
		},
		{
			name: "testing driver general call entry",
			arg:  "shared/snapp/driver/py9kdjLYB35RP4q/call/send",
			want: GeneralCallEntry,
		},
		{
			name: "testing passenger node call entry",
			arg:  "snapp/passenger/py9kdjLYB35RP4q/call/heliograph-0/send",
			want: NodeCallEntry,
		},
		{
			name: "testing driver node call entry",
			arg:  "snapp/driver/py9kdjLYB35RP4q/call/heliograph-1/send",
			want: NodeCallEntry,
		},
		{
			name: "testing passenger call",
			arg:  "snapp/passenger/py9kdjLYB35RP4q/call/receive",
			want: CallOutgoing,
		},
		{
			name: "testing driver call",
			arg:  "snapp/driver/py9kdjLYB35RP4q/call/receive",
			want: CallOutgoing,
		},
	}
	for i, tt := range tests {
		t1.Run(fmt.Sprintf("#%d %s", i, tt.name), func(t1 *testing.T) {
			t := tt.arg
			if got := t.GetType(); got != tt.want {
				t1.Errorf("GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}
