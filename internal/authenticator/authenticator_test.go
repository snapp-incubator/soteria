package authenticator_test

import (
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.snapp.ir/dispatching/snappids/v2"
	"gitlab.snapp.ir/dispatching/soteria/internal/authenticator"
	"gitlab.snapp.ir/dispatching/soteria/internal/config"
	"gitlab.snapp.ir/dispatching/soteria/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/pkg/acl"
	"gitlab.snapp.ir/dispatching/soteria/pkg/user"
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

	driverToken, err := suite.getSampleToken(user.Driver, false)
	require.NoError(err)

	suite.Tokens.Driver = driverToken

	passengerToken, err := suite.getSampleToken(user.Passenger, false)
	require.NoError(err)

	suite.Tokens.Passenger = passengerToken

	pkey0, err := suite.getPublicKey(user.Driver)
	require.NoError(err)

	suite.PublicKeys.Driver = pkey0

	pkey1, err := suite.getPublicKey(user.Passenger)
	require.NoError(err)

	suite.PublicKeys.Passenger = pkey1

	hid := &snappids.HashIDSManager{
		Salts: map[snappids.Audience]string{
			snappids.PassengerAudience:  "secret",
			snappids.DriverAudience:     "secret",
			snappids.ThirdPartyAudience: "secret",
		},
		Lengths: map[snappids.Audience]int{
			snappids.PassengerAudience:  15,
			snappids.DriverAudience:     15,
			snappids.ThirdPartyAudience: 15,
		},
	}

	cfg := config.SnappVendor()
	suite.Authenticator = authenticator.Authenticator{
		PublicKeys: map[string]*rsa.PublicKey{
			user.Driver:    pkey0,
			user.Passenger: pkey1,
		},
		AllowedAccessTypes: []acl.AccessType{acl.Pub, acl.Sub, acl.PubSub},
		Company:            "snapp",
		TopicManager:       topics.NewTopicManager(cfg.Topics, hid, "snapp", cfg.IssEntityMap, cfg.IssPeerMap),
	}
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
		require.Equal("token is invalid illegal base64 data at input byte 36", err.Error())
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
		suite.NoError(err)
		suite.True(ok)
	})

	suite.Run("testing driver subscribe on invalid superapp event topic", func() {
		ok, err := suite.Authenticator.ACL(acl.Sub, token, invalidDriverSuperappEventTopic)
		suite.Error(err)
		suite.False(ok)
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

func TestAuthenticator_ValidateTopicBySender(t *testing.T) {
	t.Parallel()

	hid := &snappids.HashIDSManager{
		Salts: map[snappids.Audience]string{
			snappids.DriverAudience: "secret",
		},
		Lengths: map[snappids.Audience]int{
			snappids.DriverAudience: 15,
		},
	}

	cfg := config.SnappVendor()

	// nolint: exhaustruct
	authenticator := authenticator.Authenticator{
		AllowedAccessTypes: []acl.AccessType{acl.Pub, acl.Sub},
		Company:            "snapp",
		TopicManager:       topics.NewTopicManager(cfg.Topics, hid, "snapp", cfg.IssEntityMap, cfg.IssPeerMap),
	}

	t.Run("testing valid driver cab event", func(t *testing.T) {
		topicTemplate := authenticator.TopicManager.ParseTopic(validDriverCabEventTopic, user.Driver, "DXKgaNQa7N5Y7bo")
		assert.True(t, topicTemplate != nil)
	})
}

// nolint: funlen
func TestAuthenticator_validateAccessType(t *testing.T) {
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
			a := authenticator.Authenticator{
				AllowedAccessTypes: tt.fields.AllowedAccessTypes,
			}
			if got := a.ValidateAccessType(tt.args.accessType); got != tt.want {
				t.Errorf("validateAccessType() = %v, want %v", got, tt.want)
			}
		})
	}
}

// nolint: goerr113, wrapcheck
func (suite *AuthenticatorTestSuite) getPublicKey(u string) (*rsa.PublicKey, error) {
	var fileName string

	switch u {
	case user.Passenger:
		fileName = "../../test/1.pem"
	case user.Driver:
		fileName = "../../test/0.pem"
	case user.None:
		fallthrough
	default:
		return nil, errors.New("invalid user, public key not found")
	}

	pem, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pem)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}

// nolint: goerr113, wrapcheck
func (suite *AuthenticatorTestSuite) getPrivateKey(u string) (*rsa.PrivateKey, error) {
	var fileName string

	switch u {
	case user.Driver:
		fileName = "../../test/0.private.pem"
	case user.Passenger:
		fileName = "../../test/1.private.pem"
	case user.None:
		fallthrough
	default:
		return nil, errors.New("invalid user, private key not found")
	}

	pem, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(pem)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func (suite *AuthenticatorTestSuite) getSampleToken(issuer string, isSuperuser bool) (string, error) {
	key, err := suite.getPrivateKey(issuer)
	if err != nil {
		suite.Require().NoError(err)
	}

	exp := time.Now().Add(time.Hour * 24 * 365 * 10)
	sub := "DXKgaNQa7N5Y7bo"

	// nolint: exhaustruct
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(exp),
		Issuer:    string(issuer),
		Subject:   sub,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(key)
	if err != nil {
		suite.Require().NoError(err)
	}

	return tokenString, nil
}
