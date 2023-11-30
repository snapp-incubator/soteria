package authenticator_test

import (
	"testing"

	"github.com/snapp-incubator/soteria/internal/authenticator"
	"github.com/snapp-incubator/soteria/internal/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestBuilderWithoutAuthenticator(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	b := authenticator.Builder{
		Vendors: []config.Vendor{},
		Logger:  zap.NewNop(),
		ValidatorConfig: config.Validator{
			URL:     "",
			Timeout: 0,
		},
	}

	_, err := b.Authenticators()
	require.ErrorIs(err, authenticator.ErrNoAuthenticator)
}

func TestBuilderInternalAuthenticator(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	b := authenticator.Builder{
		Vendors: []config.Vendor{
			{
				Company: "internal",
				Jwt: config.Jwt{
					IssName:       "iss",
					SubName:       "sub",
					SigningMethod: "HS512",
				},
				IsInternal:         true,
				UseValidator:       false,
				AllowedAccessTypes: nil,
				Topics:             nil,
				HashIDMap:          nil,
				IssEntityMap:       nil,
				IssPeerMap:         nil,
				Keys: map[string]string{
					"system": "c2VjcmV0",
				},
			},
		},
		Logger: zap.NewNop(),
		ValidatorConfig: config.Validator{
			URL:     "",
			Timeout: 0,
		},
	}

	vendors, err := b.Authenticators()
	require.NoError(err)
	require.Len(vendors, 1)
	require.Contains(vendors, "internal")
}

func TestBuilderInternalAuthenticatorWithInvalidKey(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	b := authenticator.Builder{
		Vendors: []config.Vendor{
			{
				Company: "internal",
				Jwt: config.Jwt{
					IssName:       "iss",
					SubName:       "sub",
					SigningMethod: "HS512",
				},
				IsInternal:         true,
				UseValidator:       false,
				AllowedAccessTypes: nil,
				Topics:             nil,
				HashIDMap:          nil,
				IssEntityMap:       nil,
				IssPeerMap:         nil,
				Keys: map[string]string{
					"not-system": "c2VjcmV0",
				},
			},
		},
		Logger: zap.NewNop(),
		ValidatorConfig: config.Validator{
			URL:     "",
			Timeout: 0,
		},
	}

	_, err := b.Authenticators()
	require.ErrorIs(err, authenticator.ErrAdminAuthenticatorSystemKey)
}
