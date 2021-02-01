package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/metrics"
)

type SoteriaMetrics struct {
	Handler metrics.Handler
}

// NewMetrics creates and returns all metrics needed in Soteria
func NewMetrics() *SoteriaMetrics {
	statusCodesCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "dispatching",
		Subsystem: "soteria",
		Name:      "status_codes",
		Help:      "status codes observed from soteria and its all external calls",
	}, []string{"api", "service", "function", "code"})

	statusesCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "dispatching",
		Subsystem: "soteria",
		Name:      "statuses",
		Help:      "statuses observed from soteria and its all external calls",
	}, []string{"api", "service", "function", "status", "info"})

	responseTimesSummery := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "dispatching",
		Subsystem: "soteria",
		Name:      "response_times",
		Help:      "response times observed from and its all external calls",
	}, []string{"api", "service", "function"})

	prometheus.MustRegister(statusCodesCounter)
	prometheus.MustRegister(statusesCounter)
	prometheus.MustRegister(responseTimesSummery)

	h := metrics.Handler{
		StatusCodeCounterVec: statusCodesCounter,
		StatusCounterVec:     statusesCounter,
		ResponseTimeVec:      responseTimesSummery,
	}

	return &SoteriaMetrics{h}
}
