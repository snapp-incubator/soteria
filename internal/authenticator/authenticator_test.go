package authenticator

import (
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	snappids "gitlab.snapp.ir/dispatching/snappids/v2"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/db"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/acl"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/memoize"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/user"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"testing"
	"time"
)

const (
	invalidToken                  = "ey1JhbGciOiJSUzI1NiIsInR5cCI56kpXVCJ9.eyJzdWIiOiJCRzdScDFkcnpWRE5RcjYiLCJuYW1lIjoiSm9obiBEb2UiLCJhZG1pbiI6dHJ1ZSwiaXNzIjowLCJpYXQiOjE1MTYyMzkwMjJ9.1cYXFEhcewOYFjGJYhB8dsaFO9uKEXwlM8954rkt4Tsu0lWMITbRf_hHh1l9QD4MFqD-0LwRPUYaiaemy0OClMu00G2sujLCWaquYDEP37iIt8RoOQAh8Jb5vT8LX5C3PEKvbW_i98u8HHJoFUR9CXJmzrKi48sAcOYvXVYamN0S9KoY38H-Ze37Mdu3o6B58i73krk7QHecsc2_PkCJisvUVAzb0tiInIalBc8-zI3QZSxwNLr_hjlBg1sUxTUvH5SCcRR7hxI8TxJzkOHqAHWDRO84NC_DSAoO2p04vrHpqglN9XPJ8RC2YWpfefvD2ttH554RJWu_0RlR2kAYvQ"
	validPassengerCabEventTopic   = "passenger-event-152384980615c2bd16143cff29038b67"
	invalidPassengerCabEventTopic = "passenger-event-152384980615c2bd16156cff29038b67"

	validDriverCabEventTopic   = "driver-event-152384980615c2bd16143cff29038b67"
	invalidDriverCabEventTopic = "driver-event-152384980615c2bd16156cff29038b67"

	validDriverLocationTopic   = "snapp/driver/DXKgaNQa7N5Y7bo/location"
	invalidDriverLocationTopic = "snapp/driver/DXKgaNQa9Q5Y7bo/location"

	validPassengerSuperappEventTopic   = "snapp/passenger/0956923be632d673560af9adadd2f78a/superapp"
	invalidPassengerSuperappEventTopic = "snapp/passenger/0959623be632d673560af9adadd2f78a/superapp"

	validDriverSuperappEventTopic   = "snapp/driver/0956923be632d673560af9adadd2f78a/superapp"
	invalidDriverSuperappEventTopic = "snapp/driver/0596923be632d673560af9adadd2f78a/superapp"
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
	thirdPartyToken, err := getSampleToken(user.ThirdParty, false)
	if err != nil {
		t.Fatal(err)
	}
	superuserToken, err := getSampleToken(user.ThirdParty, true)
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
	pkey100, err := getPublicKey(user.ThirdParty)
	if err != nil {
		t.Fatal(err)
	}
	key100, err := getPrivateKey(user.ThirdParty)
	if err != nil {
		t.Fatal(err)
	}
	passwordChecker := memoize.MemoizedCompareHashAndPassword()
	authenticator := Authenticator{
		PrivateKeys: map[user.Issuer]*rsa.PrivateKey{
			user.ThirdParty: key100,
		},
		PublicKeys: map[user.Issuer]*rsa.PublicKey{
			user.Driver:     pkey0,
			user.Passenger:  pkey1,
			user.ThirdParty: pkey100,
		},
		ModelHandler:           MockModelHandler{},
		CompareHashAndPassword: passwordChecker,
	}
	t.Run("testing driver token auth", func(t *testing.T) {
		ok, err := authenticator.Auth(context.Background(), driverToken)
		assert.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("testing passenger token auth", func(t *testing.T) {
		ok, err := authenticator.Auth(context.Background(), passengerToken)
		assert.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("testing third party token auth", func(t *testing.T) {
		ok, err := authenticator.Auth(context.Background(), thirdPartyToken)
		assert.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("testing superuser token auth", func(t *testing.T) {
		ok, err := authenticator.Auth(context.Background(), superuserToken)
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("testing invalid token auth", func(t *testing.T) {
		ok, err := authenticator.Auth(context.Background(), invalidToken)
		assert.Error(t, err)
		assert.False(t, ok)
	})
}

func TestAuthenticator_Token(t *testing.T) {
	key, err := getPrivateKey(user.ThirdParty)
	if err != nil {
		t.Fatal(err)
	}
	pk, err := getPublicKey(user.ThirdParty)
	if err != nil {
		t.Fatal(err)
	}
	passwordChecker := memoize.MemoizedCompareHashAndPassword()

	authenticator := Authenticator{
		PrivateKeys: map[user.Issuer]*rsa.PrivateKey{
			user.ThirdParty: key,
		},
		ModelHandler:           MockModelHandler{},
		CompareHashAndPassword: passwordChecker,
	}
	t.Run("testing getting token with valid inputs", func(t *testing.T) {
		tokenString, err := authenticator.Token(context.Background(), acl.ClientCredentials, "snappbox", "KJIikjIKbIYVGj)YihYUGIB&")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return pk, nil
		})
		assert.NoError(t, err)
		claims := token.Claims.(jwt.MapClaims)
		assert.Equal(t, "snappbox", claims["sub"].(string))
		assert.Equal(t, "100", claims["iss"].(string))

	})
	t.Run("testing getting token with valid inputs", func(t *testing.T) {
		tokenString, err := authenticator.Token(context.Background(), acl.ClientCredentials, "snappbox", "invalid secret")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return pk, nil
		})
		assert.Error(t, err)
		assert.Nil(t, token)
	})
}

func TestAuthenticator_HeraldToken(t *testing.T) {
	key, err := getPrivateKey(user.ThirdParty)
	if err != nil {
		t.Fatal(err)
	}
	pk, err := getPublicKey(user.ThirdParty)
	if err != nil {
		t.Fatal(err)
	}

	authenticator := Authenticator{
		PrivateKeys: map[user.Issuer]*rsa.PrivateKey{
			user.ThirdParty: key,
		},
	}

	t.Run("testing issuing valid herald token", func(t *testing.T) {
		allowedTopics := []acl.Topic{
			{
				Type: topics.BoxEvent,
			},
		}

		allowedEndpoints := []acl.Endpoint{
			{
				Name: "event",
			},
		}

		tokenString, err := authenticator.HeraldToken(
			"snappbox", allowedEndpoints, allowedTopics, time.Hour*24)
		assert.NoError(t, err)
		token, err := jwt.ParseWithClaims(tokenString, &acl.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return pk, nil
		})
		assert.NoError(t, err)
		actual := token.Claims.(*acl.Claims)

		expected := acl.Claims{
			StandardClaims: jwt.StandardClaims{
				Issuer:  "100",
				Subject: "snappbox",
			},
			Topics: []acl.Topic{
				{
					Type: topics.BoxEvent,
				},
			},
			Endpoints: []acl.Endpoint{
				{
					Name: "event",
				},
			},
		}
		assert.Equal(t, expected.StandardClaims.Issuer, actual.StandardClaims.Issuer)
		assert.Equal(t, expected.StandardClaims.Subject, actual.StandardClaims.Subject)
		assert.EqualValues(t, expected.Endpoints, actual.Endpoints)
		assert.EqualValues(t, expected.Topics, actual.Topics)
	})
}

func TestAuthenticator_SuperuserToken(t *testing.T) {
	key, err := getPrivateKey(user.ThirdParty)
	if err != nil {
		t.Fatal(err)
	}
	pk, err := getPublicKey(user.ThirdParty)
	if err != nil {
		t.Fatal(err)
	}
	authenticator := Authenticator{
		PrivateKeys: map[user.Issuer]*rsa.PrivateKey{
			user.ThirdParty: key,
		},
	}

	tokenString, err := authenticator.SuperuserToken("herald", time.Hour*24)
	assert.NoError(t, err)

	token, err := jwt.ParseWithClaims(tokenString, &acl.SuperuserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return pk, nil
	})
	assert.NoError(t, err)

	actual := token.Claims.(*acl.SuperuserClaims)
	expected := acl.SuperuserClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:  "100",
			Subject: "herald",
		},
		IsSuperuser: true,
	}

	assert.Equal(t, expected.StandardClaims.Issuer, actual.StandardClaims.Issuer)
	assert.Equal(t, expected.StandardClaims.Subject, actual.StandardClaims.Subject)
	assert.Equal(t, expected.IsSuperuser, actual.IsSuperuser)
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
	pkey100, err := getPublicKey(user.ThirdParty)
	if err != nil {
		t.Fatal(err)
	}
	key100, err := getPrivateKey(user.ThirdParty)
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
	thirdPartyToken, err := getSampleToken(user.ThirdParty, false)
	if err != nil {
		t.Fatal(err)
	}

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

	passwordChecker := memoize.MemoizedCompareHashAndPassword()
	authenticator := Authenticator{
		PrivateKeys: map[user.Issuer]*rsa.PrivateKey{
			user.ThirdParty: key100,
		},
		PublicKeys: map[user.Issuer]*rsa.PublicKey{
			user.Driver:     pkey0,
			user.Passenger:  pkey1,
			user.ThirdParty: pkey100,
		},
		AllowedAccessTypes:     []acl.AccessType{acl.Pub, acl.Sub},
		ModelHandler:           MockModelHandler{},
		EMQTopicManager:        snappids.NewEMQManager(hid),
		HashIDSManager:         hid,
		CompareHashAndPassword: passwordChecker,
	}
	t.Run("testing acl with invalid access type", func(t *testing.T) {
		ok, err := authenticator.Acl(context.Background(), acl.PubSub, passengerToken, "test")
		assert.Error(t, err)
		assert.False(t, ok)
		assert.Equal(t, "requested access type publish-subscribe is invalid", err.Error())
	})
	t.Run("testing acl with invalid token", func(t *testing.T) {
		ok, err := authenticator.Acl(context.Background(), acl.Pub, invalidToken, validDriverCabEventTopic)
		assert.False(t, ok)
		assert.Error(t, err)
		assert.Equal(t, "token is invalid. err: illegal base64 data at input byte 37", err.Error())
	})
	t.Run("testing acl with valid inputs", func(t *testing.T) {
		ok, err := authenticator.Acl(context.Background(), acl.Sub, passengerToken, validPassengerCabEventTopic)
		assert.NoError(t, err)
		assert.True(t, ok)
	})
	t.Run("testing acl with invalid topic", func(t *testing.T) {
		ok, err := authenticator.Acl(context.Background(), acl.Sub, passengerToken, invalidPassengerCabEventTopic)
		assert.Error(t, err)
		assert.False(t, ok)
	})
	t.Run("testing acl with invalid access type", func(t *testing.T) {
		ok, err := authenticator.Acl(context.Background(), acl.Pub, passengerToken, validPassengerCabEventTopic)
		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("testing acl with third party token", func(t *testing.T) {
		ok, err := authenticator.Acl(context.Background(), acl.Sub, thirdPartyToken, validDriverLocationTopic)
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("testing driver publish on its location topic", func(t *testing.T) {
		ok, err := authenticator.Acl(context.Background(), acl.Pub, driverToken, validDriverLocationTopic)
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("testing driver publish on invalid location topic", func(t *testing.T) {
		ok, err := authenticator.Acl(context.Background(), acl.Pub, driverToken, invalidDriverLocationTopic)
		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("testing driver subscribe on invalid cab event topic", func(t *testing.T) {
		ok, err := authenticator.Acl(context.Background(), acl.Sub, driverToken, invalidDriverCabEventTopic)
		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("testing passenger subscribe on valid superapp event topic", func(t *testing.T) {
		ok, err := authenticator.Acl(context.Background(), acl.Sub, passengerToken, validPassengerSuperappEventTopic)
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("testing passenger subscribe on invalid superapp event topic", func(t *testing.T) {
		ok, err := authenticator.Acl(context.Background(), acl.Sub, passengerToken, invalidPassengerSuperappEventTopic)
		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("testing driver subscribe on valid superapp event topic", func(t *testing.T) {
		ok, err := authenticator.Acl(context.Background(), acl.Sub, driverToken, validDriverSuperappEventTopic)
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("testing driver subscribe on invalid superapp event topic", func(t *testing.T) {
		ok, err := authenticator.Acl(context.Background(), acl.Sub, driverToken, invalidDriverSuperappEventTopic)
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

	passwordChecker := memoize.MemoizedCompareHashAndPassword()
	authenticator := Authenticator{
		AllowedAccessTypes:     []acl.AccessType{acl.Pub, acl.Sub},
		ModelHandler:           MockModelHandler{},
		EMQTopicManager:        snappids.NewEMQManager(hid),
		HashIDSManager:         hid,
		CompareHashAndPassword: passwordChecker,
	}

	t.Run("testing valid driver cab event", func(t *testing.T) {
		ok := authenticator.ValidateTopicBySender(validDriverCabEventTopic, snappids.DriverAudience, 123)
		assert.True(t, ok)
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
		t.Run(tt.name, func(t *testing.T) {
			a := Authenticator{
				AllowedAccessTypes: tt.fields.AllowedAccessTypes,
			}
			if got := a.validateAccessType(tt.args.accessType); got != tt.want {
				t.Errorf("validateAccessType() = %v, want %v", got, tt.want)
			}
		})
	}
}

type MockModelHandler struct{}

func (rmh MockModelHandler) Save(ctx context.Context, model db.Model) error {
	return nil
}

func (rmh MockModelHandler) Delete(ctx context.Context, modelName, pk string) error {
	return nil
}

func (rmh MockModelHandler) Get(ctx context.Context, modelName, pk string, v interface{}) error {
	switch pk {
	case "passenger":
		*v.(*user.User) = user.User{
			MetaData: db.MetaData{},
			Username: string(user.Passenger),
			Type:     user.EMQUser,
			Rules: []user.Rule{
				user.Rule{
					UUID:       uuid.New(),
					Topic:      topics.CabEvent,
					AccessType: acl.Sub,
				},
				user.Rule{
					UUID:       uuid.New(),
					Topic:      topics.SuperappEvent,
					AccessType: acl.Sub,
				},
			},
		}
	case "driver":
		*v.(*user.User) = user.User{
			MetaData: db.MetaData{},
			Username: string(user.Driver),
			Type:     user.EMQUser,
			Rules: []user.Rule{{
				UUID:       uuid.Nil,
				Endpoint:   "",
				Topic:      topics.DriverLocation,
				AccessType: acl.Pub,
			}, {
				UUID:       uuid.Nil,
				Endpoint:   "",
				Topic:      topics.CabEvent,
				AccessType: acl.Sub,
			},
				{
					UUID:       uuid.New(),
					Topic:      topics.SuperappEvent,
					AccessType: acl.Sub,
				},
			},
		}
	case "snappbox":
		*v.(*user.User) = user.User{
			MetaData:                db.MetaData{},
			Username:                "snapp-box",
			Password:                getSamplePassword(),
			Type:                    user.HeraldUser,
			Secret:                  "KJIikjIKbIYVGj)YihYUGIB&",
			TokenExpirationDuration: time.Hour * 72,
			Rules: []user.Rule{
				{
					UUID:       uuid.New(),
					Topic:      topics.BoxEvent,
					AccessType: acl.Sub,
				},
				{
					UUID:       uuid.New(),
					Endpoint:   "/notification",
					AccessType: acl.Pub,
				},
			},
		}
	case "colony-subscriber":
		*v.(*user.User) = user.User{
			MetaData:                db.MetaData{},
			Username:                "colony-subscriber",
			Password:                "password",
			Type:                    user.EMQUser,
			Secret:                  "secret",
			TokenExpirationDuration: 0,
			Rules: []user.Rule{
				user.Rule{
					UUID:       uuid.New(),
					Topic:      topics.DriverLocation,
					AccessType: acl.Sub,
				},
			},
		}
	}
	return nil
}

func (rmh MockModelHandler) Update(ctx context.Context, model db.Model) error {
	return nil
}

func getPublicKey(u user.Issuer) (*rsa.PublicKey, error) {
	var fileName string
	switch u {
	case user.Passenger:
		fileName = "../../test/1.pem"
	case user.Driver:
		fileName = "../../test/0.pem"
	case user.ThirdParty:
		fileName = "../../test/100.pem"
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
	case user.ThirdParty:
		fileName = "../../test/100.private.pem"
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
	if issuer == user.ThirdParty {
		sub = "colony-subscriber"
	}

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

func getSamplePassword() string {
	hash, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	return string(hash)
}
