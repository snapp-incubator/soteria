package api_test

import (
	"testing"

	"github.com/snapp-incubator/soteria/internal/api"
)

// nolint: funlen
func TestExtractVendorToken(t *testing.T) {
	t.Parallel()

	type fields struct {
		Token    string
		Username string
		Password string
	}

	tests := []struct {
		name   string
		fields fields
		vendor string
		token  string
	}{
		{
			name: "token field as token",
			fields: fields{
				Token:    "vendor:token",
				Username: "vendor:username",
				Password: "password",
			},
			vendor: "vendor",
			token:  "token",
		},
		{
			name: "username as token without vendor",
			fields: fields{
				Token:    "",
				Username: "username",
				Password: "password",
			},
			vendor: "",
			token:  "username",
		},
		{
			name: "username as token with vendor",
			fields: fields{
				Token:    "",
				Username: "vendor:username",
				Password: "",
			},
			vendor: "vendor",
			token:  "username",
		},
		{
			name: "password as token",
			fields: fields{
				Token:    "",
				Username: "",
				Password: "vendor:password",
			},
			vendor: "vendor",
			token:  "password",
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			vendor, token := api.ExtractVendorToken(tt.fields.Token, tt.fields.Username, tt.fields.Password)
			if vendor != tt.vendor {
				t.Errorf("ExtractVendorToken() vendor = %v, vendor %v", vendor, tt.vendor)
			}
			if token != tt.token {
				t.Errorf("ExtractVendorToken() token = %v, vendor %v", token, tt.token)
			}
		})
	}
}
