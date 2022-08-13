package authenticator

import (
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"gitlab.snapp.ir/dispatching/soteria/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/pkg/acl"
)

const DefaultVendor = "snapp"

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
	Issuer     string
	Sub        string
	AccessType acl.AccessType
	Topic      string
	TopicType  string
}

func (err TopicNotAllowedError) Error() string {
	return fmt.Sprintf("issuer %s with sub %s is not allowed to %s on topic %s (%s)",
		err.Issuer, err.Sub, err.AccessType, err.Topic, err.TopicType,
	)
}

type PublicKeyNotFoundError struct {
	Issuer string
}

func (err PublicKeyNotFoundError) Error() string {
	return fmt.Sprintf("cannot find issuer %s public key", err.Issuer)
}

type InvalidTopicError struct {
	Topic string
}

func (err InvalidTopicError) Error() string {
	return fmt.Sprintf("provided topic %s is not valid", err.Topic)
}

// Authenticator is responsible for Acl/Auth/Token of users.
type Authenticator struct {
	PublicKeys         map[string]*rsa.PublicKey
	AllowedAccessTypes []acl.AccessType
	TopicManager       *topics.Manager
	Company            string
}

// Auth check user authentication by checking the user's token
// isSuperuser is a flag that authenticator set it true when credentials is related to a superuser.
func (a Authenticator) Auth(tokenString string) error {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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

		issuer := fmt.Sprintf("%v", claims["iss"])

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
// nolint: funlen, cyclop
func (a Authenticator) ACL(
	accessType acl.AccessType,
	tokenString string,
	topic string,
) (bool, error) {
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

		issuer := fmt.Sprintf("%v", claims["iss"])
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

	issuer := fmt.Sprintf("%v", claims["iss"])

	if claims["sub"] == nil {
		return false, ErrSubNotFound
	}

	sub, _ := claims["sub"].(string)

	topicTemplate := a.TopicManager.ValidateTopic(topic, issuer, sub)
	if topicTemplate == nil {
		return false, InvalidTopicError{Topic: topic}
	}

	if !topicTemplate.HasAccess(issuer, accessType) {
		return false, TopicNotAllowedError{
			issuer,
			sub,
			accessType,
			topic,
			topicTemplate.Type,
		}
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
