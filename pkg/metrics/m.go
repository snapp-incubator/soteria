package metrics

// Metrics is an interface used to capture metrics in soteria.
type Metrics interface {
	// ObserveStatusCode is the method for status code metrics
	ObserveStatusCode(api, serviceName, function string, code int)

	// ObserveStatus is the method for the status of the done operations
	ObserveStatus(api, serviceName, function, status, info string)

	// ObserveResponseTime is the method for times takes by a specific function
	ObserveResponseTime(api, serviceName, function string, time float64)
}
