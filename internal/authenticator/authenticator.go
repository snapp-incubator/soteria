package authenticator

import (
	"crypto/rsa"
	"errors"
	"fmt"

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

type TopicNotAllowedError struct {
	Issuer     user.Issuer
	Sub        string
	AccessType acl.AccessType
	Topic      topics.Topic
}

func (err TopicNotAllowedError) Error() string {
	return fmt.Sprintf("issuer %s with sub %s is not allowed to %s on topic %s (%s)",
		err.Issuer, err.Sub, err.AccessType, err.Topic, err.Topic.GetType(),
	)
}

type PublicKeyNotFoundError struct {
	Issuer user.Issuer
}

func (err PublicKeyNotFoundError) Error() string {
	return fmt.Sprintf("cannot find issuer %s public key", err.Issuer)
}

type InvalidTopicError struct {
	Topic topics.Topic
}

func (err InvalidTopicError) Error() string {
	return fmt.Sprintf("provided topic %s is not valid", err.Topic)
}

// Authenticator is responsible for Acl/Auth/Token of users.
type Authenticator struct {
	PublicKeys         map[user.Issuer]*rsa.PublicKey
	AllowedAccessTypes []acl.AccessType
	ModelHandler       db.ModelHandler
	EMQTopicManager    *snappids.EMQTopicManager
	HashIDSManager     *snappids.HashIDSManager
	Company            string
}

// Auth check user authentication by checking the user's token
// isSuperuser is a flag that authenticator set it true when credentials is related to a superuser.
func (a Authenticator) Auth(tokenString string) (err error) {
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

		issuer := user.Issuer(fmt.Sprintf("%v", claims["iss"]))

		key := a.PublicKeys[issuer]
		if key == nil {
			return nil, PublicKeyNotFoundError{Issuer: issuer}
		}

		return key, nil
	})
	if err != nil {
		return fmt.Errorf("token is invalid: %w", err)
	}

	return nil
}

// ACL check a user access to a topic.
func (a Authenticator) ACL(accessType acl.AccessType,
	tokenString string, topic topics.Topic) (bool, error) {
	if !a.ValidateAccessType(accessType) {
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
			return nil, PublicKeyNotFoundError{Issuer: issuer}
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

	user := a.ModelHandler.Get(pk)

	id, err := a.HashIDSManager.DecodeHashID(sub, issuerToAudience(issuer))
	if err != nil {
		return false, ErrDecodeHashID
	}

	if !a.ValidateTopicBySender(topic, issuerToAudience(issuer), id) {
		return false, InvalidTopicError{Topic: topic}
	}

	if ok := user.CheckTopicAllowance(topic.GetTypeWithCompany(a.Company), accessType); !ok {
		return false,
			TopicNotAllowedError{issuer, sub, accessType, topic}
	}

	return true, nil
}

func (a Authenticator) ValidateAccessType(accessType acl.AccessType) bool {
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
	case topics.CallEntry:
		ch, _ = a.EMQTopicManager.CreateCallEntryTopic(id, audience)
	case topics.CallOutgoing:
		ch, _ = a.EMQTopicManager.CreateCallOutgoingTopic(id, audience)
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
	default:
		return -1
	}
}
