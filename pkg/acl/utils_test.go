package acl_test

import (
	"testing"

	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/acl"
)

// nolint: funlen
func TestValidateEndpoint(t *testing.T) {
	t.Parallel()

	type args struct {
		endpoint              string
		authorizedEndpoints   []string
		unauthorizedEndpoints []string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "#1",
			args: args{
				endpoint:              "",
				authorizedEndpoints:   []string{},
				unauthorizedEndpoints: []string{},
			},
			want: true,
		},
		{
			name: "#2",
			args: args{
				endpoint:              "/push",
				authorizedEndpoints:   []string{"/push"},
				unauthorizedEndpoints: []string{"/push"},
			},
			want: false,
		},
		{
			name: "#3",
			args: args{
				endpoint:              "/push",
				authorizedEndpoints:   []string{"/push"},
				unauthorizedEndpoints: []string{"/events"},
			},
			want: true,
		},
		{
			name: "#4",
			args: args{
				endpoint:              "/events",
				authorizedEndpoints:   []string{"/push"},
				unauthorizedEndpoints: []string{"/events"},
			},
			want: false,
		},
		{
			name: "#5",
			args: args{
				endpoint:              "/events",
				authorizedEndpoints:   []string{"/push"},
				unauthorizedEndpoints: []string{"/event"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := acl.ValidateEndpoint(tt.args.endpoint, tt.args.authorizedEndpoints,
				tt.args.unauthorizedEndpoints); got != tt.want {
				t.Errorf("ValidateEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

// nolint: funlen
func TestValidateIP(t *testing.T) {
	t.Parallel()

	type args struct {
		IP         string
		validIPs   []string
		invalidIPs []string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "#1",
			args: args{
				IP:         "192.168.24.56",
				validIPs:   []string{},
				invalidIPs: []string{},
			},
			want: false,
		},
		{
			name: "#2",
			args: args{
				IP:         "192.168.24.56",
				validIPs:   []string{"192.168.24.56"},
				invalidIPs: []string{"192.168.24.56"},
			},
			want: false,
		},
		{
			name: "#3",
			args: args{
				IP:         "192.168.24.56",
				validIPs:   []string{"192.168.24.0/8"},
				invalidIPs: []string{"192.168.24.56"},
			},
			want: false,
		},
		{
			name: "#4",
			args: args{
				IP:         "192.168.24.56",
				validIPs:   []string{"192.168.24.0/8"},
				invalidIPs: []string{"192.168.24.8"},
			},
			want: true,
		},
		{
			name: "#5",
			args: args{
				IP:         "192.168.24.56",
				validIPs:   []string{"192.168.24.56"},
				invalidIPs: []string{"192.168.24.8/8"},
			},
			want: false,
		},
		{
			name: "#6",
			args: args{
				IP:         "192.168.24.56",
				validIPs:   []string{"192.168.24.1", "192.168.0.0/16"},
				invalidIPs: []string{"192.168.24.89"},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := acl.ValidateIP(tt.args.IP, tt.args.validIPs, tt.args.invalidIPs); got != tt.want {
				t.Errorf("ValidateIP() = %v, want %v", got, tt.want)
			}
		})
	}
}
