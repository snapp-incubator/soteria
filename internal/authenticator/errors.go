package authenticator

import "github.com/snapp-incubator/soteria/internal/errors"

var (
	ErrInvalidSigningMethod = errors.ErrInvalidSigningMethod
	ErrIssNotFound          = errors.ErrIssNotFound
	ErrSubNotFound          = errors.ErrSubNotFound
	ErrInvalidClaims        = errors.ErrInvalidClaims
	ErrInvalidIP            = errors.ErrInvalidIP
	ErrInvalidAccessType    = errors.ErrInvalidAccessType
	ErrDecodeHashID         = errors.ErrDecodeHashID
	ErrInvalidSecret        = errors.ErrInvalidSecret
	ErrIncorrectPassword    = errors.ErrIncorrectPassword
)

type TopicNotAllowedError = errors.TopicNotAllowedError

type KeyNotFoundError = errors.KeyNotFoundError

type InvalidTopicError = errors.InvalidTopicError
