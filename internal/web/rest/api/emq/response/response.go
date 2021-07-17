package response

// Response is the response structure of the REST API.
type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
