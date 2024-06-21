package authenticator

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var AuthenticateCounterMetric = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "dispatching",
	Subsystem: "soteria",
	Name:      "auth_total",
	Help:      "Total number of authentication attempts",
}, []string{"company", "status"})

func IncrementWithErrorAuthCounter(company string, err error) {
	var status string
	var topicNotAllowedErrorTarget *TopicNotAllowedError
	var keyNotFoundErrorTarget *KeyNotFoundError

	if err != nil {
		status = "success"
	} else if errors.Is(err, ErrInvalidSigningMethod) {
		status = "err_invalid_signing_method"
	} else if errors.Is(err, ErrIssNotFound) {
		status = "err_iss_not_found"
	} else if errors.Is(err, ErrSubNotFound) {
		status = "err_sub_not_found"
	} else if errors.Is(err, ErrInvalidClaims) {
		status = "err_invalid_claims"
	} else if errors.Is(err, ErrInvalidIP) {
		status = "err_invalid_ip"
	} else if errors.Is(err, ErrInvalidAccessType) {
		status = "err_invalid_access_type"
	} else if errors.Is(err, ErrDecodeHashID) {
		status = "err_decode_hash_id"
	} else if errors.Is(err, ErrInvalidSecret) {
		status = "err_invalid_secret"
	} else if errors.Is(err, ErrIncorrectPassword) {
		status = "err_incorrect_password"
	} else if errors.As(err, &topicNotAllowedErrorTarget) {
		status = "topic_not_allowed_error"
	} else if errors.As(err, &keyNotFoundErrorTarget) {
		status = "key_not_found_error"
	} else {
		status = "unknown_error"
	}

	AuthenticateCounterMetric.WithLabelValues(company, status).Inc()
}
