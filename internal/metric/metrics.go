package metric

import (
	"errors"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/snapp-incubator/soteria/internal/authenticator"
)

type AutoAuthenticatorMetrics struct {
	latency *prometheus.HistogramVec
}

type APIMetrics struct {
	auth *prometheus.CounterVec
	acl  *prometheus.CounterVec
}

func NewAutoAuthenticatorMetrics() *AutoAuthenticatorMetrics {
	m := &AutoAuthenticatorMetrics{
		latency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace:                       "platform",
			Subsystem:                       "soteria",
			Name:                            "auto_auth_latency_seconds",
			Help:                            "Automatic authentication latency in seconds",
			ConstLabels:                     prometheus.Labels{},
			Buckets:                         prometheus.DefBuckets,
			NativeHistogramBucketFactor:     0,
			NativeHistogramZeroThreshold:    0,
			NativeHistogramMaxBucketNumber:  0,
			NativeHistogramMinResetDuration: 0,
			NativeHistogramMaxZeroThreshold: 0,
			NativeHistogramMaxExemplars:     0,
			NativeHistogramExemplarTTL:      0,
		}, []string{"company", "status"}),
	}

	m.register()

	return m
}

// Latency measures latency in seconds.
func (m *AutoAuthenticatorMetrics) Latency(latency float64, company string, status int) {
	m.latency.WithLabelValues(company, strconv.Itoa(status)).Observe(latency)
}

func (m *AutoAuthenticatorMetrics) register() {
	register(m.latency)
}

func NewAPIMetrics() *APIMetrics {
	m := &APIMetrics{
		auth: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace:   "platform",
			Subsystem:   "soteria",
			Name:        "auth_total",
			Help:        "Total number of authentication attempts",
			ConstLabels: prometheus.Labels{},
		}, []string{"company", "status", "source"}),
		acl: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace:   "platform",
			Subsystem:   "soteria",
			Name:        "acl_total",
			Help:        "Total number of authorization attempts",
			ConstLabels: prometheus.Labels{},
		}, []string{"company", "status"}),
	}

	m.register()

	return m
}

func (m *APIMetrics) register() {
	register(m.acl)
	register(m.auth)
}

func (m *APIMetrics) AuthSuccess(company, source string) {
	m.auth.WithLabelValues(company, "success", source).Inc()
}

// nolint:cyclop
func (m *APIMetrics) AuthFailed(company, source string, err error) {
	var (
		status                     string
		topicNotAllowedErrorTarget *authenticator.TopicNotAllowedError
		keyNotFoundErrorTarget     *authenticator.KeyNotFoundError
	)

	switch {
	case errors.Is(err, authenticator.ErrInvalidSigningMethod):
		status = "err_invalid_signing_method"
	case errors.Is(err, authenticator.ErrIssNotFound):
		status = "err_iss_not_found"
	case errors.Is(err, authenticator.ErrSubNotFound):
		status = "err_sub_not_found"
	case errors.Is(err, authenticator.ErrInvalidClaims):
		status = "err_invalid_claims"
	case errors.Is(err, authenticator.ErrInvalidIP):
		status = "err_invalid_ip"
	case errors.Is(err, authenticator.ErrInvalidAccessType):
		status = "err_invalid_access_type"
	case errors.Is(err, authenticator.ErrDecodeHashID):
		status = "err_decode_hash_id"
	case errors.Is(err, authenticator.ErrInvalidSecret):
		status = "err_invalid_secret"
	case errors.Is(err, authenticator.ErrIncorrectPassword):
		status = "err_incorrect_password"
	case errors.As(err, &topicNotAllowedErrorTarget):
		status = "topic_not_allowed_error"
	case errors.As(err, &keyNotFoundErrorTarget):
		status = "key_not_found_error"
	default:
		status = "unknown_error"
	}

	m.auth.WithLabelValues(company, status, source).Inc()
}

func (m *APIMetrics) ACLSuccess(company string) {
	m.acl.WithLabelValues(company, "success").Inc()
}

// nolint:cyclop
func (m *APIMetrics) ACLFailed(company string, err error) {
	var (
		status                     string
		topicNotAllowedErrorTarget *authenticator.TopicNotAllowedError
		keyNotFoundErrorTarget     *authenticator.KeyNotFoundError
	)

	switch {
	case errors.Is(err, authenticator.ErrInvalidSigningMethod):
		status = "err_invalid_signing_method"
	case errors.Is(err, authenticator.ErrIssNotFound):
		status = "err_iss_not_found"
	case errors.Is(err, authenticator.ErrSubNotFound):
		status = "err_sub_not_found"
	case errors.Is(err, authenticator.ErrInvalidClaims):
		status = "err_invalid_claims"
	case errors.Is(err, authenticator.ErrInvalidIP):
		status = "err_invalid_ip"
	case errors.Is(err, authenticator.ErrInvalidAccessType):
		status = "err_invalid_access_type"
	case errors.Is(err, authenticator.ErrDecodeHashID):
		status = "err_decode_hash_id"
	case errors.Is(err, authenticator.ErrInvalidSecret):
		status = "err_invalid_secret"
	case errors.Is(err, authenticator.ErrIncorrectPassword):
		status = "err_incorrect_password"
	case errors.As(err, &topicNotAllowedErrorTarget):
		status = "topic_not_allowed_error"
	case errors.As(err, &keyNotFoundErrorTarget):
		status = "key_not_found_error"
	default:
		status = "unknown_error"
	}

	m.acl.WithLabelValues(company, status).Inc()
}
