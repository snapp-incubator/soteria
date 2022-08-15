package authenticator

import (
	"errors"
	"fmt"
	"gitlab.snapp.ir/dispatching/soteria/pkg/acl"
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
