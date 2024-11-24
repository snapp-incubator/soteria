package metric_test

import (
	"errors"
	"testing"

	serrors "github.com/snapp-incubator/soteria/internal/errors"
	"github.com/snapp-incubator/soteria/internal/metric"
)

func TestAuthIncrement(t *testing.T) {
	t.Parallel()

	m := metric.NewAPIMetrics()

	m.AuthSuccess("snapp", "-")
	m.AuthFailed("snapp", "-", serrors.ErrInvalidSigningMethod)
	m.AuthFailed("snapp", "-", serrors.ErrIssNotFound)
	m.AuthFailed("snapp", "-", serrors.ErrSubNotFound)
	m.AuthFailed("snapp", "-", serrors.ErrInvalidClaims)
	m.AuthFailed("snapp", "-", serrors.ErrInvalidIP)

	m.AuthFailed("snapp", "-", serrors.ErrInvalidAccessType)
	m.AuthFailed("snapp", "-", serrors.ErrDecodeHashID)
	m.AuthFailed("snapp", "-", serrors.ErrInvalidSecret)
	m.AuthFailed("snapp", "-", serrors.ErrIncorrectPassword)
	m.AuthFailed("snapp", "-", &serrors.TopicNotAllowedError{
		Issuer:     "issuer",
		Sub:        "subject",
		AccessType: "1",
		Topic:      "topic",
		TopicType:  "pub",
	})
	m.AuthFailed("snapp", "-", &serrors.KeyNotFoundError{Issuer: "iss"})
	m.AuthFailed("snapp", "-", errors.ErrUnsupported)

	m.ACLSuccess("snapp")
	m.ACLFailed("snapp", serrors.ErrInvalidSigningMethod)
	m.ACLFailed("snapp", serrors.ErrIssNotFound)
	m.ACLFailed("snapp", serrors.ErrSubNotFound)
	m.ACLFailed("snapp", serrors.ErrInvalidClaims)
	m.ACLFailed("snapp", serrors.ErrInvalidIP)

	m.ACLFailed("snapp", serrors.ErrInvalidAccessType)
	m.ACLFailed("snapp", serrors.ErrDecodeHashID)
	m.ACLFailed("snapp", serrors.ErrInvalidSecret)
	m.ACLFailed("snapp", serrors.ErrIncorrectPassword)
	m.ACLFailed("snapp", &serrors.TopicNotAllowedError{
		Issuer:     "issuer",
		Sub:        "subject",
		AccessType: "1",
		Topic:      "topic",
		TopicType:  "pub",
	})
	m.ACLFailed("snapp", &serrors.KeyNotFoundError{Issuer: "iss"})
	m.ACLFailed("snapp", errors.ErrUnsupported)
}
