package authenticator_test

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/snapp-incubator/soteria/internal/authenticator"
	"github.com/snapp-incubator/soteria/pkg/acl"
	"github.com/stretchr/testify/suite"
)

const (
	// nolint: gosec, lll
	invalidToken = "ey1JhbGciOiJSUzI1NiIsInR5cCI56kpXVCJ9.eyJzdWIiOiJCRzdScDFkcnpWRE5RcjYiLCJuYW1lIjoiSm9obiBEb2UiLCJhZG1pbiI6dHJ1ZSwiaXNzIjowLCJpYXQiOjE1MTYyMzkwMjJ9.1cYXFEhcewOYFjGJYhB8dsaFO9uKEXwlM8954rkt4Tsu0lWMITbRf_hHh1l9QD4MFqD-0LwRPUYaiaemy0OClMu00G2sujLCWaquYDEP37iIt8RoOQAh8Jb5vT8LX5C3PEKvbW_i98u8HHJoFUR9CXJmzrKi48sAcOYvXVYamN0S9KoY38H-Ze37Mdu3o6B58i73krk7QHecsc2_PkCJisvUVAzb0tiInIalBc8-zI3QZSxwNLr_hjlBg1sUxTUvH5SCcRR7hxI8TxJzkOHqAHWDRO84NC_DSAoO2p04vrHpqglN9XPJ8RC2YWpfefvD2ttH554RJWu_0RlR2kAYvQ"

	validPassengerCabEventTopic   = "passenger-event-152384980615c2bd16143cff29038b67"
	invalidPassengerCabEventTopic = "passenger-event-152384980615c2bd16156cff29038b67"

	validDriverCabEventTopic   = "driver-event-152384980615c2bd16143cff29038b67"
	invalidDriverCabEventTopic = "driver-event-152384980615c2bd16156cff29038b67"

	validDriverLocationTopic   = "snapp/driver/DXKgaNQa7N5Y7bo/location"
	invalidDriverLocationTopic = "snapp/driver/DXKgaNQa9Q5Y7bo/location"

	validPassengerSuperappEventTopic   = "snapp/passenger/DXKgaNQa7N5Y7bo/superapp"
	invalidPassengerSuperappEventTopic = "snapp/passenger/DXKgaNQa9Q5Y7bo/superapp"

	validDriverSuperappEventTopic   = "snapp/driver/DXKgaNQa7N5Y7bo/superapp"
	invalidDriverSuperappEventTopic = "snapp/driver/DXKgaNQa9Q5Y7bo/superapp"

	validDriverSharedTopic      = "snapp/driver/DXKgaNQa7N5Y7bo/passenger-location"
	validPassengerSharedTopic   = "snapp/passenger/DXKgaNQa7N5Y7bo/driver-location"
	invalidDriverSharedTopic    = "snapp/driver/0596923be632d673560af9adadd2f78a/passenger-location"
	invalidPassengerSharedTopic = "snapp/passenger/0596923be632d673560af9adadd2f78a/driver-location"

	validDriverChatTopic      = "snapp/driver/DXKgaNQa7N5Y7bo/chat"
	validPassengerChatTopic   = "snapp/passenger/DXKgaNQa7N5Y7bo/chat"
	invalidDriverChatTopic    = "snapp/driver/0596923be632d673560af9adadd2f78a/chat"
	invalidPassengerChatTopic = "snapp/passenger/0596923be632d673560af9adadd2f78a/chat"

	validDriverCallEntryTopic         = "shared/snapp/driver/DXKgaNQa7N5Y7bo/call/send"
	validPassengerCallEntryTopic      = "shared/snapp/passenger/DXKgaNQa7N5Y7bo/call/send"
	validDriverNodeCallEntryTopic     = "snapp/driver/DXKgaNQa7N5Y7bo/call/heliograph-0/send"
	validPassengerNodeCallEntryTopic  = "snapp/passenger/DXKgaNQa7N5Y7bo/call/heliograph-0/send"
	invalidDriverCallEntryTopic       = "snapp/driver/0596923be632d673560af9adadd2f78a/call/send"
	invalidPassengerCallEntryTopic    = "snapp/passenger/0596923be632d673560af9adadd2f78a/call/send"
	validDriverCallOutgoingTopic      = "snapp/driver/DXKgaNQa7N5Y7bo/call/receive"
	validPassengerCallOutgoingTopic   = "snapp/passenger/DXKgaNQa7N5Y7bo/call/receive"
	invalidDriverCallOutgoingTopic    = "snapp/driver/0596923be632d673560af9adadd2f78a/call/receive"
	invalidPassengerCallOutgoingTopic = "snapp/passenger/0596923be632d673560af9adadd2f78a/call/receive"
)

var (
	ErrPrivateKeyNotFound = errors.New("invalid user, private key not found")
	ErrPublicKeyNotFound  = errors.New("invalid user, public key not found")
)

type AuthenticatorTestSuite struct {
	suite.Suite

	Tokens struct {
		Passenger string
		Driver    string
	}

	PublicKeys struct {
		Passenger *rsa.PublicKey
		Driver    *rsa.PublicKey
	}

	Authenticator authenticator.Authenticator
}

func (suite *AuthenticatorTestSuite) SetupSuite() {
	require := suite.Require()

	pkey0, err := getPublicKey("0")
	require.NoError(err)

	suite.PublicKeys.Driver = pkey0

	pkey1, err := getPublicKey("1")
	require.NoError(err)

	suite.PublicKeys.Passenger = pkey1

	key0, err := getPrivateKey("0")
	require.NoError(err)

	suite.PublicKeys.Driver = pkey0

	key1, err := getPrivateKey("1")
	require.NoError(err)

	driverToken, err := getSampleToken("0", key0)
	require.NoError(err)

	suite.Tokens.Driver = driverToken

	passengerToken, err := getSampleToken("1", key1)
	require.NoError(err)

	suite.Tokens.Passenger = passengerToken
}

func (suite *AuthenticatorTestSuite) TestAuth() {
	require := suite.Require()

	suite.Run("testing driver token auth", func() {
		require.NoError(suite.Authenticator.Auth(suite.Tokens.Driver))
	})

	suite.Run("testing passenger token auth", func() {
		require.NoError(suite.Authenticator.Auth(suite.Tokens.Passenger))
	})

	suite.Run("testing invalid token auth", func() {
		require.Error(suite.Authenticator.Auth(invalidToken))
	})
}

// nolint: dupl
func (suite *AuthenticatorTestSuite) TestACL_Basics() {
	require := suite.Require()

	suite.Run("testing acl with invalid access type", func() {
		ok, err := suite.Authenticator.ACL("invalid-access", suite.Tokens.Passenger, "test")
		require.Error(err)
		require.False(ok)
		require.ErrorIs(err, authenticator.ErrInvalidAccessType)
	})

	suite.Run("testing acl with invalid token", func() {
		ok, err := suite.Authenticator.ACL(acl.Pub, invalidToken, validDriverCabEventTopic)
		require.False(ok)
		require.Error(err)
		require.ErrorIs(err, jwt.ErrTokenMalformed)
	})

	suite.Run("testing acl with valid inputs", func() {
		ok, err := suite.Authenticator.ACL(acl.Sub, suite.Tokens.Passenger, validPassengerCabEventTopic)
		require.NoError(err)
		require.True(ok)
	})

	suite.Run("testing acl with invalid topic", func() {
		ok, err := suite.Authenticator.ACL(acl.Sub, suite.Tokens.Passenger, invalidPassengerCabEventTopic)
		require.Error(err)
		require.False(ok)
	})

	suite.Run("testing acl with invalid access type", func() {
		ok, err := suite.Authenticator.ACL(acl.Pub, suite.Tokens.Passenger, validPassengerCabEventTopic)
		require.Error(err)
		require.False(ok)
	})
}

// nolint: funlen
func (suite *AuthenticatorTestSuite) TestACL_Passenger() {
	require := suite.Require()
	token := suite.Tokens.Passenger

	suite.Run("testing passenger subscribe on valid superapp event topic", func() {
		ok, err := suite.Authenticator.ACL(acl.Sub, token, validPassengerSuperappEventTopic)
		require.NoError(err)
		require.True(ok)
	})

	suite.Run("testing passenger subscribe on invalid superapp event topic", func() {
		ok, err := suite.Authenticator.ACL(acl.Sub, token, invalidPassengerSuperappEventTopic)
		require.Error(err)
		require.False(ok)
	})

	suite.Run("testing passenger subscribe on valid shared location topic", func() {
		ok, err := suite.Authenticator.ACL(acl.Sub, token, validPassengerSharedTopic)
		require.NoError(err)
		require.True(ok)
	})

	suite.Run("testing passenger subscribe on invalid shared location topic", func() {
		ok, err := suite.Authenticator.ACL(acl.Sub, token, invalidPassengerSharedTopic)
		require.Error(err)
		require.False(ok)
	})

	suite.Run("testing passenger subscribe on valid chat topic", func() {
		ok, err := suite.Authenticator.ACL(acl.Sub, token, validPassengerChatTopic)
		require.NoError(err)
		require.True(ok)
	})

	suite.Run("testing passenger subscribe on invalid chat topic", func() {
		ok, err := suite.Authenticator.ACL(acl.Sub, token, invalidPassengerChatTopic)
		require.Error(err)
		require.False(ok)
	})

	suite.Run("testing passenger subscribe on valid entry call topic", func() {
		ok, err := suite.Authenticator.ACL(acl.Pub, token, validPassengerCallEntryTopic)
		require.NoError(err)
		require.True(ok)
	})

	suite.Run("testing passenger subscribe on invalid call entry topic", func() {
		ok, err := suite.Authenticator.ACL(acl.Pub, token, invalidPassengerCallEntryTopic)
		require.Error(err)
		require.False(ok)
	})

	suite.Run("testing passenger subscribe on valid outgoing call topic", func() {
		ok, err := suite.Authenticator.ACL(acl.Sub, token, validPassengerCallOutgoingTopic)
		require.NoError(err)
		require.True(ok)
	})

	suite.Run("testing passenger subscribe on valid outgoing call node topic", func() {
		ok, err := suite.Authenticator.ACL(acl.Pub, token, validPassengerNodeCallEntryTopic)
		require.NoError(err)
		require.True(ok)
	})

	suite.Run("testing passenger subscribe on invalid call outgoing topic", func() {
		ok, err := suite.Authenticator.ACL(acl.Sub, token, invalidPassengerCallOutgoingTopic)
		require.Error(err)
		require.False(ok)
	})
}

// nolint: funlen
func (suite *AuthenticatorTestSuite) TestACL_Driver() {
	require := suite.Require()
	token := suite.Tokens.Driver

	suite.Run("testing driver publish on its location topic", func() {
		ok, err := suite.Authenticator.ACL(acl.Pub, token, validDriverLocationTopic)
		require.NoError(err)
		require.True(ok)
	})

	suite.Run("testing driver publish on invalid location topic", func() {
		ok, err := suite.Authenticator.ACL(acl.Pub, token, invalidDriverLocationTopic)
		require.Error(err)
		require.False(ok)
	})

	suite.Run("testing driver subscribe on invalid cab event topic", func() {
		ok, err := suite.Authenticator.ACL(acl.Sub, token, invalidDriverCabEventTopic)
		require.Error(err)
		require.False(ok)
	})

	suite.Run("testing driver subscribe on valid superapp event topic", func() {
		ok, err := suite.Authenticator.ACL(acl.Sub, token, validDriverSuperappEventTopic)
		require.NoError(err)
		require.True(ok)
	})

	suite.Run("testing driver subscribe on invalid superapp event topic", func() {
		ok, err := suite.Authenticator.ACL(acl.Sub, token, invalidDriverSuperappEventTopic)
		require.Error(err)
		require.False(ok)
	})

	suite.Run("testing driver subscribe on valid shared location topic", func() {
		ok, err := suite.Authenticator.ACL(acl.Sub, token, validDriverSharedTopic)
		require.NoError(err)
		require.True(ok)
	})

	suite.Run("testing driver subscribe on invalid shared location topic", func() {
		ok, err := suite.Authenticator.ACL(acl.Sub, token, invalidDriverSharedTopic)
		require.Error(err)
		require.False(ok)
	})

	suite.Run("testing driver subscribe on valid chat topic", func() {
		ok, err := suite.Authenticator.ACL(acl.Sub, token, validDriverChatTopic)
		require.NoError(err)
		require.True(ok)
	})

	suite.Run("testing driver subscribe on invalid chat topic", func() {
		ok, err := suite.Authenticator.ACL(acl.Sub, token, invalidDriverChatTopic)
		require.Error(err)
		require.False(ok)
	})

	suite.Run("testing driver subscribe on valid call entry topic", func() {
		ok, err := suite.Authenticator.ACL(acl.Pub, token, validDriverCallEntryTopic)
		require.NoError(err)
		require.True(ok)
	})

	suite.Run("testing driver subscribe on invalid call entry topic", func() {
		ok, err := suite.Authenticator.ACL(acl.Pub, token, invalidDriverCallEntryTopic)
		require.Error(err)
		require.False(ok)
	})

	suite.Run("testing driver subscribe on valid call outgoing topic", func() {
		ok, err := suite.Authenticator.ACL(acl.Sub, token, validDriverCallOutgoingTopic)
		require.NoError(err)
		require.True(ok)
	})

	suite.Run("testing driver subscribe on valid call outgoing node topic", func() {
		ok, err := suite.Authenticator.ACL(acl.Pub, token, validDriverNodeCallEntryTopic)
		require.NoError(err)
		require.True(ok)
	})

	suite.Run("testing driver subscribe on invalid call outgoing topic", func() {
		ok, err := suite.Authenticator.ACL(acl.Sub, token, invalidDriverCallOutgoingTopic)
		require.Error(err)
		require.False(ok)
	})
}

// nolint: funlen
func TestManualAuthenticator_validateAccessType(t *testing.T) {
	t.Parallel()

	type fields struct {
		AllowedAccessTypes []acl.AccessType
	}

	type args struct {
		accessType acl.AccessType
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "#1 testing with no allowed access type",
			fields: fields{AllowedAccessTypes: []acl.AccessType{}},
			args:   args{accessType: acl.Sub},
			want:   false,
		},
		{
			name:   "#2 testing with no allowed access type",
			fields: fields{AllowedAccessTypes: []acl.AccessType{}},
			args:   args{accessType: acl.Pub},
			want:   false,
		},
		{
			name:   "#3 testing with no allowed access type",
			fields: fields{AllowedAccessTypes: []acl.AccessType{}},
			args:   args{accessType: acl.PubSub},
			want:   false,
		},
		{
			name:   "#4 testing with one allowed access type",
			fields: fields{AllowedAccessTypes: []acl.AccessType{acl.Pub}},
			args:   args{accessType: acl.Pub},
			want:   true,
		},
		{
			name:   "#5 testing with one allowed access type",
			fields: fields{AllowedAccessTypes: []acl.AccessType{acl.Pub}},
			args:   args{accessType: acl.Sub},
			want:   false,
		},
		{
			name:   "#6 testing with two allowed access type",
			fields: fields{AllowedAccessTypes: []acl.AccessType{acl.Pub, acl.Sub}},
			args:   args{accessType: acl.Sub},
			want:   true,
		},
		{
			name:   "#7 testing with two allowed access type",
			fields: fields{AllowedAccessTypes: []acl.AccessType{acl.Pub, acl.Sub}},
			args:   args{accessType: acl.Pub},
			want:   true,
		},
		{
			name:   "#8 testing with two allowed access type",
			fields: fields{AllowedAccessTypes: []acl.AccessType{acl.Pub, acl.Sub}},
			args:   args{accessType: acl.PubSub},
			want:   false,
		},
		{
			name:   "#9 testing with three allowed access type",
			fields: fields{AllowedAccessTypes: []acl.AccessType{acl.Pub, acl.Sub, acl.PubSub}},
			args:   args{accessType: acl.PubSub},
			want:   true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// nolint: exhaustruct
			a := authenticator.ManualAuthenticator{
				AllowedAccessTypes: tt.fields.AllowedAccessTypes,
			}
			if got := a.ValidateAccessType(tt.args.accessType); got != tt.want {
				t.Errorf("validateAccessType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getPublicKey(u string) (*rsa.PublicKey, error) {
	var fileName string

	switch u {
	case "1":
		fileName = "../../test/snapp-1.pem"
	case "0":
		fileName = "../../test/snapp-0.pem"
	case "admin":
		fileName = "../../test/snapp-admin.pem"
	default:
		return nil, ErrPublicKeyNotFound
	}

	pem, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("reading public key failed %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pem)
	if err != nil {
		return nil, fmt.Errorf("paring public key failed %w", err)
	}

	return publicKey, nil
}

func getPrivateKey(u string) (*rsa.PrivateKey, error) {
	var fileName string

	switch u {
	case "0":
		fileName = "../../test/snapp-0.private.pem"
	case "1":
		fileName = "../../test/snapp-1.private.pem"
	case "admin":
		fileName = "../../test/snapp-admin.private.pem"
	default:
		return nil, ErrPrivateKeyNotFound
	}

	pem, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("reading private key failed %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(pem)
	if err != nil {
		return nil, fmt.Errorf("paring private key failed %w", err)
	}

	return privateKey, nil
}

func getSampleToken(issuer string, key *rsa.PrivateKey) (string, error) {
	exp := time.Now().Add(time.Hour * 24 * 365 * 10)
	sub := "DXKgaNQa7N5Y7bo"

	// nolint: exhaustruct
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(exp),
		Issuer:    issuer,
		Subject:   sub,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("cannot generate a signed string %w", err)
	}

	return tokenString, nil
}
