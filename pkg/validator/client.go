package validator

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	ServiceNameHeader = "X-Service-Name"

	validateURI    = "/api/v3/internal/validate"
	authHeader     = "Authorization"
	userDataHeader = "X-User-Data"
	modeQueryParam = "mode"
)

var (
	ErrEmptyServiceName      = errors.New("x-service-name can not be empty")
	ErrInvalidJWT            = errors.New("invalid jwt")
	ErrInvalidUserDataHeader = errors.New("invalid X-User-Data header")
	ErrRequestFailed         = errors.New("validator request failed")
)

type Client struct {
	baseURL    string
	client     *http.Client
	timeout    time.Duration
	isOptional bool
}

// New creates a new Client with default attributes.
func New(url string, timeout time.Duration) Client {
	return Client{
		baseURL:    url,
		client:     new(http.Client),
		timeout:    timeout,
		isOptional: false,
	}
}

// WithOptionalValidate enables you to bypass the signature validation.
// If the validation of the JWT signature is optional for you, and you just want to extract
// the payload from the token, you can use the client in `WithOptionalValidate` mode.
func (c *Client) WithOptionalValidate() {
	c.isOptional = true
}

// Validate gets the parent context, headers, and JWT token and calls the validate API of the JWT validator service.
// The parent context is helpful in canceling the process in the upper hand (a function that used the SDK) and in case
// you have something like tracing spans in your context and want to extend these things in your custom HTTP handler.
// Otherwise, you can use `context.Background()`.
// The headers argument is used when you want to pass some headers like user-agent,
// X-Service-Name, X-App-Name, X-App-Version and
// X-App-Version-Code to the validator. It is extremely recommended to pass these headers (if you have them) because
// it increases the visibility in the logs and metrics of the JWT Validator service.
// You must place your Authorization header content in the bearerToken argument.
// Consider that the bearerToken must contain Bearer keyword and JWT.
// For `X-Service-Name` you should put your project/service name in this header.
// nolint: funlen, cyclop
func (c *Client) Validate(parentCtx context.Context, headers http.Header, bearerToken string) error {
	if headers.Get(ServiceNameHeader) == "" {
		return ErrEmptyServiceName
	}

	segments := strings.Split(bearerToken, " ")
	if len(segments) < 2 || strings.ToLower(segments[0]) != "bearer" {
		return ErrInvalidJWT
	}

	ctx, cancel := context.WithTimeout(parentCtx, c.timeout)
	defer cancel()

	url := c.baseURL + validateURI

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("validator creating request failed %w", err)
	}

	request.Header = headers
	request.Header.Set(authHeader, bearerToken)

	query := request.URL.Query()
	if c.isOptional {
		query.Add(modeQueryParam, "optional")
	}

	request.URL.RawQuery = query.Encode()

	response, err := c.client.Do(request)
	if err != nil {
		return fmt.Errorf("validator sending request failed %w", err)
	}

	closeBody(response)

	if response.StatusCode != http.StatusOK {
		return ErrRequestFailed
	}

	userDataHeader := response.Header.Get(userDataHeader)
	if userDataHeader == "" {
		return ErrInvalidJWT
	}

	return nil
}

// closeBody to avoid memory leak when reusing http connection.
func closeBody(response *http.Response) {
	if response != nil {
		_ = response.Body.Close()
	}
}
