package authenticator_test

import (
	"testing"

	"github.com/snapp-incubator/soteria/internal/authenticator"
	"github.com/snapp-incubator/soteria/internal/config"
	"github.com/snapp-incubator/soteria/internal/topics"
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
				Jwt: config.JWT{
					IssName:       "iss",
					SubName:       "sub",
					SigningMethod: "HS512",
				},
				Type:               "internal",
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
				Jwt: config.JWT{
					IssName:       "iss",
					SubName:       "sub",
					SigningMethod: "HS512",
				},
				Type:               "internal",
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

func TestBuilderManualAuthenticatorWithoutKey(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	b := authenticator.Builder{
		Vendors: []config.Vendor{
			{
				Company: "snapp",
				Jwt: config.JWT{
					IssName:       "iss",
					SubName:       "sub",
					SigningMethod: "HS512",
				},
				Type:               "manual",
				AllowedAccessTypes: []string{"pub", "sub"},
				Topics:             nil,
				HashIDMap: map[string]topics.HashData{
					"0": {
						Alphabet: "",
						Length:   15,
						Salt:     "secret",
					},
					"1": {
						Alphabet: "",
						Length:   15,
						Salt:     "secret",
					},
				},
				IssEntityMap: map[string]string{
					"0":       "driver",
					"1":       "passenger",
					"default": "",
				},
				IssPeerMap: map[string]string{
					"0":       "passenger",
					"1":       "driver",
					"default": "",
				},
				Keys: nil,
			},
		},
		Logger: zap.NewNop(),
		ValidatorConfig: config.Validator{
			URL:     "",
			Timeout: 0,
		},
	}

	_, err := b.Authenticators()
	require.ErrorIs(err, authenticator.ErrNoKeys)
}

// nolint: funlen
func TestBuilderManualAuthenticator(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	b := authenticator.Builder{
		Vendors: []config.Vendor{
			{
				Company: "snapp",
				Jwt: config.JWT{
					IssName:       "iss",
					SubName:       "sub",
					SigningMethod: "RSA512",
				},
				Type:               "manual",
				AllowedAccessTypes: []string{"pub", "sub"},
				Topics:             nil,
				HashIDMap: map[string]topics.HashData{
					"0": {
						Alphabet: "",
						Length:   15,
						Salt:     "secret",
					},
					"1": {
						Alphabet: "",
						Length:   15,
						Salt:     "secret",
					},
				},
				IssEntityMap: map[string]string{
					"0":       "driver",
					"1":       "passenger",
					"default": "",
				},
				IssPeerMap: map[string]string{
					"0":       "passenger",
					"1":       "driver",
					"default": "",
				},
				Keys: map[string]string{
					"0": `-----BEGIN PUBLIC KEY-----
			MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyG4XpV9TpDfgWJF9TiIv
			va4hNhDuqYMJO6iXLzr3y8oCvoB7zUK0EjtbLH+A3gr1kUvyZKDWT4qHTvU2Sshm
			X+ttWGK34EhCvF3Lb18yxmVDSSK8JIcTaJjMqmyubxzamQnNoWazJ7ea9BIo2YGL
			C9rgPbi1hihhdb07xPGUkJRqbWkI98xjDhKdMqiwW1hIRXm/apo++FjptvqvF84s
			ynC5gWGFHiGNICRsLJBczLEAf2Atbafigq6/tovzMabnp2yRtr1ReEgioH1RO4gX
			J7F4N5f6y/VWd8+sDOSxtS/HcnP/7g8/A54G2IbXxr+EiwOO/1F+pyMPKq7sGDSU
			DwIDAQAB
-----END PUBLIC KEY-----`,

					"1": `-----BEGIN PUBLIC KEY-----
			MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA5SeRfOdTyvQZ7N9ahFHl
        +J05r7e9fgOQ2cpOtnnsIjAjCt1dF7/NkqVifEaxABRBGG9iXIw//G4hi0TqoKqK
        aoSHMGf6q9pSRLGyB8FatxZf2RBTgrXYqVvpasbnB1ZNv858yTpRjV9NzJXYHLp8
        8Hbd/yYTR6Q7ajs11/SMLGO7KBELsI1pBz7UW/fngJ2pRmd+RkG+EcGrOIZ27TkI
        Xjtog6bgfmtV9FWxSVdKACOY0OmW+g7jIMik2eZTYG3kgCmW2odu3zRoUa7l9VwN
        YMuhTePaIWwOifzRQt8HDsAOpzqJuLCoYX7HmBfpGAnwu4BuTZgXVwpvPNb+KlgS
        pQIDAQAB
-----END PUBLIC KEY-----`,
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
	require.Contains(vendors, "snapp")
}
