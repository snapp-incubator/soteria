package authenticator

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gitlab.snapp.ir/dispatching/snappids/v2"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/db"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/acl"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/user"
)

var (
	ErrInvalidSigningMethod = errors.New("token is not valid, signing method is not RSA")
	ErrIssNotFound          = errors.New("could not found iss in token claims")
	ErrSubNotFound          = errors.New("could not found sub in token claims")
	ErrInvalidClaims        = errors.New("invalid claims")
	ErrInvalidIP            = errors.New("IP is not valid")
	ErrInvalidAccessType    = errors.New("requested access type is invalid")
	ErrDecodeHashID         = errors.New("could not decode hash id")
	ErrInvalidSecret        = errors.New("invalid secret")
	ErrIncorrectPassword    = errors.New("username or password is worng")
)

type ErrTopicNotAllowed struct {
	Issuer     user.Issuer
	Sub        string
	AccessType acl.AccessType
	Topic      topics.Topic
}

func (err ErrTopicNotAllowed) Error() string {
	return fmt.Sprintf("issuer %s with sub %s is not allowed to %s on topic %s (%s)",
		err.Issuer, err.Sub, err.AccessType, err.Topic, err.Topic.GetType(),
	)
}

type ErrPublicKeyNotFound struct {
	Issuer user.Issuer
}

func (err ErrPublicKeyNotFound) Error() string {
	return fmt.Sprintf("cannot find issuer %s public key", err.Issuer)
}

type ErrInvalidTopic struct {
	Topic topics.Topic
}

func (err ErrInvalidTopic) Error() string {
	return fmt.Sprintf("provided topic %s is not valid", err.Topic)
}

// Authenticator is responsible for Acl/Auth/Token of users.
type Authenticator struct {
	PrivateKeys            map[user.Issuer]*rsa.PrivateKey
	PublicKeys             map[user.Issuer]*rsa.PublicKey
	AllowedAccessTypes     []acl.AccessType
	ModelHandler           db.ModelHandler
	EMQTopicManager        *snappids.EMQTopicManager
	HashIDSManager         *snappids.HashIDSManager
	CompareHashAndPassword func([]byte, []byte) error
	Company                string
}

// Auth check user authentication by checking the user's token
// isSuperuser is a flag that authenticator set it true when credentials is related to a superuser.
func (a Authenticator) Auth(ctx context.Context, tokenString string) (isSuperuser bool, err error) {
	_, err = jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, ErrInvalidSigningMethod
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, ErrInvalidClaims
		}
		if claims["iss"] == nil {
			return nil, ErrIssNotFound
		}

		_, isSuperuser = claims["is_superuser"]

		issuer := user.Issuer(fmt.Sprintf("%v", claims["iss"]))

		key := a.PublicKeys[issuer]
		if key == nil {
			return nil, ErrPublicKeyNotFound{Issuer: issuer}
		}

		return key, nil
	})

	if err != nil {
		return false, fmt.Errorf("token is invalid: %w", err)
	}

	return isSuperuser, nil
}

// ACL check a user access to a topic.
func (a Authenticator) ACL(ctx context.Context, accessType acl.AccessType,
	tokenString string, topic topics.Topic) (bool, error) {
	if !a.validateAccessType(accessType) {
		return false, ErrInvalidAccessType
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, ErrInvalidSigningMethod
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, ErrInvalidClaims
		}
		if claims["iss"] == nil {
			return nil, ErrIssNotFound
		}
		if claims["sub"] == nil {
			return nil, ErrSubNotFound
		}

		issuer := user.Issuer(fmt.Sprintf("%v", claims["iss"]))
		key := a.PublicKeys[issuer]
		if key == nil {
			return nil, ErrPublicKeyNotFound{Issuer: issuer}
		}

		return key, nil
	})
	if err != nil {
		return false, fmt.Errorf("token is invalid %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, ErrInvalidClaims
	}

	if claims["iss"] == nil {
		return false, ErrIssNotFound
	}

	issuer := user.Issuer(fmt.Sprintf("%v", claims["iss"]))

	if claims["sub"] == nil {
		return false, ErrSubNotFound
	}

	sub := fmt.Sprintf("%v", claims["sub"])

	pk := primaryKey(issuer, sub)

	var u user.User
	if err := a.ModelHandler.Get(ctx, "user", pk, &u); err != nil {
		return false, fmt.Errorf("error getting user %s from db err: %w", pk, err)
	}

	// validate passenger and driver topics.
	if issuer != user.ThirdParty {
		id, err := a.HashIDSManager.DecodeHashID(sub, issuerToAudience(issuer))
		if err != nil {
			return false, ErrDecodeHashID
		}

		ok := a.ValidateTopicBySender(topic, issuerToAudience(issuer), id)
		if !ok {
			return false, ErrInvalidTopic{Topic: topic}
		}
	}

	if ok := u.CheckTopicAllowance(topic.GetTypeWithCompany(a.Company), accessType); !ok {
		return false,
			ErrTopicNotAllowed{issuer, sub, accessType, topic}
	}

	return true, nil
}

// Token function issues JWT token by taking client credentials.
func (a Authenticator) Token(ctx context.Context, _ acl.AccessType,
	username, secret string) (string, error) {
	var u user.User
	if err := a.ModelHandler.Get(ctx, "user", username, &u); err != nil {
		return "", fmt.Errorf("could not get user %s. %w", username, err)
	}

	if u.Secret != secret {
		return "", ErrInvalidSecret
	}

	// nolint: exhaustivestruct
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(u.TokenExpirationDuration).Unix(),
		Issuer:    string(user.ThirdParty),
		Subject:   username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(a.PrivateKeys[user.ThirdParty])
	if err != nil {
		return "", fmt.Errorf("could not sign the token %w", err)
	}

	return tokenString, nil
}

func (a Authenticator) HeraldToken(
	username string,
	endpoints []acl.Endpoint,
	topics []acl.Topic,
	duration time.Duration,
) (string, error) {
	// nolint: exhaustivestruct
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

	tokenString, err := token.SignedString(a.PrivateKeys[user.ThirdParty])
	if err != nil {
		return "", fmt.Errorf("could not sign the token. %w", err)
	}

	return tokenString, nil
}

func (a Authenticator) SuperuserToken(username string, duration time.Duration) (string, error) {
	// nolint: exhaustivestruct
	claims := acl.SuperuserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			Issuer:    string(user.ThirdParty),
			Subject:   username,
		},
		IsSuperuser: true,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(a.PrivateKeys[user.ThirdParty])
	if err != nil {
		return "", fmt.Errorf("could not sign the token. %w", err)
	}

	return tokenString, nil
}

func (a Authenticator) EndPointBasicAuth(ctx context.Context, username, password, endpoint string) (bool, error) {
	var u user.User
	if err := a.ModelHandler.Get(ctx, "user", username, &u); err != nil {
		return false, fmt.Errorf("could not get user from db: %w", err)
	}

	if err := a.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return false, ErrIncorrectPassword
	}

	return u.CheckEndpointAllowance(endpoint), nil
}

func (a Authenticator) EndpointIPAuth(ctx context.Context, username string, ip string, endpoint string) (bool, error) {
	var u user.User
	if err := a.ModelHandler.Get(ctx, "user", username, &u); err != nil {
		return false, fmt.Errorf("could not get user from db: %w", err)
	}

	if ok := acl.ValidateIP(ip, u.IPs, []string{}); !ok {
		return false, ErrInvalidIP
	}

	return u.CheckEndpointAllowance(endpoint), nil
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
	} else if issuer == user.Driver {
		return "driver"
	}

	return sub
}

func (a Authenticator) ValidateTopicBySender(topic topics.Topic, audience snappids.Audience, id int) bool {
	var ch snappids.Topic

	switch topic.GetType() {
	case topics.CabEvent:
		ch, _ = a.EMQTopicManager.CreateCabEventTopic(id, audience)
	case topics.DriverLocation, topics.PassengerLocation:
		ch, _ = a.EMQTopicManager.CreateLocationTopic(id, audience)
	case topics.SuperappEvent:
		ch, _ = a.EMQTopicManager.CreateSuperAppEventTopic(id, audience)
	case topics.SharedLocation:
		ch, _ = a.EMQTopicManager.CreateSharedLocationTopic(id, audience)
	case topics.Chat:
		ch, _ = a.EMQTopicManager.CreateChatTopic(id, audience)
	case topics.Call:
		ch, _ = a.EMQTopicManager.CreateCallTopic(id, audience)
	case topics.DaghighSys:
		return true
	case topics.BoxEvent:
		return true
	}

	return string(ch) == string(topic)
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
