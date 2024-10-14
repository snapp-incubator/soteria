// nolint: testpackage
package authenticator

import (
	"errors"
	"testing"
)

// nolint: paralleltest
func TestIncrementWithErrorAuthCounter_NoError(_ *testing.T) {
	IncrementWithErrorAuthCounter("snapp", nil)
}

// nolint: paralleltest
func TestIncrementWithErrorAuthCounter_ErrInvalidSigningMethod(_ *testing.T) {
	IncrementWithErrorAuthCounter("snapp", ErrInvalidSigningMethod)
}

// nolint: paralleltest
func TestIncrementWithErrorAuthCounter_ErrIssNotFound(_ *testing.T) {
	IncrementWithErrorAuthCounter("snapp", ErrIssNotFound)
}

// nolint: paralleltest
func TestIncrementWithErrorAuthCounter_ErrSubNotFound(_ *testing.T) {
	IncrementWithErrorAuthCounter("snapp", ErrSubNotFound)
}

// nolint: paralleltest
func TestIncrementWithErrorAuthCounter_ErrInvalidClaims(_ *testing.T) {
	IncrementWithErrorAuthCounter("snapp", ErrInvalidClaims)
}

// nolint: paralleltest
func TestIncrementWithErrorAuthCounter_ErrInvalidIP(_ *testing.T) {
	IncrementWithErrorAuthCounter("snapp", ErrInvalidIP)
}

// nolint: paralleltest
func TestIncrementWithErrorAuthCounter_ErrInvalidAccessType(_ *testing.T) {
	IncrementWithErrorAuthCounter("snapp", ErrInvalidAccessType)
}

// nolint: paralleltest
func TestIncrementWithErrorAuthCounter_ErrDecodeHashID(_ *testing.T) {
	IncrementWithErrorAuthCounter("snapp", ErrDecodeHashID)
}

// nolint: paralleltest
func TestIncrementWithErrorAuthCounter_ErrInvalidSecret(_ *testing.T) {
	IncrementWithErrorAuthCounter("snapp", ErrInvalidSecret)
}

// nolint: paralleltest
func TestIncrementWithErrorAuthCounter_ErrIncorrectPassword(_ *testing.T) {
	IncrementWithErrorAuthCounter("snapp", ErrIncorrectPassword)
}

// nolint: paralleltest
func TestIncrementWithErrorAuthCounter_ErrTopicNotAllowed(_ *testing.T) {
	IncrementWithErrorAuthCounter("snapp", &TopicNotAllowedError{
		Issuer:     "issuer",
		Sub:        "subject",
		AccessType: "1",
		Topic:      "topic",
		TopicType:  "pub",
	})
}

// nolint: paralleltest
func TestIncrementWithErrorAuthCounter_ErrKeyNotFound(_ *testing.T) {
	IncrementWithErrorAuthCounter("snapp", &KeyNotFoundError{Issuer: "iss"})
}

// nolint: paralleltest
func TestIncrementWithErrorAuthCounter_UnknowError(_ *testing.T) {
	IncrementWithErrorAuthCounter("snapp", errors.ErrUnsupported)
}
