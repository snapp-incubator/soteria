package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

// Handler is a implementation `Metrics`
type Handler struct {
	StatusCodeCounterVec *prometheus.CounterVec
	StatusCounterVec     *prometheus.CounterVec
	ResponseTimeVec      *prometheus.SummaryVec
}

func (h *Handler) ObserveStatusCode(api, serviceName, function string, code int) {
	h.StatusCodeCounterVec.With(prometheus.Labels{
		"api":      api,
		"service":  serviceName,
		"function": function,
		"code":     strconv.Itoa(code),
	}).Inc()
}

func (h *Handler) ObserveStatus(api, serviceName, function, status, info string) {
	h.StatusCounterVec.With(prometheus.Labels{
		"api":      api,
		"service":  serviceName,
		"function": function,
		"status":   status,
		"info":     info,
	}).Inc()
}

func (h *Handler) ObserveResponseTime(api, serviceName, function string, time float64) {
	h.ResponseTimeVec.With(prometheus.Labels{
		"api":      api,
		"service":  serviceName,
		"function": function,
	}).Observe(time)
}
