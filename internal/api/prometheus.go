package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PrometheusMiddleware collects HTTP metrics for Fiber v3.
type PrometheusMiddleware struct {
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	requestsInFlight prometheus.Gauge
}

// NewPrometheusMiddleware creates a new PrometheusMiddleware and registers its
// collectors with the given registerer.
func NewPrometheusMiddleware(reg prometheus.Registerer) *PrometheusMiddleware {
	m := &PrometheusMiddleware{
		requestsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "platform",
			Subsystem: "soteria",
			Name:      "http_requests_total",
			Help:      "Total number of HTTP requests.",
		}, []string{"status_code", "method", "path"}),
		requestDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "platform",
			Subsystem: "soteria",
			Name:      "http_request_duration_seconds",
			Help:      "Duration of HTTP requests in seconds.",
			Buckets:   prometheus.DefBuckets,
		}, []string{"status_code", "method", "path"}),
		requestsInFlight: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "platform",
			Subsystem: "soteria",
			Name:      "http_requests_in_progress",
			Help:      "Number of HTTP requests currently in progress.",
		}),
	}

	reg.MustRegister(m.requestsTotal, m.requestDuration, m.requestsInFlight)

	return m
}

// RegisterAt registers the /metrics endpoint on the given Fiber app.
func (m *PrometheusMiddleware) RegisterAt(app *fiber.App, path string) {
	app.Get(path, adaptor.HTTPHandler(promhttp.Handler()))
}

// Handler is the Fiber middleware that records metrics for each request.
func (m *PrometheusMiddleware) Handler(c fiber.Ctx) error {
	// Skip the metrics endpoint itself.
	if c.Path() == "/metrics" {
		return c.Next()
	}

	m.requestsInFlight.Inc()
	start := time.Now()

	err := c.Next()

	duration := time.Since(start).Seconds()
	status := strconv.Itoa(c.Response().StatusCode())
	method := c.Method()
	path := c.Route().Path

	if path == "" {
		path = c.Path()
	}

	m.requestsTotal.WithLabelValues(status, method, path).Inc()
	m.requestDuration.WithLabelValues(status, method, path).Observe(duration)
	m.requestsInFlight.Dec()

	if err != nil {
		return &fiber.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}
