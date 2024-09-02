package authenticator_test

import (
	"crypto/rsa"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/snapp-incubator/soteria/internal/authenticator"
	"github.com/snapp-incubator/soteria/internal/config"
	"github.com/snapp-incubator/soteria/internal/topics"
	"github.com/snapp-incubator/soteria/pkg/acl"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type ManualAuthenticatorSnappTestSuite struct {
	suite.Suite

	Tokens struct {
		Passenger string
		Driver    string
	}

	PublicKeys struct {
		Passenger *rsa.PublicKey
		Driver    *rsa.PublicKey
	}

	PrivateKeys struct {
		Passenger *rsa.PrivateKey
		Driver    *rsa.PrivateKey
	}

	Authenticator authenticator.Authenticator
}

func TestManualAuthenticator_suite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(ManualAuthenticatorSnappTestSuite))
}

func (suite *ManualAuthenticatorSnappTestSuite) SetupSuite() {
	cfg := config.SnappVendor()

	require := suite.Require()

	pkey0, err := getPublicKey("0")
	require.NoError(err)

	suite.PublicKeys.Driver = pkey0

	pkey1, err := getPublicKey("1")
	require.NoError(err)

	suite.PublicKeys.Passenger = pkey1

	key0, err := getPrivateKey("0")
	require.NoError(err)

	suite.PrivateKeys.Driver = key0

	key1, err := getPrivateKey("1")
	require.NoError(err)

	suite.PrivateKeys.Passenger = key1

	driverToken, err := getSampleToken("0", key0)
	require.NoError(err)

	suite.Tokens.Driver = driverToken

	passengerToken, err := getSampleToken("1", key1)
	require.NoError(err)

	suite.Tokens.Passenger = passengerToken

	hid, err := topics.NewHashIDManager(cfg.HashIDMap)
	require.NoError(err)

	suite.Authenticator = authenticator.ManualAuthenticator{
		Keys: map[string]any{
			topics.DriverIss:    pkey0,
			topics.PassengerIss: pkey1,
		},
		AllowedAccessTypes: []acl.AccessType{acl.Pub, acl.Sub, acl.PubSub},
		Company:            "snapp",
		Parser:             jwt.NewParser(),
		TopicManager:       topics.NewTopicManager(cfg.Topics, hid, "snapp", cfg.IssEntityMap, cfg.IssPeerMap, zap.NewNop()),
		JWTConfig: config.JWT{
			IssName:       "iss",
			SubName:       "sub",
			SigningMethod: "rsa256",
		},
	}
}

func (suite *ManualAuthenticatorSnappTestSuite) TestAuth() {
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

	suite.Run("testing token with invalid iss", func() {
		token, err := getSampleToken("-1", suite.PrivateKeys.Passenger)
		require.NoError(err)

		require.ErrorIs(suite.Authenticator.Auth(token), authenticator.KeyNotFoundError{
			Issuer: "-1",
		})
	})
}

func (suite *ManualAuthenticatorSnappTestSuite) TestACLBasics() {
	require := suite.Require()

	suite.Run("testing acl with invalid access type", func() {
		ok, err := suite.Authenticator.ACL("invalid-access", suite.Tokens.Passenger, "test")
		require.False(ok)
		require.ErrorIs(err, authenticator.ErrInvalidAccessType)
	})

	suite.Run("testing acl with invalid token", func() {
		ok, err := suite.Authenticator.ACL(acl.Pub, invalidToken, validDriverCabEventTopic)
		require.False(ok)
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
func (suite *ManualAuthenticatorSnappTestSuite) TestACLPassenger() {
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
		require.ErrorIs(err, authenticator.InvalidTopicError{
			Topic: invalidPassengerCallEntryTopic,
		})
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
func (suite *ManualAuthenticatorSnappTestSuite) TestACLDriver() {
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

func TestManualAuthenticator_ValidateTopicBySender(t *testing.T) {
	t.Parallel()

	cfg := config.SnappVendor()

	hid, err := topics.NewHashIDManager(cfg.HashIDMap)
	require.NoError(t, err)

	// nolint: exhaustruct
	authenticator := authenticator.ManualAuthenticator{
		AllowedAccessTypes: []acl.AccessType{acl.Pub, acl.Sub},
		Company:            "snapp",
		TopicManager:       topics.NewTopicManager(cfg.Topics, hid, "snapp", cfg.IssEntityMap, cfg.IssPeerMap, zap.NewNop()),
	}

	t.Run("testing valid driver cab event", func(t *testing.T) {
		t.Parallel()

		topicTemplate := authenticator.TopicManager.ParseTopic(
			validDriverCabEventTopic,
			topics.DriverIss,
			"DXKgaNQa7N5Y7bo",
			nil,
		)
		require.NotNil(t, topicTemplate)
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
