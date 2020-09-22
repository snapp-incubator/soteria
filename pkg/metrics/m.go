package metrics

// Metrics is the interface for observing metrics.
type Metrics interface {
	Http
}

// Http is a interface used to capture metrics related to HTTP.
type Http interface {
	// ObserveStatusCode is the method for HTTP status code metrics
	ObserveStatusCode(serviceName string, function string, code int)

	// ObserveStatus is the method for the status of the done operations.
	ObserveStatus(serviceName string, function string, status string, info string)

	// ObserveResponseTime is the method for times takes by a specific function
	ObserveResponseTime(serviceName string, function string, time float64)
}
