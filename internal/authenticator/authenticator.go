package authenticator

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	snappids "gitlab.snapp.ir/dispatching/snappids/v2"
	"gitlab.snapp.ir/dispatching/soteria/internal/db"
	"gitlab.snapp.ir/dispatching/soteria/pkg/acl"
	"gitlab.snapp.ir/dispatching/soteria/pkg/errors"
	"gitlab.snapp.ir/dispatching/soteria/pkg/user"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"
)

// Authenticator is responsible for Acl/Auth/Token of users
type Authenticator struct {
	PrivateKeys        map[string]*rsa.PrivateKey
	AllowedAccessTypes []user.AccessType
	ModelHandler       db.ModelHandler
	EMQTopicManager    *snappids.EMQTopicManager
	HashIDSManager     *snappids.HashIDSManager
}

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
		u := user.User{}
		err = a.ModelHandler.Get("user", issuer, &u)
		if err != nil {
			return false, fmt.Errorf("error getting issuer from db err: %v", err)
		}
		key := u.PublicKey
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
func (a Authenticator) Acl(accessType user.AccessType, tokenString, topic string) (bool, error) {
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
		u := user.User{}
		err := a.ModelHandler.Get("user", primaryKey(user.Issuer(issuer), sub), &u)
		if err != nil {
			return false, fmt.Errorf("error getting user from db err: %v", err)
		}
		key := u.PublicKey
		if key == nil {
			return nil, fmt.Errorf("cannot find user %v public key", issuer)
		}
		topicMan := a.getTopicMan(topic)
		if ok := u.CheckTopicAllowance(topicMan, user.Issuer(issuer), sub, topic, accessType); !ok {
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
func (a Authenticator) Token(accessType user.AccessType, username, secret string) (tokenString string, err error) {
	if accessType == user.ClientCredentials {
		accessType = user.Sub
	}
	u := user.User{}
	err = a.ModelHandler.Get("user", username, &u)
	if err != nil {
		return "", fmt.Errorf("could not get user. err: %v", err)
	}
	if u.Secret != secret {
		return "", fmt.Errorf("invlaid secret %v", secret)
	}

	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(u.TokenExpirationDuration).Unix(),
		Issuer:    user.ThirdParty,
		Subject:   username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err = token.SignedString(a.PrivateKeys[user.ThirdParty])
	if err != nil {
		return "", fmt.Errorf("could not sign the token. err; %v", err)
	}
	return tokenString, nil
}

func (a Authenticator) EndPointBasicAuth(username, password, endpoint string) (bool, error) {
	var u user.User
	if err := a.ModelHandler.Get("user", username, &u); err != nil {
		return false, errors.CreateError(errors.DatabaseGetFailure, err.Error())
	}

	if err := bcrypt.CompareHashAndPassword(u.Password, []byte(password)); err != nil {
		return false, errors.CreateError(errors.WrongUsernameOrPassword, "wrong password")
	}
	ok := u.CheckEndpointAllowance(endpoint)
	if !ok {
		return false, nil
	}
	return true, nil
}

func (a Authenticator) EndpointIPAuth(username string, ip string, endpoint string) (bool, error) {
	var u user.User
	if err := a.ModelHandler.Get("user", username, &u); err != nil {
		return false, errors.CreateError(errors.DatabaseGetFailure, err.Error())
	}
	ok := acl.ValidateIP(ip, u.IPs, []string{})
	if !ok {
		return false, errors.CreateError(errors.IPMisMatch, "ip mismatch")
	}
	ok = u.CheckEndpointAllowance(endpoint)
	if !ok {
		return false, nil
	}
	return true, nil
}

func (a Authenticator) validateAccessType(accessType user.AccessType) bool {
	for _, allowedAccessType := range a.AllowedAccessTypes {
		if allowedAccessType == accessType {
			return true
		}
	}
	return false
}

func primaryKey(issuer user.Issuer, sub string) string {
	if issuer == user.Passenger || issuer == user.Driver {
		return string(issuer)
	}
	return sub
}

func (a Authenticator) getTopicMan(topic string) user.TopicMan {
	matched, _ := regexp.Match(`(\w+)-event-[a-zA-Z0-9]+`, []byte(topic))
	if matched {
		return func(issuer user.Issuer, sub string) string {
			id, _ := a.HashIDSManager.DecodeHashID(sub, toAudience(issuer))
			ch, _ := a.EMQTopicManager.CreateCabEventTopic(id, toAudience(issuer))
			return string(ch)
		}
	}
	matched, _ = regexp.Match(`snapp/driver/[a-zA-Z0-9]+/location`, []byte(topic))
	if matched {
		return func(issuer user.Issuer, sub string) string {
			id, _ := a.HashIDSManager.DecodeHashID(sub, toAudience(issuer))
			ch, _ := a.EMQTopicManager.CreateLocationTopic(id, toAudience(issuer))
			return string(ch)
		}
	}
	return nil
}

func toAudience(issuer user.Issuer) snappids.Audience {
	switch issuer {
	case user.Passenger:
		return snappids.PassengerAudience
	case user.Driver:
		return snappids.DriverAudience
	case user.ThirdParty:
		return snappids.ThirdPartyAudience
	default:
		return -1
	}
}
