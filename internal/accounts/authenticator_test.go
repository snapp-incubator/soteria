package accounts

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"gitlab.snapp.ir/dispatching/snappids"
	"gitlab.snapp.ir/dispatching/soteria/internal/db"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"testing"
	"time"
)

func TestAuthenticator_Auth(t *testing.T) {
	validToken, err := getSampleToken(true)
	if err != nil {
		t.Fatal(err)
	}
	invalidToken, err := getSampleToken(false)
	if err != nil {
		t.Fatal(err)
	}
	key, err := getPrivateKey(ThirdParty)
	if err != nil {
		t.Fatal(err)
	}
	authenticator := Authenticator{
		PrivateKeys: map[string]*rsa.PrivateKey{
			ThirdParty: key,
		},
		ModelHandler: MockModelHandler{},
	}
	t.Run("testing invalid token", func(t *testing.T) {
		ok, err := authenticator.Auth(invalidToken)
		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("testing valid auth", func(t *testing.T) {
		ok, err := authenticator.Auth(validToken)
		assert.NoError(t, err)
		assert.True(t, ok)
	})
}

func TestAuthenticator_Token(t *testing.T) {
	key, err := getPrivateKey(ThirdParty)
	if err != nil {
		t.Fatal(err)
	}
	pk, err := getPublicKey(ThirdParty)
	if err != nil {
		t.Fatal(err)
	}
	authenticator := Authenticator{
		PrivateKeys: map[string]*rsa.PrivateKey{
			ThirdParty: key,
		},
		ModelHandler: MockModelHandler{},
	}
	t.Run("testing getting token with valid inputs", func(t *testing.T) {
		tokenString, err := authenticator.Token(ClientCredentials, "snappbox", "KJIikjIKbIYVGj)YihYUGIB&")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return pk, nil
		})
		assert.NoError(t, err)
		claims := token.Claims.(jwt.MapClaims)
		assert.Equal(t, "snappbox", claims["sub"].(string))
		assert.Equal(t, "100", claims["iss"].(string))
	})
	t.Run("testing getting token with valid inputs", func(t *testing.T) {
		tokenString, err := authenticator.Token(ClientCredentials, "snappbox", "invalid secret")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return pk, nil
		})
		assert.Error(t, err)
		assert.Nil(t, token)
	})
}

func TestAuthenticator_Acl(t *testing.T) {
	key, err := getPrivateKey(ThirdParty)
	if err != nil {
		t.Fatal(err)
	}
	tokenString, err := getSampleToken(true)
	if err != nil {
		t.Fatal(err)
	}
	invalidTokenString, err := getSampleToken(false)
	if err != nil {
		t.Fatal(t, err)
	}
	authenticator := Authenticator{
		PrivateKeys: map[string]*rsa.PrivateKey{
			ThirdParty: key,
		},
		AllowedAccessTypes: []string{Pub, Sub},
		ModelHandler:       MockModelHandler{},
		SnappIDsManager: snappids.NewEMQManager(map[snappids.Audience]string{
			snappids.DriverAudience: "12345678901234567890123456789012",
		}, map[snappids.Audience]int{
			snappids.DriverAudience: 15,
		}),
	}
	t.Run("testing acl with invalid access type", func(t *testing.T) {
		ok, err := authenticator.Acl(PubSub, tokenString, "test")
		assert.Error(t, err)
		assert.False(t, ok)
		assert.Equal(t, "requested access type 3 is invalid", err.Error())
	})
	t.Run("testing acl with invalid token", func(t *testing.T) {
		ok, err := authenticator.Acl(Pub, invalidTokenString, "driver-event-5ab8f6e552c445d0c8d38f9f38ca4e2b")
		assert.False(t, ok)
		assert.Error(t, err)
		assert.Equal(t, "illegal base64 data at input byte 37", err.Error())
	})
	t.Run("testing acl with valid inputs", func(t *testing.T) {
		ok, err := authenticator.Acl(Sub, tokenString, "driver-event-5ab8f6e552c445d0c8d38f9f38ca4e2b")
		assert.NoError(t, err)
		assert.True(t, ok)
	})
	t.Run("testing acl with invalid topic", func(t *testing.T) {
		ok, err := authenticator.Acl(Sub, tokenString, "driver-event-5ab8f6e552c4423fc8d38f9f38ca4e2b")
		assert.Error(t, err)
		assert.False(t, ok)
	})
	t.Run("testing acl with invalid access type", func(t *testing.T) {
		ok, err := authenticator.Acl(Pub, tokenString, "driver-event-5ab8f6e552c445d0c8d38f9f38ca4e2b")
		assert.Error(t, err)
		assert.False(t, ok)
	})

}

func TestAuthenticator_validateAccessType(t *testing.T) {
	type fields struct {
		AllowedAccessTypes []string
	}
	type args struct {
		accessType string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "#1 testing with no allowed access type",
			fields: fields{AllowedAccessTypes: []string{}},
			args:   args{accessType: Sub},
			want:   false,
		},
		{
			name:   "#2 testing with no allowed access type",
			fields: fields{AllowedAccessTypes: []string{}},
			args:   args{accessType: Pub},
			want:   false,
		},
		{
			name:   "#3 testing with no allowed access type",
			fields: fields{AllowedAccessTypes: []string{}},
			args:   args{accessType: PubSub},
			want:   false,
		},
		{
			name:   "#4 testing with one allowed access type",
			fields: fields{AllowedAccessTypes: []string{Pub}},
			args:   args{accessType: Pub},
			want:   true,
		},
		{
			name:   "#5 testing with one allowed access type",
			fields: fields{AllowedAccessTypes: []string{Pub}},
			args:   args{accessType: Sub},
			want:   false,
		},
		{
			name:   "#6 testing with two allowed access type",
			fields: fields{AllowedAccessTypes: []string{Pub, Sub}},
			args:   args{accessType: Sub},
			want:   true,
		},
		{
			name:   "#7 testing with two allowed access type",
			fields: fields{AllowedAccessTypes: []string{Pub, Sub}},
			args:   args{accessType: Pub},
			want:   true,
		},
		{
			name:   "#8 testing with two allowed access type",
			fields: fields{AllowedAccessTypes: []string{Pub, Sub}},
			args:   args{accessType: PubSub},
			want:   false,
		},
		{
			name:   "#9 testing with three allowed access type",
			fields: fields{AllowedAccessTypes: []string{Pub, Sub, PubSub}},
			args:   args{accessType: PubSub},
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

func (rmh MockModelHandler) Save(model db.Model) error {
	return nil
}

func (rmh MockModelHandler) Delete(modelName, pk string) error {
	return nil
}

func (rmh MockModelHandler) Get(modelName, pk string, v interface{}) error {
	key0, _ := getPublicKey(Driver)
	key1, _ := getPublicKey(Passenger)
	key100, _ := getPublicKey(ThirdParty)
	switch pk {
	case Passenger:
		*v.(*User) = User{
			MetaData:  db.MetaData{},
			Username:  Passenger,
			Type:      EMQUser,
			PublicKey: key1,
		}
	case Driver:
		*v.(*User) = User{
			MetaData:  db.MetaData{},
			Username:  Driver,
			Type:      EMQUser,
			PublicKey: key0,
			Rules: []Rule{{
				UID:          1,
				Endpoint:     "",
				TopicPattern: `(\w+)-event-(\w*\d*|\d*\w*)`,
				AccessType:   Sub,
			}},
		}
	case "snappbox":
		*v.(*User) = User{
			MetaData:                db.MetaData{},
			Username:                "snappbox",
			Password:                getSamplePassword(),
			Type:                    HeraldUser,
			PublicKey:               key100,
			Secret:                  "KJIikjIKbIYVGj)YihYUGIB&",
			TokenExpirationDuration: time.Hour * 72,
		}
	}
	return nil
}

func getPublicKey(user string) (*rsa.PublicKey, error) {
	var fileName string
	switch user {
	case Passenger:
		fileName = "../../test/1.test.pem"
	case Driver:
		fileName = "../../test/1.test.pem"
	case ThirdParty:
		fileName = "../../test/100.test.pem"
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

func getPrivateKey(user string) (*rsa.PrivateKey, error) {
	var fileName string
	switch user {
	case ThirdParty:
		fileName = "../../test/100.test.private.pem"
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

func getSampleToken(isValid bool) (string, error) {
	var fileName string
	if isValid {
		fileName = "../../test/token.valid.sample"
	} else {
		fileName = "../../test/token.invalid.sample"
	}
	token, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(token), nil
}

func getSamplePassword() []byte {
	hash, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	return hash
}
