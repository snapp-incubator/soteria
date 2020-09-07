package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

// HttpCall is a implementation of `Http` and `Metrics`
type HttpCall struct {
	StatusCodeCounterVec *prometheus.CounterVec
	StatusCounterVec     *prometheus.CounterVec
	ResponseTimeVec      *prometheus.SummaryVec
}

func (s *HttpCall) ObserveStatusCode(serviceName string, function string, code int) {
	s.StatusCodeCounterVec.With(prometheus.Labels{
		"service":  serviceName,
		"function": function,
		"code":     strconv.Itoa(code),
	}).Inc()
}

func (s *HttpCall) ObserveStatus(serviceName string, function string, status string, info string) {
	s.StatusCounterVec.With(prometheus.Labels{
		"service":  serviceName,
		"function": function,
		"status":   status,
		"info":     info,
	}).Inc()
}

func (s *HttpCall) ObserveResponseTime(serviceName string, function string, time float64) {
	s.ResponseTimeVec.With(prometheus.Labels{
		"service":  serviceName,
		"function": function,
	}).Observe(time)
}
