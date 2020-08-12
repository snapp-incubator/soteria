package accounts

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gitlab.snapp.ir/dispatching/snappids"
	"gitlab.snapp.ir/dispatching/soteria/internal/db"
	"gitlab.snapp.ir/dispatching/soteria/pkg/errors"
	"gitlab.snapp.ir/dispatching/soteria/pkg/acl"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"
)

// Authenticator is responsible for Acl/Auth/Token of users
type Authenticator struct {
	PrivateKeys        map[string]*rsa.PrivateKey
	AllowedAccessTypes []string
	ModelHandler       db.ModelHandler
	SnappIDsManager    *snappids.Manager
}

// TopicMan is a function that takes issuer and subject as inputs and generates topic name
type TopicMan func(issuer, sub string) string

const (
	// Issuers
	Driver     = "0"
	Passenger  = "1"
	ThirdParty = "100"

	// Access Types
	Sub    = "1"
	Pub    = "2"
	PubSub = "3"

	ClientCredentials = "client_credentials"
)

// Auth check user authentication by checking the user's token
func (a Authenticator) Auth(tokenString string) (bool, error) {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("token is not valid, signing method is not RSA")
		}
		claims := token.Claims.(jwt.MapClaims)
		if claims["iss"] == nil {
			return nil, fmt.Errorf("could not found iss in token claims")
		}
		issuer := fmt.Sprintf("%v", claims["iss"])
		user := User{}
		err = a.ModelHandler.Get("user", issuer, &user)
		if err != nil {
			return false, fmt.Errorf("error getting issuer from db err: %v", err)
		}
		key := user.PublicKey
		if key == nil {
			return nil, fmt.Errorf("cannot find issuer %v public key", issuer)
		}
		return key, nil
	})
	if err != nil {
		return false, fmt.Errorf("token is invalid err: %v", err)
	}
	return true, nil
}

// ACL check a user access to a topic
func (a Authenticator) Acl(accessType, tokenString, topic string) (bool, error) {
	if !a.validateAccessType(accessType) {
		return false, fmt.Errorf("requested access type %s is invalid", accessType)
	}
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("token is not valid, signing method is not RSA")
		}
		claims := token.Claims.(jwt.MapClaims)
		if claims["iss"] == nil {
			return nil, fmt.Errorf("could not found iss in token claims")
		}
		if claims["sub"] == nil {
			return nil, fmt.Errorf("could not find sub in token claims")
		}
		issuer := fmt.Sprintf("%v", claims["iss"])
		sub := fmt.Sprintf("%v", claims["sub"])
		user := User{}
		err := a.ModelHandler.Get("user", primaryKey(issuer, sub), &user)
		if err != nil {
			return false, fmt.Errorf("error getting user from db err: %v", err)
		}
		key := user.PublicKey
		if key == nil {
			return nil, fmt.Errorf("cannot find user %v public key", issuer)
		}
		topicMan := a.getTopicMan(topic)
		if ok := user.CheckTopicAllowance(topicMan, issuer, sub, topic, accessType); !ok {
			return nil, fmt.Errorf("issuer %v with sub %v is not allowed to %v on topic %v", issuer, sub, accessType, topic)
		}
		return key, nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

// Token function issues JWT token by taking client credentials
func (a Authenticator) Token(accessType, username, secret string) (tokenString string, err error) {
	if accessType == ClientCredentials {
		accessType = Sub
	}
	user := User{}
	err = a.ModelHandler.Get("user", username, &user)
	if err != nil {
		return "", fmt.Errorf("could not get user. err: %v", err)
	}
	if user.Secret != secret {
		return "", fmt.Errorf("invlaid secret %v", secret)
	}

	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(user.TokenExpirationDuration).Unix(),
		Issuer:    ThirdParty,
		Subject:   username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err = token.SignedString(a.PrivateKeys[ThirdParty])
	if err != nil {
		return "", fmt.Errorf("could not sign the token. err; %v", err)
	}
	return tokenString, nil
}

func (a Authenticator) EndPointBasicAuth(username, password, endpoint string) (bool, error) {
	var user User
	if err := ModelHandler.Get("user", username, &user); err != nil {
		return false, errors.CreateError(errors.DatabaseGetFailure, err.Error())
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		return false, errors.CreateError(errors.WrongUsernameOrPassword, "wrong password")
	}
	ok := user.CheckEndpointAllowance(endpoint)
	if !ok {
		return false, nil
	}
	return true, nil
}

func (a Authenticator) EndpointIPAuth(username string, ip string, endpoint string) (bool, error) {
	var user User
	if err := ModelHandler.Get("user", username, &user); err != nil {
		return false, errors.CreateError(errors.DatabaseGetFailure, err.Error())
	}
	ok := acl.ValidateIP(ip, user.IPs, []string{})
	if !ok {
		return false, errors.CreateError(errors.IPMisMatch, "ip mismatch")
	}
	ok = user.CheckEndpointAllowance(endpoint)
	if !ok {
		return false, nil
	}
	return true, nil
}

func (a Authenticator) validateAccessType(accessType string) bool {
	for _, allowedAccessType := range a.AllowedAccessTypes {
		if allowedAccessType == accessType {
			return true
		}
	}
	return false
}

func primaryKey(issuer, sub string) string {
	if issuer == Passenger || issuer == Driver {
		return issuer
	}
	return sub
}

func (a Authenticator) getTopicMan(topic string) TopicMan {
	matched, _ := regexp.Match(`(\w+)-event-(\w*\d*|\d*\w*)`, []byte(topic))
	if matched {
		return func(issuer, sub string) string {
			id, _ := a.SnappIDsManager.DecodeHashID(sub, toAudience(issuer))
			ch, _ := a.SnappIDsManager.CreateChannelName(id, toAudience(issuer))
			return ch
		}
	}
	matched, _ = regexp.Match(`snapp/driver/(\w*\d*|\d*\w*)/location`, []byte(topic))
	if matched {
		return func(issuer, sub string) string {
			ch, _ := a.SnappIDsManager.CreateDriverLocationChannelName(sub)
			return ch
		}
	}
	return nil
}

func toAudience(issuer string) snappids.Audience {
	switch issuer {
	case Passenger:
		return snappids.PassengerAudience
	case Driver:
		return snappids.DriverAudience
	case ThirdParty:
		return snappids.ThirdPartyAudience
	default:
		return -1
	}
}
