package factory

import (
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.snapp.ir/dispatching/soteria/pkg/metrics"
)

func GetMetrics(name string) metrics.Metrics {
	switch name {
	case "http":
		return getHttpMetric()
	default:
		return nil
	}
}

// getHttpMetric returns HttpCall metrics
func getHttpMetric() *metrics.HttpCall {
	h := &metrics.HttpCall{}
	statusCodesOps := prometheus.CounterOpts{
		Namespace: "dispatching",
		Subsystem: "soteria",
		Name:      "status_codes",
		Help:      "status codes observed from privent and its all external calls",
	}
	statusesOpts := prometheus.CounterOpts{
		Namespace: "dispatching",
		Subsystem: "soteria",
		Name:      "statuses",
		Help:      "statuses observed from privent and its all external calls",
	}
	sumOpts := prometheus.SummaryOpts{
		Namespace: "dispatching",
		Subsystem: "soteria",
		Name:      "response_times",
		Help:      "response times observed from privent and its all external calls",
	}
	h.StatusCodeCounterVec = prometheus.NewCounterVec(statusCodesOps, []string{"service", "function", "code"})
	h.StatusCounterVec = prometheus.NewCounterVec(statusesOpts, []string{"service", "function", "status", "info"})
	h.ResponseTimeVec = prometheus.NewSummaryVec(sumOpts, []string{"service", "function"})
	prometheus.MustRegister(h.StatusCodeCounterVec)
	prometheus.MustRegister(h.StatusCounterVec)
	prometheus.MustRegister(h.ResponseTimeVec)
	return h
}
