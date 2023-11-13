package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/snapp-incubator/soteria/internal/api"
	"github.com/snapp-incubator/soteria/internal/authenticator"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap"
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

func TestAuthv2(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	app := fiber.New()

	a := api.API{
		Authenticators: map[string]authenticator.Authenticator{},
		DefaultVendor:  "snapp",
		Tracer:         noop.NewTracerProvider().Tracer(""),
		Logger:         zap.NewExample(),
	}

	app.Post("/v2/auth", a.Authv2)

	t.Run("bad request because it doesn't have json heaer", func(t *testing.T) {
		t.Parallel()

		body, err := json.Marshal(api.AuthRequest{
			Token:    "",
			Username: "not-found:token",
			Password: "",
		})
		require.NoError(err)

		req := httptest.NewRequest(http.MethodPost, "/v2/auth", bytes.NewReader(body))

		resp, err := app.Test(req)
		require.NoError(err)

		defer resp.Body.Close()

		require.Equal(http.StatusBadRequest, resp.StatusCode)
	})
}
