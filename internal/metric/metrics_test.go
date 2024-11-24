package metric_test

import (
	"errors"
	"testing"

	serror "github.com/snapp-incubator/soteria/internal/error"
	"github.com/snapp-incubator/soteria/internal/metric"
)

func TestAuthIncrement(t *testing.T) {
	t.Parallel()

	m := metric.NewAPIMetrics()

	m.AuthSuccess("snapp", "-")
	m.AuthFailed("snapp", "-", serror.ErrInvalidSigningMethod)
	m.AuthFailed("snapp", "-", serror.ErrIssNotFound)
	m.AuthFailed("snapp", "-", serror.ErrSubNotFound)
	m.AuthFailed("snapp", "-", serror.ErrInvalidClaims)
	m.AuthFailed("snapp", "-", serror.ErrInvalidIP)

	m.AuthFailed("snapp", "-", serror.ErrInvalidAccessType)
	m.AuthFailed("snapp", "-", serror.ErrDecodeHashID)
	m.AuthFailed("snapp", "-", serror.ErrInvalidSecret)
	m.AuthFailed("snapp", "-", serror.ErrIncorrectPassword)
	m.AuthFailed("snapp", "-", &serror.TopicNotAllowedError{
		Issuer:     "issuer",
		Sub:        "subject",
		AccessType: "1",
		Topic:      "topic",
		TopicType:  "pub",
	})
	m.AuthFailed("snapp", "-", &serror.KeyNotFoundError{Issuer: "iss"})
	m.AuthFailed("snapp", "-", errors.ErrUnsupported)

	m.ACLSuccess("snapp")
	m.ACLFailed("snapp", serror.ErrInvalidSigningMethod)
	m.ACLFailed("snapp", serror.ErrIssNotFound)
	m.ACLFailed("snapp", serror.ErrSubNotFound)
	m.ACLFailed("snapp", serror.ErrInvalidClaims)
	m.ACLFailed("snapp", serror.ErrInvalidIP)

	m.ACLFailed("snapp", serror.ErrInvalidAccessType)
	m.ACLFailed("snapp", serror.ErrDecodeHashID)
	m.ACLFailed("snapp", serror.ErrInvalidSecret)
	m.ACLFailed("snapp", serror.ErrIncorrectPassword)
	m.ACLFailed("snapp", &serror.TopicNotAllowedError{
		Issuer:     "issuer",
		Sub:        "subject",
		AccessType: "1",
		Topic:      "topic",
		TopicType:  "pub",
	})
	m.ACLFailed("snapp", &serror.KeyNotFoundError{Issuer: "iss"})
	m.ACLFailed("snapp", errors.ErrUnsupported)
}
