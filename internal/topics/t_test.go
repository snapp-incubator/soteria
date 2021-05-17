package topics

import "testing"

func TestTopic_GetType(t1 *testing.T) {
	tests := []struct {
		name string
		arg  Topic
		want Type
	}{
		{
			name: "#1 testing cab event",
			arg:  Topic("passenger-event-123456789abcdefgABCDEFG"),
			want: CabEvent,
		},
		{
			name: "#2 testing cab event",
			arg:  "driver-event-123456789abcdefgABCDEFG",
			want: CabEvent,
		},
		{
			name: "#3 testing invalid event",
			arg:  "-event-123456789abcdefgABCDEFG",
			want: "",
		},
		{
			name: "#4 testing driver location",
			arg:  "snapp/driver/sfhsdkifs475sfhs/location",
			want: DriverLocation,
		},
		{
			name: "#5 testing passenger location",
			arg:  "snapp/passenger/sfhsdkifs475sfhs/location",
			want: PassengerLocation,
		},
		{
			name: "#6 testing invalid location",
			arg:  "snapp/thirdparty/sfhsdkifs475sfhs/location",
			want: "",
		},
		{
			name: "#7 testing superapp event",
			arg:  "snapp/passenger/fhdyfuiksdf5456456adljada/superapp",
			want: SuperappEvent,
		},
		{
			name: "#8 testing superapp event",
			arg:  "snapp/driver/+/location",
			want: DriverLocation,
		},
		{
			name: "#9 testing daghigh sys",
			arg:  "$SYS/brokers/+/clients/+/disconnected",
			want: DaghighSys,
		},
		{
			name: "#10 testing daghigh sys",
			arg:  "$SYS/brokers/+/clients/+/connected",
			want: DaghighSys,
		},
		{
			name: "#11 testing daghigh sys",
			arg:  "$share/hello/$SYS/brokers/+/clients/+/connected",
			want: DaghighSys,
		},
		{
			name: "#12 testing gossiper passenger location",
			arg:  "snapp/passenger/py9kdjLYB35RP4q/driver-location",
			want: GossiperLocation,
		},
		{
			name: "#13 testing gossiper driver location",
			arg:  "snapp/driver/py9kdjLYB35RP4q/passenger-location",
			want: GossiperLocation,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := tt.arg
			if got := t.GetType(); got != tt.want {
				t1.Errorf("GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}
