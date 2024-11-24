package authenticator

import "github.com/snapp-incubator/soteria/internal/error"

var (
	ErrInvalidSigningMethod = error.ErrInvalidSigningMethod
	ErrIssNotFound          = error.ErrIssNotFound
	ErrSubNotFound          = error.ErrSubNotFound
	ErrInvalidClaims        = error.ErrInvalidClaims
	ErrInvalidIP            = error.ErrInvalidIP
	ErrInvalidAccessType    = error.ErrInvalidAccessType
	ErrDecodeHashID         = error.ErrDecodeHashID
	ErrInvalidSecret        = error.ErrInvalidSecret
	ErrIncorrectPassword    = error.ErrIncorrectPassword
)

type TopicNotAllowedError = error.TopicNotAllowedError

type KeyNotFoundError = error.KeyNotFoundError

type InvalidTopicError = error.InvalidTopicError
