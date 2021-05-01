package authenticator

import (
	"context"
	"crypto/rsa"
	errs "errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	snappids "gitlab.snapp.ir/dispatching/snappids/v2"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/db"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/acl"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/user"
)

var TopicNotAllowed = errs.New("topic is not allowed")

// Authenticator is responsible for Acl/Auth/Token of users
type Authenticator struct {
	PrivateKeys            map[user.Issuer]*rsa.PrivateKey
	PublicKeys             map[user.Issuer]*rsa.PublicKey
	AllowedAccessTypes     []acl.AccessType
	ModelHandler           db.ModelHandler
	EMQTopicManager        *snappids.EMQTopicManager
	HashIDSManager         *snappids.HashIDSManager
	CompareHashAndPassword func([]byte, []byte) error
}

// Auth check user authentication by checking the user's token
// isSuperuser is a flag that authenticator set it true when credentials is related to a superuser.
func (a Authenticator) Auth(ctx context.Context, tokenString string) (isSuperuser bool, err error) {
	_, err = jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("token is not valid, signing method is not RSA")
		}
		claims := token.Claims.(jwt.MapClaims)
		if claims["iss"] == nil {
			return nil, fmt.Errorf("could not found iss in token claims")
		}
		_, isSuperuser = claims["is_superuser"]
		issuer := user.Issuer(fmt.Sprintf("%v", claims["iss"]))
		key := a.PublicKeys[issuer]
		if key == nil {
			return nil, fmt.Errorf("cannot find issuer %s public key", issuer)
		}
		return key, nil
	})
	if err != nil {
		return false, fmt.Errorf("token is invalid err: %w", err)
	}
	return isSuperuser, nil
}

// ACL check a user access to a topic
func (a Authenticator) Acl(ctx context.Context, accessType acl.AccessType, tokenString string, topic topics.Topic) (bool, error) {
	if !a.validateAccessType(accessType) {
		return false, fmt.Errorf("requested access type %s is invalid", accessType)
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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

		issuer := user.Issuer(fmt.Sprintf("%v", claims["iss"]))
		key := a.PublicKeys[issuer]
		if key == nil {
			return nil, fmt.Errorf("cannot find user %v public key", issuer)
		}
		return key, nil
	})
	if err != nil {
		return false, fmt.Errorf("token is invalid. err: %w", err)
	}

	claims := token.Claims.(jwt.MapClaims)
	issuer := user.Issuer(fmt.Sprintf("%v", claims["iss"]))
	sub := fmt.Sprintf("%v", claims["sub"])
	pk := primaryKey(issuer, sub)
	u := user.User{}
	err = a.ModelHandler.Get(ctx, "user", pk, &u)
	if err != nil {
		return false, fmt.Errorf("error getting user %s from db err: %w", pk, err)
	}
	if issuer != user.ThirdParty {
		id, err := a.HashIDSManager.DecodeHashID(sub, issuerToAudience(issuer))
		if err != nil {
			return false, fmt.Errorf("could not decode hash id")
		}
		ok := a.ValidateTopicBySender(topic, issuerToAudience(issuer), id)
		if !ok {
			return false, fmt.Errorf("provided topic %v is not valid", topic)
		}
	}

	if ok := u.CheckTopicAllowance(topic.GetType(), accessType); !ok {
		return false,
			fmt.Errorf("%w. issuer %s with sub %s is not allowed to %s on topic %s", TopicNotAllowed, issuer, sub, accessType, topic)
	}
	return true, nil
}

// Token function issues JWT token by taking client credentials
func (a Authenticator) Token(ctx context.Context, accessType acl.AccessType, username, secret string) (tokenString string, err error) {
	if accessType == acl.ClientCredentials {
		accessType = acl.Sub
	}
	u := user.User{}
	err = a.ModelHandler.Get(ctx, "user", username, &u)
	if err != nil {
		return "", fmt.Errorf("could not get user %s. err: %w", username, err)
	}
	if u.Secret != secret {
		return "", fmt.Errorf("invlaid secret %v", secret)
	}
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(u.TokenExpirationDuration).Unix(),
		Issuer:    string(user.ThirdParty),
		Subject:   username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err = token.SignedString(a.PrivateKeys[user.ThirdParty])
	if err != nil {
		return "", fmt.Errorf("could not sign the token. err; %v", err)
	}
	return tokenString, nil
}

func (a Authenticator) HeraldToken(
	username string,
	endpoints []acl.Endpoint,
	topics []acl.Topic,
	duration time.Duration) (tokenString string, err error) {

	claims := acl.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			Issuer:    string(user.ThirdParty),
			Subject:   username,
		},
		Topics:    topics,
		Endpoints: endpoints,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err = token.SignedString(a.PrivateKeys[user.ThirdParty])
	if err != nil {
		return "", fmt.Errorf("could not sign the token. err; %v", err)
	}
	return tokenString, nil
}

func (a Authenticator) SuperuserToken(username string, duration time.Duration) (tokenString string, err error) {
	claims := acl.SuperuserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			Issuer:    string(user.ThirdParty),
			Subject:   username,
		},
		IsSuperuser: true,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err = token.SignedString(a.PrivateKeys[user.ThirdParty])
	if err != nil {
		return "", fmt.Errorf("could not sign the token. err; %v", err)
	}
	return tokenString, nil
}

func (a Authenticator) EndPointBasicAuth(ctx context.Context, username, password, endpoint string) (bool, error) {
	var u user.User
	if err := a.ModelHandler.Get(ctx, "user", username, &u); err != nil {
		return false, fmt.Errorf("could not get user from db: %w", err)
	}

	if err := a.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return false, fmt.Errorf("username or password is worng")
	}
	ok := u.CheckEndpointAllowance(endpoint)
	if !ok {
		return false, nil
	}
	return true, nil
}

func (a Authenticator) EndpointIPAuth(ctx context.Context, username string, ip string, endpoint string) (bool, error) {
	var u user.User
	if err := a.ModelHandler.Get(ctx, "user", username, &u); err != nil {
		return false, fmt.Errorf("could not get user from db: %w", err)
	}
	ok := acl.ValidateIP(ip, u.IPs, []string{})
	if !ok {
		return false, fmt.Errorf("IP is not valid")
	}
	ok = u.CheckEndpointAllowance(endpoint)
	if !ok {
		return false, nil
	}
	return true, nil
}

func (a Authenticator) validateAccessType(accessType acl.AccessType) bool {
	for _, allowedAccessType := range a.AllowedAccessTypes {
		if allowedAccessType == accessType {
			return true
		}
	}
	return false
}

func primaryKey(issuer user.Issuer, sub string) string {
	if issuer == user.Passenger {
		return "passenger"
	}
	if issuer == user.Driver {
		return "driver"
	}
	return sub
}

func (a Authenticator) ValidateTopicBySender(topic topics.Topic, audience snappids.Audience, id int) bool {
	var ch snappids.Topic
	switch topic.GetType() {
	case topics.CabEvent:
		ch, _ = a.EMQTopicManager.CreateCabEventTopic(id, audience)
	case topics.DriverLocation:
		ch, _ = a.EMQTopicManager.CreateLocationTopic(id, audience)
	case topics.SuperappEvent:
		ch, _ = a.EMQTopicManager.CreateSuperAppEventTopic(id, audience)
	case topics.GossiperLocation:
		ch, _ = a.EMQTopicManager.CreateGossiperTopic(id, audience)
	case topics.BoxEvent:
		return true
	}
	if string(ch) != string(topic) {
		return false
	}
	return true
}

func issuerToAudience(issuer user.Issuer) snappids.Audience {
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
