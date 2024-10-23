package metric_test

import (
	"errors"
	"testing"

	"github.com/snapp-incubator/soteria/internal/authenticator"
	"github.com/snapp-incubator/soteria/internal/metric"
)

func TestAuthIncrement(t *testing.T) {
	t.Parallel()

	m := metric.NewAPIMetrics()

	m.AuthSuccess("snapp", "-")
	m.AuthFailed("snapp", "-", authenticator.ErrInvalidSigningMethod)
	m.AuthFailed("snapp", "-", authenticator.ErrIssNotFound)
	m.AuthFailed("snapp", "-", authenticator.ErrSubNotFound)
	m.AuthFailed("snapp", "-", authenticator.ErrInvalidClaims)
	m.AuthFailed("snapp", "-", authenticator.ErrInvalidIP)

	m.AuthFailed("snapp", "-", authenticator.ErrInvalidAccessType)
	m.AuthFailed("snapp", "-", authenticator.ErrDecodeHashID)
	m.AuthFailed("snapp", "-", authenticator.ErrInvalidSecret)
	m.AuthFailed("snapp", "-", authenticator.ErrIncorrectPassword)
	m.AuthFailed("snapp", "-", &authenticator.TopicNotAllowedError{
		Issuer:     "issuer",
		Sub:        "subject",
		AccessType: "1",
		Topic:      "topic",
		TopicType:  "pub",
	})
	m.AuthFailed("snapp", "-", &authenticator.KeyNotFoundError{Issuer: "iss"})
	m.AuthFailed("snapp", "-", errors.ErrUnsupported)

	m.ACLSuccess("snapp")
	m.ACLFailed("snapp", authenticator.ErrInvalidSigningMethod)
	m.ACLFailed("snapp", authenticator.ErrIssNotFound)
	m.ACLFailed("snapp", authenticator.ErrSubNotFound)
	m.ACLFailed("snapp", authenticator.ErrInvalidClaims)
	m.ACLFailed("snapp", authenticator.ErrInvalidIP)

	m.ACLFailed("snapp", authenticator.ErrInvalidAccessType)
	m.ACLFailed("snapp", authenticator.ErrDecodeHashID)
	m.ACLFailed("snapp", authenticator.ErrInvalidSecret)
	m.ACLFailed("snapp", authenticator.ErrIncorrectPassword)
	m.ACLFailed("snapp", &authenticator.TopicNotAllowedError{
		Issuer:     "issuer",
		Sub:        "subject",
		AccessType: "1",
		Topic:      "topic",
		TopicType:  "pub",
	})
	m.ACLFailed("snapp", &authenticator.KeyNotFoundError{Issuer: "iss"})
	m.ACLFailed("snapp", errors.ErrUnsupported)
}
