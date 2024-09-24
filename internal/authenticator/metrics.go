package authenticator

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	AuthenticateCounter *prometheus.CounterVec
	ACLCounter          *prometheus.CounterVec
}

func NewMetrics() Metrics {
	m := Metrics{
		AuthenticateCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: "dispatching",
			Subsystem: "soteria",
			Name:      "auth_total",
			Help:      "Total number of authentication attempts",
		}, []string{"company", "status"}),

		ACLCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: "dispatching",
			Subsystem: "soteria",
			Name:      "acl_total",
			Help:      "Total number of acl attempts",
		}, []string{"company", "status"}),
	}

	return m
}

func (m Metrics) ACLIncrementWithError(company string, err error) {
	var (
		status                     string
		topicNotAllowedErrorTarget *TopicNotAllowedError
		keyNotFoundErrorTarget     *KeyNotFoundError
	)

	switch {
	case err == nil:
		status = "success"
	case errors.Is(err, ErrInvalidSigningMethod):
		status = "err_invalid_signing_method"
	case errors.Is(err, ErrIssNotFound):
		status = "err_iss_not_found"
	case errors.Is(err, ErrSubNotFound):
		status = "err_sub_not_found"
	case errors.Is(err, ErrInvalidClaims):
		status = "err_invalid_claims"
	case errors.Is(err, ErrInvalidIP):
		status = "err_invalid_ip"
	case errors.Is(err, ErrInvalidAccessType):
		status = "err_invalid_access_type"
	case errors.Is(err, ErrDecodeHashID):
		status = "err_decode_hash_id"
	case errors.Is(err, ErrInvalidSecret):
		status = "err_invalid_secret"
	case errors.Is(err, ErrIncorrectPassword):
		status = "err_incorrect_password"
	case errors.As(err, &topicNotAllowedErrorTarget):
		status = "topic_not_allowed_error"
	case errors.As(err, &keyNotFoundErrorTarget):
		status = "key_not_found_error"
	default:
		status = "unknown_error"
	}

	m.ACLCounter.WithLabelValues(company, status).Inc()
}

func (m Metrics) AuthIncrementWithError(company string, err error) {
	var (
		status                     string
		topicNotAllowedErrorTarget *TopicNotAllowedError
		keyNotFoundErrorTarget     *KeyNotFoundError
	)

	switch {
	case err == nil:
		status = "success"
	case errors.Is(err, ErrInvalidSigningMethod):
		status = "err_invalid_signing_method"
	case errors.Is(err, ErrIssNotFound):
		status = "err_iss_not_found"
	case errors.Is(err, ErrSubNotFound):
		status = "err_sub_not_found"
	case errors.Is(err, ErrInvalidClaims):
		status = "err_invalid_claims"
	case errors.Is(err, ErrInvalidIP):
		status = "err_invalid_ip"
	case errors.Is(err, ErrInvalidAccessType):
		status = "err_invalid_access_type"
	case errors.Is(err, ErrDecodeHashID):
		status = "err_decode_hash_id"
	case errors.Is(err, ErrInvalidSecret):
		status = "err_invalid_secret"
	case errors.Is(err, ErrIncorrectPassword):
		status = "err_incorrect_password"
	case errors.As(err, &topicNotAllowedErrorTarget):
		status = "topic_not_allowed_error"
	case errors.As(err, &keyNotFoundErrorTarget):
		status = "key_not_found_error"
	default:
		status = "unknown_error"
	}

	m.AuthenticateCounter.WithLabelValues(company, status).Inc()
}
