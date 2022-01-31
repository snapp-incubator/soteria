package authenticator_test

import (
	"crypto/rsa"
	"fmt"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/config"
	"io/ioutil"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"gitlab.snapp.ir/dispatching/snappids/v2"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/authenticator"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/acl"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/user"
)

const (
	invalidToken                  = "ey1JhbGciOiJSUzI1NiIsInR5cCI56kpXVCJ9.eyJzdWIiOiJCRzdScDFkcnpWRE5RcjYiLCJuYW1lIjoiSm9obiBEb2UiLCJhZG1pbiI6dHJ1ZSwiaXNzIjowLCJpYXQiOjE1MTYyMzkwMjJ9.1cYXFEhcewOYFjGJYhB8dsaFO9uKEXwlM8954rkt4Tsu0lWMITbRf_hHh1l9QD4MFqD-0LwRPUYaiaemy0OClMu00G2sujLCWaquYDEP37iIt8RoOQAh8Jb5vT8LX5C3PEKvbW_i98u8HHJoFUR9CXJmzrKi48sAcOYvXVYamN0S9KoY38H-Ze37Mdu3o6B58i73krk7QHecsc2_PkCJisvUVAzb0tiInIalBc8-zI3QZSxwNLr_hjlBg1sUxTUvH5SCcRR7hxI8TxJzkOHqAHWDRO84NC_DSAoO2p04vrHpqglN9XPJ8RC2YWpfefvD2ttH554RJWu_0RlR2kAYvQ"
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

func TestAuthenticator_Auth(t *testing.T) {
	driverToken, err := getSampleToken(user.Driver, false)
	if err != nil {
		t.Fatal(err)
	}

	passengerToken, err := getSampleToken(user.Passenger, false)
	if err != nil {
		t.Fatal(err)
	}

	pkey0, err := getPublicKey(user.Driver)
	if err != nil {
		t.Fatal(err)
	}

	pkey1, err := getPublicKey(user.Passenger)
	if err != nil {
		t.Fatal(err)
	}

	// nolint: exhaustivestruct
	authenticator := authenticator.Authenticator{
		PublicKeys: map[user.Issuer]*rsa.PublicKey{
			user.Driver:    pkey0,
			user.Passenger: pkey1,
		},
		ModelHandler: MockModelHandler{},
	}

	t.Run("testing driver token auth", func(t *testing.T) {
		err := authenticator.Auth(driverToken)
		assert.NoError(t, err)
	})

	t.Run("testing passenger token auth", func(t *testing.T) {
		err := authenticator.Auth(passengerToken)
		assert.NoError(t, err)
	})

	t.Run("testing invalid token auth", func(t *testing.T) {
		err := authenticator.Auth(invalidToken)
		assert.Error(t, err)
	})
}

func TestAuthenticator_Acl(t *testing.T) {
	pkey0, err := getPublicKey(user.Driver)
	if err != nil {
		t.Fatal(err)
	}
	pkey1, err := getPublicKey(user.Passenger)
	if err != nil {
		t.Fatal(err)
	}
	passengerToken, err := getSampleToken(user.Passenger, false)
	if err != nil {
		t.Fatal(err)
	}
	driverToken, err := getSampleToken(user.Driver, false)
	if err != nil {
		t.Fatal(err)
	}

	cfg := config.New()

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

	auth := authenticator.Authenticator{
		PublicKeys: map[user.Issuer]*rsa.PublicKey{
			user.Driver:    pkey0,
			user.Passenger: pkey1,
		},
		AllowedAccessTypes: []acl.AccessType{acl.Pub, acl.Sub, acl.PubSub},
		ModelHandler:       MockModelHandler{},
		Company:            "snapp",
		TopicManager:       topics.NewTopicManager(cfg.Topics, hid, "snapp"),
	}
	t.Run("testing acl with invalid access type", func(t *testing.T) {
		ok, err := auth.ACL("invalid-access", passengerToken, "test")
		assert.Error(t, err)
		assert.False(t, ok)
		assert.Equal(t, authenticator.ErrInvalidAccessType.Error(), err.Error())
	})
	t.Run("testing acl with invalid token", func(t *testing.T) {
		ok, err := auth.ACL(acl.Pub, invalidToken, validDriverCabEventTopic)
		assert.False(t, ok)
		assert.Error(t, err)
		assert.Equal(t, "token is invalid illegal base64 data at input byte 36", err.Error())
	})
	t.Run("testing acl with valid inputs", func(t *testing.T) {
		ok, err := auth.ACL(acl.Sub, passengerToken, validPassengerCabEventTopic)
		assert.NoError(t, err)
		assert.True(t, ok)
	})
	t.Run("testing acl with invalid topic", func(t *testing.T) {
		ok, err := auth.ACL(acl.Sub, passengerToken, invalidPassengerCabEventTopic)
		assert.Error(t, err)
		assert.False(t, ok)
	})
	t.Run("testing acl with invalid access type", func(t *testing.T) {
		ok, err := auth.ACL(acl.Pub, passengerToken, validPassengerCabEventTopic)
		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("testing driver publish on its location topic", func(t *testing.T) {
		ok, err := auth.ACL(acl.Pub, driverToken, validDriverLocationTopic)
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("testing driver publish on invalid location topic", func(t *testing.T) {
		ok, err := auth.ACL(acl.Pub, driverToken, invalidDriverLocationTopic)
		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("testing driver subscribe on invalid cab event topic", func(t *testing.T) {
		ok, err := auth.ACL(acl.Sub, driverToken, invalidDriverCabEventTopic)
		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("testing passenger subscribe on valid superapp event topic", func(t *testing.T) {
		ok, err := auth.ACL(acl.Sub, passengerToken, validPassengerSuperappEventTopic)
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("testing passenger subscribe on invalid superapp event topic", func(t *testing.T) {
		ok, err := auth.ACL(acl.Sub, passengerToken, invalidPassengerSuperappEventTopic)
		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("testing driver subscribe on valid superapp event topic", func(t *testing.T) {
		ok, err := auth.ACL(acl.Sub, driverToken, validDriverSuperappEventTopic)
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("testing driver subscribe on invalid superapp event topic", func(t *testing.T) {
		ok, err := auth.ACL(acl.Sub, driverToken, invalidDriverSuperappEventTopic)
		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("testing driver subscribe on valid shared location topic", func(t *testing.T) {
		ok, err := auth.ACL(acl.Sub, driverToken, validDriverSharedTopic)
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("testing passenger subscribe on valid shared location topic", func(t *testing.T) {
		ok, err := auth.ACL(acl.Sub, passengerToken, validPassengerSharedTopic)
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("testing driver subscribe on invalid shared location topic", func(t *testing.T) {
		ok, err := auth.ACL(acl.Sub, driverToken, invalidDriverSharedTopic)
		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("testing passenger subscribe on invalid shared location topic", func(t *testing.T) {
		ok, err := auth.ACL(acl.Sub, passengerToken, invalidPassengerSharedTopic)
		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("testing driver subscribe on valid chat topic", func(t *testing.T) {
		ok, err := auth.ACL(acl.Sub, driverToken, validDriverChatTopic)
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("testing passenger subscribe on valid chat topic", func(t *testing.T) {
		ok, err := auth.ACL(acl.Sub, passengerToken, validPassengerChatTopic)
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("testing driver subscribe on invalid chat topic", func(t *testing.T) {
		ok, err := auth.ACL(acl.Sub, driverToken, invalidDriverChatTopic)
		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("testing passenger subscribe on invalid chat topic", func(t *testing.T) {
		ok, err := auth.ACL(acl.Sub, passengerToken, invalidPassengerChatTopic)
		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("testing driver subscribe on valid call entry topic", func(t *testing.T) {
		ok, err := auth.ACL(acl.Pub, driverToken, validDriverCallEntryTopic)
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("testing passenger subscribe on valid entry call topic", func(t *testing.T) {
		ok, err := auth.ACL(acl.Pub, passengerToken, validPassengerCallEntryTopic)
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("testing driver subscribe on invalid call entry topic", func(t *testing.T) {
		ok, err := auth.ACL(acl.Pub, driverToken, invalidDriverCallEntryTopic)
		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("testing passenger subscribe on invalid call entry topic", func(t *testing.T) {
		ok, err := auth.ACL(acl.Pub, passengerToken, invalidPassengerCallEntryTopic)
		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("testing driver subscribe on valid call outgoing topic", func(t *testing.T) {
		ok, err := auth.ACL(acl.Sub, driverToken, validDriverCallOutgoingTopic)
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("testing passenger subscribe on valid outgoing call topic", func(t *testing.T) {
		ok, err := auth.ACL(acl.Sub, passengerToken, validPassengerCallOutgoingTopic)
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("testing driver subscribe on valid call outgoing node topic", func(t *testing.T) {
		ok, err := auth.ACL(acl.Pub, driverToken, validDriverNodeCallEntryTopic)
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("testing passenger subscribe on valid outgoing call node topic", func(t *testing.T) {
		ok, err := auth.ACL(acl.Pub, passengerToken, validPassengerNodeCallEntryTopic)
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("testing driver subscribe on invalid call outgoing topic", func(t *testing.T) {
		ok, err := auth.ACL(acl.Sub, driverToken, invalidDriverCallOutgoingTopic)
		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("testing passenger subscribe on invalid call outgoing topic", func(t *testing.T) {
		ok, err := auth.ACL(acl.Sub, passengerToken, invalidPassengerCallOutgoingTopic)
		assert.Error(t, err)
		assert.False(t, ok)
	})
}

func TestAuthenticator_ValidateTopicBySender(t *testing.T) {
	hid := &snappids.HashIDSManager{
		Salts: map[snappids.Audience]string{
			snappids.DriverAudience: "secret",
		},
		Lengths: map[snappids.Audience]int{
			snappids.DriverAudience: 15,
		},
	}

	cfg := config.New()

	// nolint: exhaustivestruct
	authenticator := authenticator.Authenticator{
		AllowedAccessTypes: []acl.AccessType{acl.Pub, acl.Sub},
		ModelHandler:       MockModelHandler{},
		Company:            "snapp",
		TopicManager:       topics.NewTopicManager(cfg.Topics, hid, "snapp"),
	}

	t.Run("testing valid driver cab event", func(t *testing.T) {
		audience, audienceStr := topics.IssuerToAudience(user.Driver)
		topicTemplate := authenticator.TopicManager.ValidateTopic(validDriverCabEventTopic, audienceStr, audience, "DXKgaNQa7N5Y7bo")
		assert.True(t, topicTemplate != nil)
	})
}

func TestAuthenticator_validateAccessType(t *testing.T) {
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

			// nolint: exhaustivestruct
			a := authenticator.Authenticator{
				AllowedAccessTypes: tt.fields.AllowedAccessTypes,
			}
			if got := a.ValidateAccessType(tt.args.accessType); got != tt.want {
				t.Errorf("validateAccessType() = %v, want %v", got, tt.want)
			}
		})
	}
}

type MockModelHandler struct{}

func (rmh MockModelHandler) Get(pk string) user.User {
	var u user.User

	switch pk {
	case "passenger":
		u = user.User{
			Username: string(user.Passenger),
			Rules: []user.Rule{
				{
					Topic:  topics.CabEvent,
					Access: acl.Sub,
				},
				{
					Topic:  topics.SuperappEvent,
					Access: acl.Sub,
				},
				{
					Topic:  topics.PassengerLocation,
					Access: acl.Pub,
				},
				{
					Topic:  topics.SharedLocation,
					Access: acl.Sub,
				},
				{
					Topic:  topics.Chat,
					Access: acl.Sub,
				},
				{
					Topic:  topics.GeneralCallEntry,
					Access: acl.Pub,
				},
				{
					Topic:  topics.NodeCallEntry,
					Access: acl.Pub,
				},
				{
					Topic:  topics.CallOutgoing,
					Access: acl.Sub,
				},
			},
		}
	case "driver":
		u = user.User{
			Username: string(user.Driver),
			Rules: []user.Rule{
				{
					Topic:  topics.DriverLocation,
					Access: acl.Pub,
				},
				{
					Topic:  topics.CabEvent,
					Access: acl.Sub,
				},
				{
					Topic:  topics.SuperappEvent,
					Access: acl.Sub,
				},
				{
					Topic:  topics.PassengerLocation,
					Access: acl.Pub,
				},
				{
					Topic:  topics.SharedLocation,
					Access: acl.Sub,
				},
				{
					Topic:  topics.Chat,
					Access: acl.Sub,
				},
				{
					Topic:  topics.GeneralCallEntry,
					Access: acl.Pub,
				},
				{
					Topic:  topics.NodeCallEntry,
					Access: acl.Pub,
				},
				{
					Topic:  topics.CallOutgoing,
					Access: acl.Sub,
				},
			},
		}
	}

	return u
}

func getPublicKey(u user.Issuer) (*rsa.PublicKey, error) {
	var fileName string

	switch u {
	case user.Passenger:
		fileName = "../../test/1.pem"
	case user.Driver:
		fileName = "../../test/0.pem"
	default:
		return nil, fmt.Errorf("invalid user, public key not found")
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

func getPrivateKey(u user.Issuer) (*rsa.PrivateKey, error) {
	var fileName string
	switch u {
	case user.Driver:
		fileName = "../../test/0.private.pem"
	case user.Passenger:
		fileName = "../../test/1.private.pem"
	default:
		return nil, fmt.Errorf("invalid user, private key not found")
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

func getSampleToken(issuer user.Issuer, isSuperuser bool) (string, error) {
	key, err := getPrivateKey(issuer)
	if err != nil {
		panic(err)
	}

	exp := time.Now().Add(time.Hour * 24 * 365 * 10).Unix()
	sub := "DXKgaNQa7N5Y7bo"

	var claims jwt.Claims
	if isSuperuser {
		claims = jwt.MapClaims{
			"exp":          exp,
			"iss":          string(issuer),
			"sub":          sub,
			"is_superuser": true,
		}
	} else {
		claims = jwt.StandardClaims{
			ExpiresAt: exp,
			Issuer:    string(issuer),
			Subject:   sub,
		}
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		panic(err)
	}
	return tokenString, nil
}
