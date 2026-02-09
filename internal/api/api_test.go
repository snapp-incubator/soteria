package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/snapp-incubator/soteria/internal/api"
	"github.com/snapp-incubator/soteria/internal/authenticator"
	"github.com/snapp-incubator/soteria/internal/clientid"
	"github.com/snapp-incubator/soteria/internal/config"
	"github.com/snapp-incubator/soteria/internal/metric"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap"
)

func getSampleToken(key string) (string, error) {
	exp := time.Now().Add(time.Hour * 24 * 365 * 10)
	sub := "DXKgaNQa7N5Y7bo"

	// nolint: exhaustruct
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(exp),
		Issuer:    "Colony",
		Subject:   sub,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("cannot generate a signed string %w", err)
	}

	return tokenString, nil
}

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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			vendor, token := api.ExtractVendorToken(tt.fields.Token, tt.fields.Username, tt.fields.Password)
			if vendor != tt.vendor || token != tt.token {
				t.Errorf("ExtractVendorToken() vendor = %v, vendor %v", vendor, tt.vendor)
			}
		})
	}
}

type APITestSuite struct {
	suite.Suite

	app *fiber.App
	key string
}

func (suite *APITestSuite) SetupSuite() {
	suite.key = "secret"

	app := fiber.New()

	a := api.API{
		Authenticators: map[string]authenticator.Authenticator{
			"snapp-admin": authenticator.AdminAuthenticator{
				Key:     []byte(suite.key),
				Company: "snapp-admin",
				JwtConfig: config.JWT{
					IssName:       "iss",
					SubName:       "sub",
					SigningMethod: "HS512",
				},
				Parser: jwt.NewParser(),
			},
		},
		DefaultVendor: "snapp-admin",
		Tracer:        noop.NewTracerProvider().Tracer(""),
		Logger:        zap.NewExample(),
		Metrics:       metric.NewAPIMetrics(),
		Parser: clientid.NewParser(clientid.Config{
			Patterns: map[string]string{},
		}),
	}

	app.Post("/v2/auth", a.Authv2)

	suite.app = app
}

func (suite *APITestSuite) TestToken() {
	require := suite.Require()

	token, err := getSampleToken(suite.key)
	require.NoError(err)

	cases := []struct {
		username    string
		action      string
		isSuperuser bool
		sendHeader  bool
	}{
		{username: "snapp-admin:" + token, action: "allow", isSuperuser: true, sendHeader: true},
		{username: token, action: "allow", isSuperuser: true, sendHeader: true},
		{username: "not-found:token", action: "deny", isSuperuser: false, sendHeader: true},
		{username: "token", action: "deny", isSuperuser: false, sendHeader: true},
		{username: "snapp-admin:token", action: "deny", isSuperuser: false, sendHeader: true},
		{username: "snapp-admin:" + token, action: "deny", isSuperuser: false, sendHeader: false},
	}

	for _, c := range cases {
		suite.Run("username "+c.username, func() {
			body, err := json.Marshal(api.AuthRequest{
				Token:    "",
				Username: c.username,
				Password: "",
				ClientID: "",
			})
			require.NoError(err)

			req := httptest.NewRequest(http.MethodPost, "/v2/auth", bytes.NewReader(body))
			if c.sendHeader {
				req.Header.Add("Content-Type", "application/json")
			}

			resp, err := suite.app.Test(req)
			require.NoError(err)

			defer func() {
				_ = resp.Body.Close()
			}()

			require.Equal(http.StatusOK, resp.StatusCode)

			data, err := io.ReadAll(resp.Body)
			require.NoError(err)

			var authResp api.AuthResponse

			require.NoError(json.Unmarshal(data, &authResp))

			require.Equal(c.action, authResp.Result)
			require.Equal(c.isSuperuser, authResp.IsSuperuser)
		})
	}
}

func TestAPITestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(APITestSuite))
}

func TestMetricsEndpoint(t *testing.T) {
	t.Parallel()

	reg := prometheus.NewRegistry()
	app := fiber.New()

	prom := api.NewPrometheusMiddleware(reg)
	prom.RegisterAt(app, "/metrics")
	app.Use(prom.Handler)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	defer func() {
		_ = resp.Body.Close()
	}()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Contains(t, string(body), "platform_soteria_http_requests_in_progress")
}
